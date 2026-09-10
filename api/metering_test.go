package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/SecurityDo/ingext_api/client"
	fsb "github.com/SecurityDo/ingext_api/fsb"
	"github.com/SecurityDo/ingext_api/model"
)

// newMeteringServiceForTest serves a fixed response and records the path and
// kargs of each call.
func newMeteringServiceForTest(t *testing.T, response interface{}, path *string, kargs *map[string]interface{}) *MeteringService {
	t.Helper()
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if path != nil {
			*path = r.URL.Path
		}
		var req fsb.CallRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatalf("failed to decode request: %v", err)
		}
		if kargs != nil && req.Kargs != nil {
			m := map[string]interface{}{}
			if err := json.Unmarshal(req.Kargs.GetBytes(), &m); err != nil {
				t.Fatalf("failed to decode kargs: %v", err)
			}
			*kargs = m
		}
		body, err := json.Marshal(response)
		if err != nil {
			t.Fatalf("failed to marshal response: %v", err)
		}
		w.Header().Set("Content-Type", "application/json")
		payload := fsb.CallResponse{Verdict: "OK", Response: fsb.NewJNodeByte(body)}
		if err := json.NewEncoder(w).Encode(payload); err != nil {
			t.Fatalf("failed to write response: %v", err)
		}
	}))
	t.Cleanup(ts.Close)
	return NewMeteringService(client.NewIngextClient(ts.URL, "", false, nil))
}

func i64(v int64) *int64 { return &v }

func TestMeteringDailyUsage(t *testing.T) {
	var path string
	var kargs map[string]interface{}
	resp := model.UsageResponse{
		SchemaVersion: "ingext-metering/v1", Account: "titan",
		From: "2026-09-08", To: "2026-09-09", StoreAvailable: true,
		Days: []*model.UsageDay{
			{BillingDate: "2026-09-08", State: model.DayStateMissing},
			{BillingDate: "2026-09-09", State: model.DayStateClosed, Status: model.DayStatusFinal,
				Meters: map[string]*int64{
					model.MeterEventwatch: i64(172209747),
					model.MeterProcessed:  i64(0),
					model.MeterInput:      nil,
				},
				PaidUsers: []model.PaidUserCount{{Provider: model.ProviderMicrosoft365, Quantity: 11}},
			},
		},
	}
	svc := newMeteringServiceForTest(t, resp, &path, &kargs)

	got, err := svc.DailyUsage("2026-09-08", "2026-09-09", &UsageOptions{IncludeOpen: true})
	if err != nil {
		t.Fatalf("DailyUsage: %v", err)
	}
	if path != "/api/ds/metering_daily_list" {
		t.Fatalf("path = %s", path)
	}
	if kargs["from"] != "2026-09-08" || kargs["to"] != "2026-09-09" || kargs["includeOpen"] != true {
		t.Fatalf("kargs = %+v", kargs)
	}
	if len(got.Days) != 2 {
		t.Fatalf("days = %d", len(got.Days))
	}
}

// The distinction the whole API exists to preserve: an unmeasured meter is not
// a zero. Bytes must report it as absent, and a measured zero as present.
func TestUsageDayBytesSeparatesNullFromZero(t *testing.T) {
	d := &model.UsageDay{Meters: map[string]*int64{
		model.MeterEventwatch: i64(500),
		model.MeterProcessed:  i64(0),
		model.MeterInput:      nil,
	}}
	if v, ok := d.Bytes(model.MeterEventwatch); !ok || v != 500 {
		t.Fatalf("eventwatch = %d ok=%v", v, ok)
	}
	if v, ok := d.Bytes(model.MeterProcessed); !ok || v != 0 {
		t.Fatalf("a MEASURED zero must report as present: %d ok=%v", v, ok)
	}
	if _, ok := d.Bytes(model.MeterInput); ok {
		t.Fatal("an unmeasured meter must report as absent, not as 0")
	}
	if _, ok := d.Bytes(model.MeterLakeSearch); ok {
		t.Fatal("a meter missing from the map must report as absent")
	}
	// And on a nil day, rather than panicking.
	var nilDay *model.UsageDay
	if _, ok := nilDay.Bytes(model.MeterEventwatch); ok {
		t.Fatal("nil day must report absent")
	}
}

