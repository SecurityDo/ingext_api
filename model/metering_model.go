package model

// Billing metering (usage) API models.
//
// These mirror the ledger written by ingext_remote's metering package: the
// account's api pod measures its own usage once a day from its cluster's
// VictoriaMetrics and records it in the shared Postgres. Ingext reports
// quantities and evidence and stores no prices, no SKUs and nothing from
// Stripe -- turning these numbers into money belongs to the caller.
//
// THE ONE RULE A CONSUMER MUST NOT GET WRONG: a meter that was not measured is
// null, and null is not zero. Every byte field here is a *int64 for that
// reason. Treating a nil as 0 turns "we could not measure this tenant" into
// "this tenant used nothing", which is a silent under-bill. Use the Bytes
// helper, or check the pointer.

// Meter names, as they appear in Meters and MeterStatuses.
const (
	MeterEventwatch         = "eventwatch_bytes"
	MeterProcessed          = "processed_bytes"
	MeterDeleted            = "deleted_bytes"
	MeterInput              = "input_bytes"
	MeterPlatformDatalake   = "platform_datalake_bytes"
	MeterLakeIngress        = "lake_ingress_bytes"
	MeterLakeSearch         = "lake_search_bytes"
	MeterLakeRealtimeSearch = "lake_realtime_search_bytes"
)

// MeterNames is the canonical order.
var MeterNames = []string{
	MeterEventwatch, MeterProcessed, MeterDeleted, MeterInput,
	MeterPlatformDatalake, MeterLakeIngress, MeterLakeSearch, MeterLakeRealtimeSearch,
}

// Per-meter status. Zero and Unavailable are the distinction that matters:
// "measured, and it was nothing" versus "we do not know".
const (
	MeterStatusMeasured      = "measured"
	MeterStatusZero          = "zero"
	MeterStatusUnavailable   = "unavailable"
	MeterStatusNotApplicable = "not_applicable"
)

// Day states in a UsageDay.
const (
	// DayStateClosed: finished and immutable.
	DayStateClosed = "closed"
	// DayStateOpen: today, provisional, recomputed every few minutes.
	DayStateOpen = "open"
	// DayStateMissing: the day elapsed and no row exists. Materialised by the
	// API rather than omitted, so a short array cannot be mistaken for a short
	// month. It is NOT a day of zero usage.
	DayStateMissing = "missing"
	// DayStateInProgress: today, before any row has been written.
	DayStateInProgress = "in_progress"
)

// Day-level status on a closed row.
const (
	// DayStatusFinal: every meter resolved. The only billable status.
	DayStatusFinal = "final"
	// DayStatusIncomplete: written only when a day aged out of the metrics
	// retention window without ever resolving. It records a permanent gap and
	// is never billable.
	DayStatusIncomplete = "incomplete"
)

// Paid-user providers.
const (
	ProviderMicrosoft365    = "microsoft365"
	ProviderGoogleWorkspace = "google_workspace"
)

// Paid-user observation status.
const (
	// UserStatusOK: the count is exact.
	UserStatusOK = "ok"
	// UserStatusVerifiedPaidLowerBound: every account counted is definitely
	// paid, but at least one active account holds only SKUs the server's table
	// does not know, so the true figure may be higher. Still final and still
	// billable.
	UserStatusVerifiedPaidLowerBound = "verified_paid_lower_bound"
)

// PaidUserCount is one provider's headcount for one day.
//
// A provider that is not integrated on the account produces NO entry at all,
// which is different from an entry of 0. Absence means "no directory here";
// zero means "a directory with nobody licensed in it".
type PaidUserCount struct {
	Provider string `json:"provider"`
	Quantity int64  `json:"quantity"`
	Status   string `json:"status"`
	Method   string `json:"method"`
}

// UsageDay is one account-day of the ledger.
type UsageDay struct {
	BillingDate string `json:"billingDate"` // YYYY-MM-DD, UTC
	State       string `json:"state"`
	Status      string `json:"status,omitempty"`
	Collector   string `json:"collector,omitempty"`

	PeriodStart string `json:"periodStart,omitempty"`
	// ObservedEnd is the exclusive end actually covered. On an open day it is
	// the last completed hour, and it only ever moves forward.
	ObservedEnd string `json:"observedEnd,omitempty"`
	PeriodEnd   string `json:"periodEnd,omitempty"`

	// Meters is keyed by the Meter* constants. A nil value means NOT MEASURED.
	Meters        map[string]*int64 `json:"meters,omitempty"`
	MeterStatuses map[string]string `json:"meterStatuses,omitempty"`

	PaidUsers []PaidUserCount `json:"paidUsers,omitempty"`

	ContentSha256 string `json:"contentSha256,omitempty"`
}

// Bytes returns a meter's value and whether it was measured at all.
//
// Prefer this over indexing Meters directly: the zero value of the map lookup
// is a nil pointer, and dereferencing or defaulting it to 0 is exactly the
// mistake that silently under-bills.
func (d *UsageDay) Bytes(meter string) (value int64, measured bool) {
	if d == nil || d.Meters == nil {
		return 0, false
	}
	p, ok := d.Meters[meter]
	if !ok || p == nil {
		return 0, false
	}
	return *p, true
}

// Billable reports whether this day may be used to compute a charge.
func (d *UsageDay) Billable() bool {
	return d != nil && d.State == DayStateClosed && d.Status == DayStatusFinal
}

