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
