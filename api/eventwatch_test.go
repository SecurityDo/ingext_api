package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/SecurityDo/ingext_api/client"
	fsb "github.com/SecurityDo/ingext_api/fsb"
	"github.com/SecurityDo/ingext_api/model"
)

// capturedFilterCall is the behavior_filter_dao kargs as it arrives on the wire.
type capturedFilterCall struct {
	Action string `json:"action"`
	Args   *struct {
		Id    string                      `json:"id"`
		Entry *model.BehaviorEventFilterT `json:"entry"`
		Flag  bool                        `json:"flag"`
	} `json:"args"`
}

// newEventWatchServiceForTest serves a fixed response and records the kargs of
// each behavior_filter_dao call into got.
func newEventWatchServiceForTest(t *testing.T, response interface{}, got *capturedFilterCall) *EventWatchService {
	t.Helper()
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/ds/behavior_filter_dao" {
			t.Errorf("unexpected path %s", r.URL.Path)
		}
		var req fsb.CallRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatalf("failed to decode request: %v", err)
		}
		if req.Kargs == nil {
			t.Fatalf("missing kargs in call request")
		}
		if err := json.Unmarshal(req.Kargs.GetBytes(), got); err != nil {
			t.Fatalf("failed to decode kargs: %v", err)
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

	return NewEventWatchService(client.NewIngextClient(ts.URL, "", false, nil))
}

func TestEventWatchService_ListBehaviorFilter(t *testing.T) {
	var got capturedFilterCall
	resp := BehaviorFilterListResponse{
		Entries: []*model.BehaviorEventFilterT{
			{BehaviorRule: "AWS Root Login", Name: "break-glass", Action: model.FilterActionDiscard},
			{BehaviorRule: "*", Name: "scanner", Action: model.FilterActionSkipSummary},
		},
	}
	svc := newEventWatchServiceForTest(t, resp, &got)

	entries, err := svc.ListBehaviorFilter()
	if err != nil {
		t.Fatalf("ListBehaviorFilter returned error: %v", err)
	}
	if got.Action != "list" {
		t.Fatalf("expected action list, got %q", got.Action)
	}
	// "list" takes no args at all; the server rejects an args field it cannot use.
	if got.Args != nil {
		t.Fatalf("expected no args on list, got %+v", got.Args)
	}
	if len(entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(entries))
	}
	if entries[1].BehaviorRule != "*" || entries[1].Name != "scanner" {
		t.Fatalf("unexpected second entry: %+v", entries[1])
	}
}

func TestEventWatchService_GetBehaviorFilter(t *testing.T) {
	var got capturedFilterCall
	resp := BehaviorFilterGetResponse{
		Entry: &model.BehaviorEventFilterT{BehaviorRule: "AWS Root Login", Name: "break-glass"},
	}
	svc := newEventWatchServiceForTest(t, resp, &got)

	entry, err := svc.GetBehaviorFilter("AWS Root Login", "break-glass")
	if err != nil {
		t.Fatalf("GetBehaviorFilter returned error: %v", err)
	}
	if got.Action != "get" {
		t.Fatalf("expected action get, got %q", got.Action)
	}
	if got.Args == nil {
		t.Fatalf("missing args on get")
	}
	// The rule and the name are joined by a slash into the single id arg; the
	// server splits them back apart to rebuild the etcd key.
	if got.Args.Id != "AWS Root Login/break-glass" {
		t.Fatalf("expected the composite id, got %q", got.Args.Id)
	}
	if entry == nil || entry.Name != "break-glass" {
		t.Fatalf("unexpected entry: %+v", entry)
	}
}

func TestEventWatchService_AddBehaviorFilter(t *testing.T) {
	var got capturedFilterCall
	svc := newEventWatchServiceForTest(t, map[string]interface{}{}, &got)

	entry := &model.BehaviorEventFilterT{
		BehaviorRule: "AWS Root Login",
		Name:         "break-glass",
		Action:       model.FilterActionDiscard,
		Filters: []*model.BehaviorEventFilterEntryT{
			{Field: "key", Values: []string{"root"}, MatchType: model.MatchTypeString},
		},
	}
	if err := svc.AddBehaviorFilter(entry); err != nil {
		t.Fatalf("AddBehaviorFilter returned error: %v", err)
	}
	if got.Action != "add" {
		t.Fatalf("expected action add, got %q", got.Action)
	}
	if got.Args == nil || got.Args.Entry == nil {
		t.Fatalf("missing entry on add: %+v", got.Args)
	}
	// add identifies the filter through the entry alone; id stays empty.
	if got.Args.Id != "" {
		t.Fatalf("expected empty id on add, got %q", got.Args.Id)
	}
	if len(got.Args.Entry.Filters) != 1 || got.Args.Entry.Filters[0].Field != "key" {
		t.Fatalf("filters did not round-trip: %+v", got.Args.Entry.Filters)
	}
}

