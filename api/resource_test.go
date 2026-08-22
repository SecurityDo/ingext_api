package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/SecurityDo/ingext_api/client"
	fsb "github.com/SecurityDo/ingext_api/fsb"
)

func newResourceServiceForTest(t *testing.T, handler http.HandlerFunc) *ResourceService {
	t.Helper()
	ts := httptest.NewServer(handler)
	t.Cleanup(ts.Close)

	flClient := client.NewIngextClient(ts.URL, "", false, nil)
	return NewResourceService(flClient)
}

func TestResourceService_DeleteResourceDump(t *testing.T) {
	svc := newResourceServiceForTest(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/ds/ingext_resource_dump_delete" {
			t.Fatalf("unexpected path %s", r.URL.Path)
		}

		var req fsb.CallRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatalf("failed to decode request: %v", err)
		}
		var kargs struct {
			Customer string `json:"customer"`
		}
		if err := json.Unmarshal(req.Kargs.GetBytes(), &kargs); err != nil {
			t.Fatalf("failed to decode kargs: %v", err)
		}
		if kargs.Customer != "office365" {
			t.Fatalf("unexpected customer %q", kargs.Customer)
		}

		body, err := json.Marshal(ResourceDumpDeleteResponse{Deleted: 3})
		if err != nil {
			t.Fatalf("failed to marshal response: %v", err)
		}
		respPayload := fsb.CallResponse{Verdict: "OK", Response: fsb.NewJNodeByte(body)}
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(respPayload); err != nil {
			t.Fatalf("failed to write response: %v", err)
		}
	})

	deleted, err := svc.DeleteResourceDump("office365")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if deleted != 3 {
		t.Fatalf("expected 3 dumps deleted, got %d", deleted)
	}
}

// TestResourceService_DeleteResourceDumpUnknownCustomer locks in that a customer
// with no dumps is a successful no-op, not an error — the count is the only
// signal that separates the two.
func TestResourceService_DeleteResourceDumpUnknownCustomer(t *testing.T) {
	svc := newResourceServiceForTest(t, func(w http.ResponseWriter, r *http.Request) {
		respPayload := fsb.CallResponse{Verdict: "OK", Response: fsb.NewJNodeByte([]byte(`{"deleted":0}`))}
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(respPayload); err != nil {
			t.Fatalf("failed to write response: %v", err)
		}
	})

	deleted, err := svc.DeleteResourceDump("never-collected")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if deleted != 0 {
		t.Fatalf("expected 0 dumps deleted, got %d", deleted)
	}
}

// TestResourceService_DeleteResourceDumpEmptyCustomer checks the empty customer
// is rejected client-side: the platform would reject it too, but a request that
// can only fail is not worth sending.
func TestResourceService_DeleteResourceDumpEmptyCustomer(t *testing.T) {
	svc := newResourceServiceForTest(t, func(w http.ResponseWriter, r *http.Request) {
		t.Fatalf("expected no request for an empty customer")
	})

	if _, err := svc.DeleteResourceDump(""); err == nil {
		t.Fatalf("expected an error for an empty customer")
	}
}
