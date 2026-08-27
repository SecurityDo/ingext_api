package api

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	ingextAPI "github.com/SecurityDo/ingext_api/api"
	"github.com/SecurityDo/ingext_api/client"
	fsb "github.com/SecurityDo/ingext_api/fsb"
	"github.com/SecurityDo/ingext_api/model"
)

// newProcessorDAORecorder serves the stored entry to a "get" and records the
// entry an "update" sends back.
func newProcessorDAORecorder(t *testing.T, stored *model.FPLScript) (*Client, *[]ingextAPI.GenericDAORequest[model.FPLScript]) {
	t.Helper()
	var calls []ingextAPI.GenericDAORequest[model.FPLScript]

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/ds/platform_processor_dao" {
			t.Fatalf("unexpected path %s", r.URL.Path)
		}
		var req fsb.CallRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatalf("failed to decode request: %v", err)
		}
		var dao ingextAPI.GenericDAORequest[model.FPLScript]
		if err := json.Unmarshal(req.Kargs.GetBytes(), &dao); err != nil {
			t.Fatalf("failed to decode kargs: %v", err)
		}
		calls = append(calls, dao)

		var body []byte
		var err error
		if dao.Action == "get" {
			body, err = json.Marshal(ingextAPI.ProcessorEntryResponse{Entry: stored})
		} else {
			body, err = json.Marshal(struct{}{})
		}
		if err != nil {
			t.Fatalf("failed to marshal response: %v", err)
		}
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(fsb.CallResponse{Verdict: "OK", Response: fsb.NewJNodeByte(body)}); err != nil {
			t.Fatalf("failed to write response: %v", err)
		}
	}))
	t.Cleanup(ts.Close)

	c := NewClient(slog.New(slog.NewTextHandler(nopWriter{}, nil)))
	c.ingextClient = client.NewIngextClient(ts.URL, "", false, nil)
	return c, &calls
}

type nopWriter struct{}

func (nopWriter) Write(p []byte) (int, error) { return len(p), nil }

// TestClientUpdateProcessor covers the reason update reads before it writes: a
// processor carries an id, a group and the repository it was imported from, and
// the DAO update stores the entry it is handed.
func TestClientUpdateProcessor(t *testing.T) {
	stored := &model.FPLScript{
		ID:          720042,
		Repository:  "Fluency",
		Group:       "Mimecast",
		GitPath:     "fplProcessors/code/Mimecast_Adjustments.js",
		Name:        "Mimecast_Adjustments",
		Type:        "fpl_processor",
		Description: "Standard: Mimecast integration event adjustments",
		ScriptText:  "function main() { /* old */ }",
	}
	c, calls := newProcessorDAORecorder(t, stored)

	if err := c.UpdateProcessor("Mimecast_Adjustments", "function main() { /* new */ }", "", ""); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(*calls) != 2 || (*calls)[0].Action != "get" || (*calls)[1].Action != "update" {
		t.Fatalf("expected a get then an update, got %+v", *calls)
	}
	if id := (*calls)[0].Args.Id; id != "Mimecast_Adjustments" {
		t.Fatalf("the entry is read back by name, got %q", id)
	}
	sent := (*calls)[1].Args.Entry
	if sent == nil {
		t.Fatalf("update sent no entry")
	}
	if sent.ScriptText != "function main() { /* new */ }" {
		t.Fatalf("script not replaced: %q", sent.ScriptText)
	}
	// Everything the caller did not change survives.
	if sent.ID != 720042 || sent.Repository != "Fluency" || sent.Group != "Mimecast" || sent.GitPath != stored.GitPath {
		t.Fatalf("update dropped fields the caller never touched: %+v", sent)
	}
	// An empty type or description keeps the stored one.
	if sent.Type != "fpl_processor" || sent.Description != stored.Description {
		t.Fatalf("empty flags overwrote the stored values: %+v", sent)
	}

	if err := c.UpdateProcessor("Mimecast_Adjustments", "x", "fpl_report", "new description"); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	sent = (*calls)[3].Args.Entry
	if sent.Type != "fpl_report" || sent.Description != "new description" {
		t.Fatalf("explicit type or description not applied: %+v", sent)
	}
}

// TestClientUpdateProcessorMissing keeps update from silently becoming an add.
func TestClientUpdateProcessorMissing(t *testing.T) {
	c, calls := newProcessorDAORecorder(t, nil)

	err := c.UpdateProcessor("nosuchprocessor", "function main() {}", "", "")
	if err == nil {
		t.Fatalf("expected an error for a processor that does not exist")
	}
	if len(*calls) != 1 {
		t.Fatalf("nothing should be written when the entry is missing, got %+v", *calls)
	}
}
