package api

import (
	"fmt"
	"strings"
	"time"

	"github.com/SecurityDo/ingext_api/client"
	"github.com/SecurityDo/ingext_api/model"
)

// MeteringService calls the metering_* (billing usage) endpoints.
//
// Ingext measures usage and stores quantities and evidence; it stores no
// prices, no SKUs and nothing from Stripe. These bindings report what was
// measured and, just as importantly, what was NOT -- see model.UsageDay.Bytes
// and UsageResponse.StoreAvailable.
//
// Three ways to reach a tenant, all the same endpoints:
//
//	direct        client points at the tenant, no gridaccount
//	via provider  client points at the provider, SetGridAccount("<tenant>")
//	provider-wide GridUsage, which fans out over every tenant in one call
type MeteringService struct {
	client *client.IngextClient
}

// NewMeteringService constructs a MeteringService backed by the provided client.
func NewMeteringService(client *client.IngextClient) *MeteringService {
	return &MeteringService{client: client}
}

func (s *MeteringService) call(function string, payload interface{}, out interface{}) error {
	return ApiCall(s.client, function, payload, out)
}

func (s *MeteringService) gridCall(function string, payload interface{}, out interface{}) error {
	return ApiCallWithPrefix(s.client, "api/grid", function, payload, out)
}

// maxUsageRangeDays matches the server's own limit. A whole-month total in this
// design is the sum of closed days, never one wide query, so a request spanning
// years is a misunderstanding rather than a big report.
const maxUsageRangeDays = 400

// normalizeRange accepts either date form on both ends and returns the wire
// form. Mixed forms are fine -- an operator pasting a dayIndex from an index
// name next to a hand-typed date should not have to care.
func normalizeRange(from, to string) (string, string, error) {
	nf, err := NormalizeDate(from)
	if err != nil {
		return "", "", fmt.Errorf("from: %w", err)
	}
	nt, err := NormalizeDate(to)
	if err != nil {
		return "", "", fmt.Errorf("to: %w", err)
	}
	f, _ := time.Parse(dateLayout, nf)
	t, _ := time.Parse(dateLayout, nt)
	if t.Before(f) {
		return "", "", fmt.Errorf("to (%s) is before from (%s)", nt, nf)
	}
	if t.Sub(f) > maxUsageRangeDays*24*time.Hour {
		return "", "", fmt.Errorf("range %s..%s exceeds %d days", nf, nt, maxUsageRangeDays)
	}
	return nf, nt, nil
}

// UsageOptions narrows a usage query.
type UsageOptions struct {
	// TenantKey selects an MSSP sub-customer. Always "" today: the platform
	// metrics carry no customer label, so capacity cannot be subdivided.
	TenantKey string
	// IncludeOpen adds today's provisional row. Its numbers grow through the
	// day and are NOT billable -- only a closed, final day is.
	IncludeOpen bool
}

// DailyUsage returns the daily ledger for [from, to] inclusive, both YYYY-MM-DD
// in UTC.
//
// Days that elapsed with no row come back with State "missing" rather than
// being omitted, so a caller cannot mistake a short array for a short month.
// Check StoreAvailable before summing anything: when it is false the ledger was
// unreachable and an empty Days says nothing about the tenant's usage.
func (s *MeteringService) DailyUsage(from, to string, opts *UsageOptions) (*model.UsageResponse, error) {
	from, to, err := normalizeRange(from, to)
	if err != nil {
		return nil, err
	}
	kargs := map[string]interface{}{"from": from, "to": to}
	if opts != nil {
		if opts.TenantKey != "" {
			kargs["tenantKey"] = opts.TenantKey
		}
		if opts.IncludeOpen {
			kargs["includeOpen"] = true
		}
	}
	var resp model.UsageResponse
	if err := s.call("metering_daily_list", kargs, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// Attempts returns the collection attempt ledger for [from, to].
//
// This is the operational view. A day absent from DailyUsage and a day whose
// collection failed look identical there; here they do not. Grep ErrorCode to
// find every tenant-day affected by one cause.
func (s *MeteringService) Attempts(from, to string, limit int, opts *UsageOptions) (*model.AttemptsResponse, error) {
	from, to, err := normalizeRange(from, to)
	if err != nil {
		return nil, err
	}
	kargs := map[string]interface{}{"from": from, "to": to}
	if limit > 0 {
		kargs["limit"] = limit
	}
	if opts != nil && opts.TenantKey != "" {
		kargs["tenantKey"] = opts.TenantKey
	}
	var resp model.AttemptsResponse
	if err := s.call("metering_attempts", kargs, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// CollectDay collects and closes one past UTC day on demand.
//
// It cannot rewrite a day that is already closed -- the ledger refuses that --
// so the worst it can do is re-derive one that failed. Closed false in the
// reply is a normal outcome meaning some meter did not resolve; read Attempts
// to find out which.
//
// The date must be strictly in the past: a day still in progress has no
// complete final hour and the server refuses it.
func (s *MeteringService) CollectDay(date string, opts *UsageOptions) (*model.CollectDayResponse, error) {
	date, err := NormalizeDate(date)
	if err != nil {
		return nil, err
	}
	d, _ := time.Parse(dateLayout, date)
	today := time.Now().UTC().Truncate(24 * time.Hour)
	if !d.Before(today) {
		return nil, fmt.Errorf("date %s is not in the past; a day that has not finished cannot be closed", date)
	}
	kargs := map[string]interface{}{"date": date}
	if opts != nil && opts.TenantKey != "" {
		kargs["tenantKey"] = opts.TenantKey
	}
	var resp model.CollectDayResponse
	if err := s.call("metering_collect_day", kargs, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// SinkClassification reports how the collector currently sorts this account's
// datasinks into meters.
//
// Worth checking first whenever processed_bytes looks wrong: it is the one
// meter whose value depends on a judgement about the topology rather than on a
// metric label.
func (s *MeteringService) SinkClassification() (*model.SinkClassification, error) {
	var resp model.SinkClassification
	if err := s.call("metering_sink_classification", map[string]interface{}{}, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// GridUsage collects the daily ledger for every tenant of a provider in one
// call. It must be issued against a PROVIDER site, not a tenant.
//
// accounts narrows the set; empty means every tenant the caller may see. It can
// only narrow -- anything named that is out of scope comes back in OutOfScope
// rather than being silently dropped.
//
// Check Complete() before totalling anything: a tenant that could not be
// reached appears with an Error and no days, and treating that as zero usage
// under-bills the provider.
func (s *MeteringService) GridUsage(from, to string, accounts []string, opts *UsageOptions) (*model.GridUsageResponse, error) {
	from, to, err := normalizeRange(from, to)
	if err != nil {
		return nil, err
	}
	kargs := map[string]interface{}{"from": from, "to": to}
	if len(accounts) > 0 {
		cleaned := make([]string, 0, len(accounts))
		for _, a := range accounts {
			if a = strings.TrimSpace(a); a != "" {
				cleaned = append(cleaned, a)
			}
		}
		kargs["accounts"] = cleaned
	}
	if opts != nil {
		if opts.TenantKey != "" {
			kargs["tenantKey"] = opts.TenantKey
		}
		if opts.IncludeOpen {
			kargs["includeOpen"] = true
		}
	}
	var resp model.GridUsageResponse
	if err := s.gridCall("grid_metering_daily_list", kargs, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// timeNowUTCDate is the current UTC calendar day, split out so tests can state
// "today" without duplicating the format string.
func timeNowUTCDate() string {
	return time.Now().UTC().Format("2006-01-02")
}