func TestEventWatchService_ToggleAndDeleteBehaviorFilter(t *testing.T) {
	for _, tc := range []struct {
		action string
		call   func(svc *EventWatchService) error
	}{
		{"toggle", func(svc *EventWatchService) error { return svc.ToggleBehaviorFilter("*", "scanner") }},
		{"delete", func(svc *EventWatchService) error { return svc.DeleteBehaviorFilter("*", "scanner") }},
	} {
		t.Run(tc.action, func(t *testing.T) {
			var got capturedFilterCall
			svc := newEventWatchServiceForTest(t, map[string]interface{}{}, &got)

			if err := tc.call(svc); err != nil {
				t.Fatalf("%s returned error: %v", tc.action, err)
			}
			if got.Action != tc.action {
				t.Fatalf("expected action %s, got %q", tc.action, got.Action)
			}
			if got.Args == nil || got.Args.Id != "*/scanner" {
				t.Fatalf("unexpected args: %+v", got.Args)
			}
		})
	}
}

func TestEventWatchService_BehaviorFilterIDGuards(t *testing.T) {
	// The server splits the id on the first slash, so these inputs cannot round
	// trip into a key. Catch them here rather than spending a request on an
	// unhelpful "invalid filter name" from the backend.
	svc := NewEventWatchService(client.NewIngextClient("http://127.0.0.1:1", "", false, nil))
	for _, tc := range []struct{ rule, name string }{
		{"", "scanner"},
		{"*", ""},
		{"a/b", "scanner"},
	} {
		if _, err := svc.GetBehaviorFilter(tc.rule, tc.name); err == nil {
			t.Fatalf("expected error for rule %q name %q", tc.rule, tc.name)
		}
		if err := svc.DeleteBehaviorFilter(tc.rule, tc.name); err == nil {
			t.Fatalf("expected error for rule %q name %q", tc.rule, tc.name)
		}
	}

	// A slash inside the filter name is fine: the server keeps everything after
	// the first separator.
	if got := BehaviorFilterID("rule", "a/b"); got != "rule/a/b" {
		t.Fatalf("unexpected composite id %q", got)
	}
}

func TestEventWatchService_BehaviorFilterGridAccount(t *testing.T) {
	var got capturedFilterCall
	var gotQuery string
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotQuery = r.URL.Query().Get("gridaccount")
		var req fsb.CallRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatalf("failed to decode request: %v", err)
		}
		if err := json.Unmarshal(req.Kargs.GetBytes(), &got); err != nil {
			t.Fatalf("failed to decode kargs: %v", err)
		}
		body, _ := json.Marshal(BehaviorFilterListResponse{})
		payload := fsb.CallResponse{Verdict: "OK", Response: fsb.NewJNodeByte(body)}
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(payload); err != nil {
			t.Fatalf("failed to write response: %v", err)
		}
	}))
	t.Cleanup(ts.Close)

	flClient := client.NewIngextClient(ts.URL, "", false, nil)
	flClient.SetGridAccount("tenant1")
	if _, err := NewEventWatchService(flClient).ListBehaviorFilter(); err != nil {
		t.Fatalf("ListBehaviorFilter returned error: %v", err)
	}
	// The tenant is routed through the query string only; it never enters kargs.
	if gotQuery != "tenant1" {
		t.Fatalf("expected gridaccount=tenant1, got %q", gotQuery)
	}
	if got.Args != nil {
		t.Fatalf("expected no args, got %+v", got.Args)
	}
}

// capturedRuleTest is the eventwatch_rule_test kargs as it arrives on the wire.
type capturedRuleTest struct {
	Bucket *model.EventWatchBucket `json:"bucket"`
	Input  json.RawMessage         `json:"input"`
}

// newRuleTestRecorder serves a fixed response per endpoint path and records the
// eventwatch_rule_test kargs.
func newRuleTestRecorder(t *testing.T, ruleTest interface{}, stored *model.EventWatchBucket) (*EventWatchService, *capturedRuleTest, *[]string) {
	t.Helper()
	captured := &capturedRuleTest{}
	var paths []string

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		paths = append(paths, r.URL.Path)
		var req fsb.CallRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatalf("failed to decode request: %v", err)
		}

		var response interface{}
		switch r.URL.Path {
		case "/api/ds/eventwatch_rule_test":
			if err := json.Unmarshal(req.Kargs.GetBytes(), captured); err != nil {
				t.Fatalf("failed to decode kargs: %v", err)
			}
			response = ruleTest
		case "/api/ds/eventwatch_bucket_dao":
			response = EventWatchRuleGetResponse{Entry: stored}
		default:
			t.Fatalf("unexpected path %s", r.URL.Path)
		}

		body, err := json.Marshal(response)
		if err != nil {
			t.Fatalf("failed to marshal response: %v", err)
		}
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(fsb.CallResponse{Verdict: "OK", Response: fsb.NewJNodeByte(body)}); err != nil {
			t.Fatalf("failed to write response: %v", err)
		}
	}))
	t.Cleanup(ts.Close)

	return NewEventWatchService(client.NewIngextClient(ts.URL, "", false, nil)), captured, &paths
}

