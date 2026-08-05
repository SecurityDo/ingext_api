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

func newDatalakeServiceForTest(t *testing.T, handler http.HandlerFunc) *DatalakeService {
	t.Helper()
	ts := httptest.NewServer(handler)
	t.Cleanup(ts.Close)

	flClient := client.NewIngextClient(ts.URL, "", false, nil)
	return NewDatalakeService(flClient)
}

func TestDatalakeService_ListDataTables(t *testing.T) {
	svc := newDatalakeServiceForTest(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/ds/list_data_tables" {
			t.Fatalf("unexpected path %s", r.URL.Path)
		}

		var req fsb.CallRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatalf("failed to decode request: %v", err)
		}
		// The function takes no arguments, but kargs must still be present and
		// must be an empty object — the platform rejects a missing kargs.
		if req.Kargs == nil {
			t.Fatalf("missing kargs in call request")
		}
		var kargs map[string]interface{}
		if err := json.Unmarshal(req.Kargs.GetBytes(), &kargs); err != nil {
			t.Fatalf("failed to decode kargs: %v", err)
		}
		if len(kargs) != 0 {
			t.Fatalf("expected empty kargs, got %+v", kargs)
		}

		response := model.ListTableResponse{
			StreamTables:   []*model.DataTable{{Name: "OfficeActivity", Description: "M365 audit log"}},
			ResourceTables: []*model.DataTable{{Name: "azureUser", Description: ""}},
		}
		body, err := json.Marshal(response)
		if err != nil {
			t.Fatalf("failed to marshal response: %v", err)
		}
		respPayload := fsb.CallResponse{Verdict: "OK", Response: fsb.NewJNodeByte(body)}
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(respPayload); err != nil {
			t.Fatalf("failed to write response: %v", err)
		}
	})

	tables, err := svc.ListDataTables()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(tables.StreamTables) != 1 || tables.StreamTables[0].Name != "OfficeActivity" {
		t.Fatalf("unexpected stream tables: %+v", tables.StreamTables)
	}
	if tables.StreamTables[0].Description != "M365 audit log" {
		t.Fatalf("stream table description lost: %+v", tables.StreamTables[0])
	}
	if len(tables.ResourceTables) != 1 || tables.ResourceTables[0].Name != "azureUser" {
		t.Fatalf("unexpected resource tables: %+v", tables.ResourceTables)
	}
}

// TestDatalakeService_ListDataTablesEmpty covers the documented empty case: an
// account with nothing to query returns "{}", not empty arrays.
func TestDatalakeService_ListDataTablesEmpty(t *testing.T) {
	svc := newDatalakeServiceForTest(t, func(w http.ResponseWriter, r *http.Request) {
		respPayload := fsb.CallResponse{Verdict: "OK", Response: fsb.NewJNodeByte([]byte(`{}`))}
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(respPayload); err != nil {
			t.Fatalf("failed to write response: %v", err)
		}
	})

	tables, err := svc.ListDataTables()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(tables.StreamTables) != 0 || len(tables.ResourceTables) != 0 {
		t.Fatalf("expected an empty catalog, got %+v", tables)
	}
}