// Only a closed, final day may be charged for. Open days grow; incomplete days
// record a permanent gap.
func TestBillableDays(t *testing.T) {
	r := &model.UsageResponse{Days: []*model.UsageDay{
		{BillingDate: "1", State: model.DayStateClosed, Status: model.DayStatusFinal},
		{BillingDate: "2", State: model.DayStateClosed, Status: model.DayStatusIncomplete},
		{BillingDate: "3", State: model.DayStateOpen, Status: model.DayStatusFinal},
		{BillingDate: "4", State: model.DayStateMissing},
	}}
	got := r.BillableDays()
	if len(got) != 1 || got[0].BillingDate != "1" {
		t.Fatalf("billable = %+v", got)
	}
}

// A total over days where the meter was unmeasured must report how many it
// skipped, so the caller knows the figure is a lower bound.
func TestMeterTotalReportsSkippedDays(t *testing.T) {
	r := &model.UsageResponse{Days: []*model.UsageDay{
		{State: model.DayStateClosed, Status: model.DayStatusFinal,
			Meters: map[string]*int64{model.MeterEventwatch: i64(100)}},
		{State: model.DayStateClosed, Status: model.DayStatusFinal,
			Meters: map[string]*int64{model.MeterEventwatch: i64(0)}},
		{State: model.DayStateClosed, Status: model.DayStatusFinal,
			Meters: map[string]*int64{model.MeterEventwatch: nil}},
		{State: model.DayStateMissing},
	}}
	total, days, skipped := r.MeterTotal(model.MeterEventwatch)
	if total != 100 || days != 2 || skipped != 1 {
		t.Fatalf("total=%d days=%d skipped=%d, want 100/2/1", total, days, skipped)
	}
}

func TestMeteringRangeValidation(t *testing.T) {
	svc := NewMeteringService(client.NewIngextClient("http://127.0.0.1:1", "", false, nil))
	for _, c := range []struct{ from, to, why string }{
		{"2026-9-8", "2026-09-09", "non-padded date"},
		{"08/09/2026", "2026-09-09", "wrong format"},
		{"2026-09-10", "2026-09-09", "to before from"},
		{"2020-01-01", "2026-09-09", "range too long"},
	} {
		if _, err := svc.DailyUsage(c.from, c.to, nil); err == nil {
			t.Errorf("%s: expected an error", c.why)
		}
	}
}

// A day that has not finished has no complete final hour, so closing it is
// refused before the call is made.
func TestCollectDayRefusesTodayAndFuture(t *testing.T) {
	svc := NewMeteringService(client.NewIngextClient("http://127.0.0.1:1", "", false, nil))
	today := timeNowUTCDate()
	if _, err := svc.CollectDay(today, nil); err == nil {
		t.Error("closing today must be refused")
	}
	if _, err := svc.CollectDay("2099-01-01", nil); err == nil {
		t.Error("closing a future day must be refused")
	}
}

func TestGridUsageGoesToTheGridPrefix(t *testing.T) {
	var path string
	var kargs map[string]interface{}
	resp := model.GridUsageResponse{
		SchemaVersion: "ingext-metering-grid/v1", Grid: "provider1",
		Accounts: []*model.GridUsageAccount{
			{Account: "a", StoreAvailable: true},
			{Account: "b", Error: "timed out after 30s", State: "unavailable"},
		},
		Summary: model.GridUsageSummary{Requested: 2, Succeeded: 1, Failed: 1},
	}
	svc := newMeteringServiceForTest(t, resp, &path, &kargs)

	got, err := svc.GridUsage("2026-09-08", "2026-09-09", []string{"a", " b ", ""}, nil)
	if err != nil {
		t.Fatalf("GridUsage: %v", err)
	}
	if path != "/api/grid/grid_metering_daily_list" {
		t.Fatalf("path = %s (a provider-level call must not go to /api/ds)", path)
	}
	accts, _ := kargs["accounts"].([]interface{})
	if len(accts) != 2 || accts[0] != "a" || accts[1] != "b" {
		t.Fatalf("accounts = %+v (blank entries dropped, whitespace trimmed)", accts)
	}
	// A partially failed fan-out must announce itself.
	if got.Complete() {
		t.Fatal("a response with failed tenants must not report Complete")
	}
}

func TestGridUsageCompleteWhenNothingFailed(t *testing.T) {
	resp := model.GridUsageResponse{Summary: model.GridUsageSummary{Requested: 3, Succeeded: 3}}
	svc := newMeteringServiceForTest(t, resp, nil, nil)
	got, err := svc.GridUsage("2026-09-08", "2026-09-09", nil, nil)
	if err != nil {
		t.Fatalf("GridUsage: %v", err)
	}
	if !got.Complete() {
		t.Fatal("a fan-out with no failures must report Complete")
	}
}