// UsageResponse is the reply from metering_daily_list.
type UsageResponse struct {
	SchemaVersion string `json:"schemaVersion"`
	Account       string `json:"account"`
	TenantKey     string `json:"tenantKey,omitempty"`
	From          string `json:"from"`
	To            string `json:"to"`

	// StoreAvailable is false when the ledger could not be reached. Days will
	// be empty, and that emptiness means NOTHING about the tenant's usage --
	// see Reason. Never sum a response with StoreAvailable false.
	StoreAvailable bool   `json:"storeAvailable"`
	Reason         string `json:"reason,omitempty"`

	Days []*UsageDay `json:"days"`
}

// BillableDays returns only the days that may be charged for.
func (r *UsageResponse) BillableDays() []*UsageDay {
	if r == nil {
		return nil
	}
	var out []*UsageDay
	for _, d := range r.Days {
		if d.Billable() {
			out = append(out, d)
		}
	}
	return out
}

// MeterTotal sums one meter over the billable days, and reports how many days
// were skipped because the meter was not measured on them.
//
// The skipped count is returned rather than hidden because a total over an
// incomplete month is a lower bound, and a caller that does not know how many
// days are missing cannot tell that.
func (r *UsageResponse) MeterTotal(meter string) (total int64, days int, skipped int) {
	for _, d := range r.BillableDays() {
		if v, ok := d.Bytes(meter); ok {
			total += v
			days++
			continue
		}
		skipped++
	}
	return total, days, skipped
}

// Attempt is one recorded collection try, whatever the outcome.
//
// The attempts ledger is what separates "this tenant-day has no data" from
// "nobody ever looked". A failure is an Attempt with no observation, never an
// observation of zero.
type Attempt struct {
	BillingDate string                 `json:"billingDate"`
	MeterFamily string                 `json:"meterFamily"` // capacity | paid_users
	Provider    string                 `json:"provider,omitempty"`
	Status      string                 `json:"status"` // succeeded | failed | skipped
	Method      string                 `json:"method"`
	ErrorCode   string                 `json:"errorCode,omitempty"`
	Actor       string                 `json:"actor"`
	StartedAt   string                 `json:"startedAt"`
	CompletedAt string                 `json:"completedAt"`
	Evidence    map[string]interface{} `json:"evidence,omitempty"`
}

// AttemptsResponse is the reply from metering_attempts.
type AttemptsResponse struct {
	StoreAvailable bool       `json:"storeAvailable"`
	Attempts       []*Attempt `json:"attempts"`
}

// CollectDayResponse is the reply from metering_collect_day.
//
// Closed false is a normal outcome, not an error: some meter did not resolve,
// so the day was deliberately not written and will be retried. Read the
// attempts ledger to find out which meter and why.
type CollectDayResponse struct {
	BillingDate string `json:"billingDate"`
	Closed      bool   `json:"closed"`
	Actor       string `json:"actor"`
}

// SinkClassification is the reply from metering_sink_classification: how the
// collector currently sorts this account's datasinks into meters.
//
// Ambiguous holds sinks whose meter differs between the platform's resident and
// job execution modes. They are counted as eventwatch (which under-counts by at
// most that sink) rather than processed (which would double-bill it), and
// listed here so the choice is visible rather than silent.
type SinkClassification struct {
	Account    string   `json:"account"`
	Processed  []string `json:"processed"`
	Eventwatch []string `json:"eventwatch"`
	Datalake   []string `json:"datalake"`
	Ambiguous  []string `json:"ambiguous,omitempty"`
	Rejected   []string `json:"rejected,omitempty"`
}

// GridUsageAccount is one tenant's slice of a provider-level response.
//
// Error is set when that tenant could not be reached. Such an entry carries no
// days and MUST NOT be read as a tenant with no usage.
type GridUsageAccount struct {
	Account        string      `json:"account"`
	State          string      `json:"state,omitempty"`
	Error          string      `json:"error,omitempty"`
	StoreAvailable bool        `json:"storeAvailable"`
	Reason         string      `json:"reason,omitempty"`
	Days           []*UsageDay `json:"days,omitempty"`
}

// GridUsageSummary counts how the fan-out went.
type GridUsageSummary struct {
	Requested int `json:"requested"`
	Succeeded int `json:"succeeded"`
	// Failed non-zero means the document is INCOMPLETE. Any provider total
	// computed from it is a lower bound.
	Failed int `json:"failed"`
}

// GridUsageResponse is the reply from grid_metering_daily_list.
type GridUsageResponse struct {
	SchemaVersion string              `json:"schemaVersion"`
	Grid          string              `json:"grid"`
	From          string              `json:"from"`
	To            string              `json:"to"`
	TenantKey     string              `json:"tenantKey,omitempty"`
	Accounts      []*GridUsageAccount `json:"accounts"`
	Summary       GridUsageSummary    `json:"summary"`
	// OutOfScope lists requested accounts this caller may not see. They are
	// reported rather than dropped so "you may not see this" cannot be mistaken
	// for "this has no usage".
	OutOfScope []string `json:"outOfScope,omitempty"`
}

// Complete reports whether every requested tenant answered. A false here means
// any total derived from the response is a lower bound.
func (r *GridUsageResponse) Complete() bool {
	return r != nil && r.Summary.Failed == 0
}

// DayIndex renders the day as a YYYYMMDD dayIndex, the form ingext uses in
// index names, dump paths and the lake's @dayIndex field.
func (d *UsageDay) DayIndex() string {
	if d == nil {
		return ""
	}
	s := d.BillingDate
	if len(s) == 10 && s[4] == '-' && s[7] == '-' {
		return s[0:4] + s[5:7] + s[8:10]
	}
	return s
}