func TestEventWatchService_TestRuleEvent(t *testing.T) {
	hit := EventWatchRuleTestResult{
		Hit: true,
		Signals: []string{
			`{"signal":"behavior:AD_Event_Log_Cleared","ts":1787853120,"count":1,"key":"dc01","valueMap":{"@fields.Channel":"System"}}`,
		},
		BehaviorEvent: &model.BehaviorEvent{
			Timestamp:    1787853120000,
			Key:          "dc01",
			BehaviorRule: "AD_Event_Log_Cleared",
			Behavior:     "security alert",
		},
	}
	svc, captured, _ := newRuleTestRecorder(t, hit, nil)

	rule := &model.EventWatchBucket{Name: "AD_Event_Log_Cleared", EventType: "event"}
	event := json.RawMessage(`{"@eventType": "nxlogAD"}`)
	resp, err := svc.TestRuleEvent(rule, event)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if !resp.Hit {
		t.Fatalf("unexpected result %+v", resp)
	}
	if captured.Bucket == nil || captured.Bucket.Name != "AD_Event_Log_Cleared" {
		t.Fatalf("the whole rule goes out under bucket, got %+v", captured.Bucket)
	}
	// The event is a JSON object under input, not a string holding one: the
	// endpoint accepts a string and then ignores it.
	var got map[string]interface{}
	if err := json.Unmarshal(captured.Input, &got); err != nil {
		t.Fatalf("input is not a JSON object: %v", err)
	}
	if got["@eventType"] != "nxlogAD" {
		t.Fatalf("event not sent verbatim: %v", got)
	}

	signals, err := resp.DecodeSignals()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(signals) != 1 {
		t.Fatalf("unexpected signals %+v", signals)
	}
	if signals[0].Signal != "behavior:AD_Event_Log_Cleared" || signals[0].Key != "dc01" || signals[0].Count != 1 {
		t.Fatalf("signal not decoded: %+v", signals[0])
	}
	if signals[0].ValueMap["@fields.Channel"] != "System" {
		t.Fatalf("valueMap not decoded: %+v", signals[0].ValueMap)
	}
}

// TestEventWatchService_TestRuleGuards covers the two inputs the endpoint
// answers with a panic rather than an error, and the one it silently ignores.
func TestEventWatchService_TestRuleGuards(t *testing.T) {
	svc, _, paths := newRuleTestRecorder(t, EventWatchRuleTestResult{}, nil)

	if _, err := svc.TestRule(&EventWatchRuleTestRequest{}); err == nil {
		t.Fatalf("expected an error for a missing rule")
	}
	if _, err := svc.TestRuleEvent(&model.EventWatchBucket{EventType: "event"}, nil); err == nil {
		t.Fatalf("expected an error for a rule with no name")
	}
	for _, bad := range []string{`[{"a":1}]`, `"an event"`, `   `} {
		_, err := svc.TestRuleEvent(&model.EventWatchBucket{Name: "r"}, json.RawMessage(bad))
		if err == nil {
			t.Fatalf("expected an error for input %q", bad)
		}
	}
	if len(*paths) != 0 {
		t.Fatalf("nothing should reach the endpoint, got %v", *paths)
	}

	// A rule with no event is allowed: the endpoint answers hit false.
	if _, err := svc.TestRuleEvent(&model.EventWatchBucket{Name: "r"}, nil); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

// TestEventWatchService_TestDeployedRule reads the rule before testing it.
func TestEventWatchService_TestDeployedRule(t *testing.T) {
	stored := &model.EventWatchBucket{ID: 101861, Name: "AD_Event_Log_Cleared", Group: "AD"}
	svc, captured, paths := newRuleTestRecorder(t, EventWatchRuleTestResult{Hit: true}, stored)

	if _, err := svc.TestDeployedRule("AD_Event_Log_Cleared", json.RawMessage(`{"a":1}`)); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(*paths) != 2 || (*paths)[0] != "/api/ds/eventwatch_bucket_dao" || (*paths)[1] != "/api/ds/eventwatch_rule_test" {
		t.Fatalf("expected a dao get then a rule test, got %v", *paths)
	}
	if captured.Bucket == nil || captured.Bucket.ID != 101861 {
		t.Fatalf("the stored rule should be the one tested, got %+v", captured.Bucket)
	}

	// An empty name would return whichever rule the DAO lists first.
	if _, err := svc.TestDeployedRule("  ", json.RawMessage(`{"a":1}`)); err == nil {
		t.Fatalf("expected an error for an empty name")
	}
}