// dayIndex (YYYYMMDD) is how operators address a day -- it is what appears in
// index names, dump paths and the lake's @dayIndex -- so both forms are
// accepted everywhere a date is taken.
func TestNormalizeDateAcceptsBothForms(t *testing.T) {
	for _, c := range []struct{ in, want string }{
		{"2026-09-09", "2026-09-09"},
		{"20260909", "2026-09-09"},
		{"  20260909  ", "2026-09-09"},
	} {
		got, err := NormalizeDate(c.in)
		if err != nil || got != c.want {
			t.Errorf("NormalizeDate(%q) = %q, %v; want %q", c.in, got, err, c.want)
		}
	}
	// Parsed, not pattern-matched: an impossible date must not become a range
	// the server answers with a confusingly empty ledger.
	for _, bad := range []string{"2026-02-30", "20260230", "2026-13-01", "20261301", "", "26-09-09", "2026/09/09"} {
		if _, err := NormalizeDate(bad); err == nil {
			t.Errorf("NormalizeDate(%q) should have failed", bad)
		}
	}
}

func TestDayIndexRendering(t *testing.T) {
	if got := DayIndex("2026-09-09"); got != "20260909" {
		t.Fatalf("DayIndex = %q", got)
	}
	d := &model.UsageDay{BillingDate: "2026-09-09"}
	if got := d.DayIndex(); got != "20260909" {
		t.Fatalf("UsageDay.DayIndex = %q", got)
	}
	var nilDay *model.UsageDay
	if got := nilDay.DayIndex(); got != "" {
		t.Fatalf("nil day DayIndex = %q", got)
	}
}

func TestMonthRange(t *testing.T) {
	// A month wholly in the past expands to all of it, both spellings.
	for _, m := range []string{"2020-02", "202002"} {
		from, to, partial, err := MonthRange(m)
		if err != nil {
			t.Fatalf("MonthRange(%q): %v", m, err)
		}
		if from != "2020-02-01" || to != "2020-02-29" { // leap year, checked deliberately
			t.Fatalf("MonthRange(%q) = %s..%s", m, from, to)
		}
		if partial {
			t.Fatalf("MonthRange(%q) reported partial for a finished month", m)
		}
	}

	// The current month is clamped to yesterday and says so -- today has not
	// finished and is never billable.
	now := time.Now().UTC()
	from, to, partial, err := MonthRange(now.Format("2006-01"))
	if err == nil {
		if !partial && now.Day() > 1 {
			t.Fatalf("current month %s..%s should report partial", from, to)
		}
		yesterday := now.AddDate(0, 0, -1).Format("2006-01-02")
		if partial && to != yesterday {
			t.Fatalf("current month ends %s, want yesterday %s", to, yesterday)
		}
	}

	// A future month has no finished day at all, which is an error rather than
	// an empty range.
	if _, _, _, err := MonthRange("2099-01"); err == nil {
		t.Fatal("a future month should be refused")
	}
	for _, bad := range []string{"2026-13", "202613", "2026", "2026-09-01", ""} {
		if _, _, _, err := MonthRange(bad); err == nil {
			t.Errorf("MonthRange(%q) should have failed", bad)
		}
	}
}

// A range may mix the two forms: pasting a dayIndex from an index name next to
// a hand-typed date should just work.
func TestRangeAcceptsMixedForms(t *testing.T) {
	var kargs map[string]interface{}
	svc := newMeteringServiceForTest(t, model.UsageResponse{StoreAvailable: true}, nil, &kargs)
	if _, err := svc.DailyUsage("20260901", "2026-09-09", nil); err != nil {
		t.Fatalf("mixed forms: %v", err)
	}
	if kargs["from"] != "2026-09-01" || kargs["to"] != "2026-09-09" {
		t.Fatalf("kargs = %+v; both ends must reach the wire normalised", kargs)
	}
}

func TestCollectDayAcceptsDayIndex(t *testing.T) {
	var kargs map[string]interface{}
	svc := newMeteringServiceForTest(t, model.CollectDayResponse{Closed: true}, nil, &kargs)
	if _, err := svc.CollectDay("20200101", nil); err != nil {
		t.Fatalf("CollectDay with dayIndex: %v", err)
	}
	if kargs["date"] != "2020-01-01" {
		t.Fatalf("date = %v, want the normalised form", kargs["date"])
	}
}
