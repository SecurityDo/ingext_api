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

// newFPLServiceForTest serves body (a raw platform response) for function and
// records its kargs into got.
func newFPLServiceForTest(t *testing.T, function, body string, got *map[string]interface{}) *FPLService {
	t.Helper()
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/ds/"+function {
			t.Errorf("unexpected path %s", r.URL.Path)
		}
		var req fsb.CallRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatalf("failed to decode request: %v", err)
		}
		if err := json.Unmarshal(req.Kargs.GetBytes(), got); err != nil {
			t.Fatalf("failed to decode kargs: %v", err)
		}
		w.Header().Set("Content-Type", "application/json")
		payload := fsb.CallResponse{Verdict: "OK", Response: fsb.NewJNodeByte([]byte(body))}
		if err := json.NewEncoder(w).Encode(payload); err != nil {
			t.Fatalf("failed to write response: %v", err)
		}
	}))
	t.Cleanup(ts.Close)
	return NewFPLService(client.NewIngextClient(ts.URL, "", false, nil))
}

// Shapes as titan returned them (platform 9.0.12), trimmed.
const archiveEntryJSON = `{"id":1,"account":"titan","reportName":"System_Event_Ingress_Report","taskID":76165,` +
	`"state":"completed","rangeFrom":"2026-10-22T00:00:00Z","rangeTo":"2026-10-23T00:00:00Z",` +
	`"createdOn":"2026-10-23T04:46:03.095986Z","completedOn":"2026-10-23T04:46:03.901712Z","took":69,"cost":0,` +
	`"total":52,"resultBytes":1633,"resultPath":"System_Event_Ingress_Report/2026/10/22/76165.json.gz",` +
	`"env":{"timezone":"UTC","scriptHash":"1ebc100eac76fb64","processorName":"System_Event_Ingress_Report",` +
	`"arguments":[{"name":"from","type":"string","value":"2026-10-22T00:00:00Z"}],"reportConfig":null}}`

func TestFPLService_ListReportArchive(t *testing.T) {
	var got map[string]interface{}
	body := `{"backfill":null,"coverage":{"firstDay":"2026-10-20","lastDay":"2026-10-22","expected":3,"completed":1,` +
		`"aborted":0,"running":0,"missing":2,"missingDays":["2026-10-20","2026-10-21"],"scriptHashes":["1ebc100eac76fb64"]},` +
		`"enabled":true,"entries":[` + archiveEntryJSON + `],"name":"System_Event_Ingress_Report","timezone":"UTC"}`
	svc := newFPLServiceForTest(t, "list_fpl_report_archive", body, &got)

	resp, err := svc.ListReportArchive(&model.ReportArchiveRequest{Name: "System_Event_Ingress_Report", From: "2026-10-20"})
	if err != nil {
		t.Fatal(err)
	}
	if got["name"] != "System_Event_Ingress_Report" || got["from"] != "2026-10-20" {
		t.Fatalf("unexpected kargs %v", got)
	}
	// omitempty: the platform reads absent fields as defaults.
	if _, ok := got["to"]; ok {
		t.Fatalf("empty to should not be sent: %v", got)
	}
	if resp.Coverage == nil || resp.Coverage.Missing != 2 || resp.Coverage.MissingDays[1] != "2026-10-21" {
		t.Fatalf("coverage: %+v", resp.Coverage)
	}
	if resp.Backfill != nil {
		t.Fatalf("a null backfill should decode to nil, got %+v", resp.Backfill)
	}
	if len(resp.Entries) != 1 || resp.Entries[0].Total != 52 || resp.Entries[0].Env.ScriptHash != "1ebc100eac76fb64" {
		t.Fatalf("entries: %+v", resp.Entries)
	}
}

func TestFPLService_ListArchivedReports(t *testing.T) {
	var got map[string]interface{}
	body := `{"reports":[{"name":"System_Event_Ingress_Report","enabled":true,"timezone":"UTC",` +
		`"firstDay":"2026-10-22","lastDay":"2026-10-22","days":1,"completed":1}]}`
	svc := newFPLServiceForTest(t, "list_fpl_report_archive", body, &got)

	resp, err := svc.ListReportArchive(&model.ReportArchiveRequest{})
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 0 {
		t.Fatalf("expected empty kargs, got %v", got)
	}
	if len(resp.Reports) != 1 || resp.Reports[0].Days != 1 || resp.Name != "" {
		t.Fatalf("reports: %+v", resp)
	}
}

func TestFPLService_DownloadReportArchive(t *testing.T) {
	var got map[string]interface{}
	body := `{"name":"System_Event_Ingress_Report","timezone":"UTC","coverage":{"firstDay":"2026-10-22","lastDay":"2026-10-22",` +
		`"expected":1,"completed":1,"aborted":0,"running":0,"missing":0},"days":[{"entry":` + archiveEntryJSON + `,` +
		`"result":{"objects":[{"name":"totalStats","table":{"rowCount":1,"columns":[{"name":"totalIngressSize","isVariable":false}],` +
		`"rows":[{"totalEventTypes":6,"totalIngressSize":4658507283}],` +
		// fplRows carry numeric value types; a string Type failed to decode every real result.
		`"fplRows":[{"key":{"type":1},"values":[{"v":{"type":3,"float":4658507283}},{"v":{"type":2,"int":6}}]}]}}],"error":""}}]}`
	svc := newFPLServiceForTest(t, "download_fpl_report_archive", body, &got)

	resp, err := svc.DownloadReportArchive(&model.ReportArchiveRequest{
		Name: "System_Event_Ingress_Report", Day: "2026-10-22", Objects: []string{"totalStats"}, Limit: 5})
	if err != nil {
		t.Fatal(err)
	}
	if got["day"] != "2026-10-22" || got["limit"] != float64(5) || got["objects"].([]interface{})[0] != "totalStats" {
		t.Fatalf("unexpected kargs %v", got)
	}
	if len(resp.Days) != 1 || resp.Days[0].Result == nil || resp.Days[0].Result.Objects[0].Table.Rows[0]["totalIngressSize"] != float64(4658507283) {
		t.Fatalf("days: %+v", resp.Days)
	}
}

func TestFPLService_BackfillReportArchive(t *testing.T) {
	var got map[string]interface{}
	body := `{"name":"System_Event_Ingress_Report","queued":["2026-09-22","2026-09-23"],"skipped":[],"timezone":"UTC"}`
	svc := newFPLServiceForTest(t, "backfill_fpl_report_archive", body, &got)

	resp, err := svc.BackfillReportArchive(&model.ReportArchiveRequest{
		Name: "System_Event_Ingress_Report", From: "2026-09-22", To: "2026-09-23", Force: true})
	if err != nil {
		t.Fatal(err)
	}
	if got["force"] != true || got["to"] != "2026-09-23" {
		t.Fatalf("unexpected kargs %v", got)
	}
	if len(resp.Queued) != 2 || len(resp.Skipped) != 0 {
		t.Fatalf("resp: %+v", resp)
	}
}
