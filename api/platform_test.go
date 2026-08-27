package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/SecurityDo/ingext_api/client"
	fsb "github.com/SecurityDo/ingext_api/fsb"
)

// capturedProcessorTest is the platform_processor_test kargs as it arrives on
// the wire. source is a JSON string, not an object.
type capturedProcessorTest struct {
	Name   string `json:"name"`
	Script string `json:"script"`
	Source string `json:"source"`
	Type   string `json:"type"`
	Tenant string `json:"tenant"`
}

func newProcessorTestRecorder(t *testing.T, response FPLProcessorTestResult) (*PlatformService, *capturedProcessorTest) {
	t.Helper()
	captured := &capturedProcessorTest{}

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/ds/platform_processor_test" {
			t.Fatalf("unexpected path %s", r.URL.Path)
		}
		var req fsb.CallRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatalf("failed to decode request: %v", err)
		}
		if req.Kargs == nil {
			t.Fatalf("missing kargs in call request")
		}
		if err := json.Unmarshal(req.Kargs.GetBytes(), captured); err != nil {
			t.Fatalf("failed to decode kargs: %v", err)
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

	return NewPlatformService(client.NewIngextClient(ts.URL, "", false, nil)), captured
}

// TestPlatformService_TestProcessorObject locks in the request the console
// sends: the document envelope goes out as a JSON string in source, and type
// names the runtime that reads it.
func TestPlatformService_TestProcessorObject(t *testing.T) {
	svc, captured := newProcessorTestRecorder(t, FPLProcessorTestResult{
		NewContent: `{"obj":{"@type":"event"},"props":{},"size":21,"source":""}`,
		Status:     "pass",
	})

	script := "function main({obj, size}) { return { status: \"pass\" } }"
	resp, err := svc.TestProcessorObject(script, json.RawMessage("{\n  \"@cloudtrail\": {\n    \"eventName\": \"DescribeInstanceStatus\"\n  }\n}"))
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if resp.Status != "pass" {
		t.Fatalf("unexpected status %q", resp.Status)
	}

	if captured.Script != script {
		t.Fatalf("script not sent verbatim: %q", captured.Script)
	}
	if captured.Type != "fpl_processor" {
		t.Fatalf("type must name the runtime, got %q", captured.Type)
	}
	if captured.Name != "" {
		t.Fatalf("name is not looked up by this endpoint and should stay off the wire, got %q", captured.Name)
	}

	// source is a string holding the envelope, not a nested object.
	var doc FPLTestDocument
	if err := json.Unmarshal([]byte(captured.Source), &doc); err != nil {
		t.Fatalf("source is not a JSON string holding the envelope: %v", err)
	}
	if got := string(doc.Obj); got != `{"@cloudtrail":{"eventName":"DescribeInstanceStatus"}}` {
		t.Fatalf("obj not compacted into the envelope: %s", got)
	}
	// size is the compact length of obj, which is what main({obj, size}) reads.
	if doc.Size != len(doc.Obj) {
		t.Fatalf("size %d does not match the compact obj length %d", doc.Size, len(doc.Obj))
	}
	if doc.Props == nil || len(doc.Props) != 0 {
		t.Fatalf("props should be sent as {} like the console does, got %v", doc.Props)
	}
	if !strings.Contains(captured.Source, `"props":{}`) {
		t.Fatalf("props marshalled as null rather than {}: %s", captured.Source)
	}
}

// TestPlatformService_TestProcessorDefaultsType covers the caller that builds
// its own source, e.g. the raw payload an fpl_receiver reads.
func TestPlatformService_TestProcessorDefaultsType(t *testing.T) {
	svc, captured := newProcessorTestRecorder(t, FPLProcessorTestResult{Status: "pass"})

	req := &FPLProcessorTestRequest{Script: "function main() {}", Source: "raw payload"}
	if _, err := svc.TestProcessor(req); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if captured.Type != "fpl_processor" {
		t.Fatalf("an unset type should go out as fpl_processor, got %q", captured.Type)
	}
	if req.Type != "" {
		t.Fatalf("the caller's request was mutated: %+v", req)
	}
	if captured.Source != "raw payload" {
		t.Fatalf("source not sent verbatim: %q", captured.Source)
	}

	req = &FPLProcessorTestRequest{Script: "function main() {}", Source: "raw payload", Type: "fpl_receiver"}
	if _, err := svc.TestProcessor(req); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if captured.Type != "fpl_receiver" {
		t.Fatalf("an explicit type was overridden: %q", captured.Type)
	}
}

func TestNewFPLTestDocument(t *testing.T) {
	doc, err := NewFPLTestDocument(json.RawMessage(`{"a": 1, "b": "x"}`))
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if string(doc.Obj) != `{"a":1,"b":"x"}` {
		t.Fatalf("obj not compacted: %s", doc.Obj)
	}
	if doc.Size != 15 {
		t.Fatalf("unexpected size %d", doc.Size)
	}
	source, err := doc.Encode()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if source != `{"obj":{"a":1,"b":"x"},"props":{},"size":15,"source":""}` {
		t.Fatalf("unexpected envelope: %s", source)
	}

	// The script destructures the document as main({obj, size}), so anything
	// that is not a JSON object is rejected before the call.
	for _, bad := range []string{``, `   `, `[{"a":1}]`, `"a string"`, `{"a": }`} {
		if _, err := NewFPLTestDocument(json.RawMessage(bad)); err == nil {
			t.Fatalf("expected an error for %q", bad)
		}
	}
}
