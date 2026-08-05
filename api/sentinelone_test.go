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

func newSentinelOneServiceForTest(t *testing.T, handler http.HandlerFunc) (*SentinelOneService, *client.IngextClient) {
	t.Helper()
	ts := httptest.NewServer(handler)
	t.Cleanup(ts.Close)

	flClient := client.NewIngextClient(ts.URL, "", false, nil)
	return NewSentinelOneService(flClient), flClient
}

// investigateKargs decodes the raw kargs object off the wire. The tests below
// assert on the raw map rather than on InvestigateAlertRequest because the point
// is which keys were *emitted* — unmarshalling into the typed struct would hide
// exactly the omissions that carry meaning.
func investigateKargs(t *testing.T, r *http.Request) map[string]interface{} {
	t.Helper()
	var req fsb.CallRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		t.Fatalf("failed to decode request: %v", err)
	}
	if req.Kargs == nil {
		t.Fatalf("missing kargs in call request")
	}
	var kargs map[string]interface{}
	if err := json.Unmarshal(req.Kargs.GetBytes(), &kargs); err != nil {
		t.Fatalf("failed to decode kargs: %v", err)
	}
	return kargs
}

func writeInvestigation(t *testing.T, w http.ResponseWriter, inv *model.AlertInvestigation) {
	t.Helper()
	body, err := json.Marshal(inv)
	if err != nil {
		t.Fatalf("failed to marshal response: %v", err)
	}
	respPayload := fsb.CallResponse{Verdict: "OK", Response: fsb.NewJNodeByte(body)}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(respPayload); err != nil {
		t.Fatalf("failed to write response: %v", err)
	}
}

func TestSentinelOneService_InvestigateAlertDefaults(t *testing.T) {
	svc, _ := newSentinelOneServiceForTest(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/ds/investigate_sentinelone_alert" {
			t.Fatalf("unexpected path %s", r.URL.Path)
		}

		kargs := investigateKargs(t, r)
		if kargs["alertId"] != "alert-1" {
			t.Fatalf("unexpected alertId %v", kargs["alertId"])
		}
		// No options given: nothing should be synthesized, or the platform would
		// see explicit values where the caller meant "use your defaults".
		if _, ok := kargs["options"]; ok {
			t.Fatalf("options should be omitted when unset: %+v", kargs)
		}
		if _, ok := kargs["integrationName"]; ok {
			t.Fatalf("integrationName should be omitted when unset: %+v", kargs)
		}

		writeInvestigation(t, w, &model.AlertInvestigation{
			Alert:            &model.AlertSummary{ID: "alert-1", Severity: "high"},
			Associations:     &model.Associations{Method: model.S1MethodExternalID, Confidence: model.S1ConfidenceHigh},
			CollectionStatus: &model.CollectionStatus{Threats: model.S1StatusComplete, Activities: model.S1StatusComplete},
		})
	})

	inv, err := svc.InvestigateAlertByID("alert-1", "")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if inv.Alert == nil || inv.Alert.ID != "alert-1" {
		t.Fatalf("unexpected alert: %+v", inv.Alert)
	}
	if inv.Associations.Confidence != model.S1ConfidenceHigh {
		t.Fatalf("unexpected confidence: %+v", inv.Associations)
	}
	if inv.CollectionStatus.Threats != model.S1StatusComplete {
		t.Fatalf("unexpected collection status: %+v", inv.CollectionStatus)
	}
}

// TestSentinelOneService_InvestigateAlertTriStateFlags is the client-side half of
// the platform's mcpParity_test: an explicit false must reach the wire, and an
// omitted flag must stay off the wire so the platform default survives.
func TestSentinelOneService_InvestigateAlertTriStateFlags(t *testing.T) {
	svc, _ := newSentinelOneServiceForTest(t, func(w http.ResponseWriter, r *http.Request) {
		kargs := investigateKargs(t, r)
		if kargs["integrationName"] != "s1-prod" {
			t.Fatalf("unexpected integrationName %v", kargs["integrationName"])
		}

		options, ok := kargs["options"].(map[string]interface{})
		if !ok {
			t.Fatalf("expected an options object, got %+v", kargs["options"])
		}
		if v, ok := options["includeThreats"]; !ok || v != false {
			t.Fatalf("explicit false must survive, got %v (present=%v)", v, ok)
		}
		if v, ok := options["includeRelatedAlerts"]; !ok || v != true {
			t.Fatalf("explicit true must survive, got %v (present=%v)", v, ok)
		}
		if _, ok := options["includeActivities"]; ok {
			t.Fatalf("omitted flag must not reach the wire: %+v", options)
		}
		if _, ok := options["includeEndpoint"]; ok {
			t.Fatalf("omitted flag must not reach the wire: %+v", options)
		}
		if v, _ := options["activityWindowMinutes"].(float64); v != 120 {
			t.Fatalf("unexpected activityWindowMinutes %v", options["activityWindowMinutes"])
		}
		if v, _ := options["maxActivities"].(float64); v != 500 {
			t.Fatalf("unexpected maxActivities %v", options["maxActivities"])
		}

		writeInvestigation(t, w, &model.AlertInvestigation{
			Alert: &model.AlertSummary{ID: "alert-1"},
			CollectionStatus: &model.CollectionStatus{
				Threats:       model.S1StatusSkipped,
				Activities:    model.S1StatusComplete,
				RelatedAlerts: model.S1StatusComplete,
			},
		})
	})

	no, yes := false, true
	inv, err := svc.InvestigateAlert(&InvestigateAlertRequest{
		AlertID:         "alert-1",
		IntegrationName: "s1-prod",
		Options: &InvestigateAlertOptions{
			IncludeThreats:        &no,
			IncludeRelatedAlerts:  &yes,
			ActivityWindowMinutes: 120,
			MaxActivities:         500,
		},
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if inv.CollectionStatus.Threats != model.S1StatusSkipped {
		t.Fatalf("expected threats to be skipped: %+v", inv.CollectionStatus)
	}
}

// TestSentinelOneService_InvestigateAlertPartial covers the documented contract
// that a failed section is still a successful call.
func TestSentinelOneService_InvestigateAlertPartial(t *testing.T) {
	svc, _ := newSentinelOneServiceForTest(t, func(w http.ResponseWriter, r *http.Request) {
		writeInvestigation(t, w, &model.AlertInvestigation{
			Alert:        &model.AlertSummary{ID: "alert-1"},
			Associations: &model.Associations{Method: model.S1MethodNone, Confidence: model.S1ConfidenceNone},
			CollectionStatus: &model.CollectionStatus{
				Threats:    model.S1StatusFailed,
				Activities: model.S1StatusTruncated,
				Warnings:   []string{"threat lookup failed"},
			},
		})
	})

	inv, err := svc.InvestigateAlertByID("alert-1", "")
	if err != nil {
		t.Fatalf("a failed section must not fail the call, got %v", err)
	}
	if inv.CollectionStatus.Threats != model.S1StatusFailed {
		t.Fatalf("unexpected threat status: %+v", inv.CollectionStatus)
	}
	if len(inv.CollectionStatus.Warnings) != 1 {
		t.Fatalf("warnings lost: %+v", inv.CollectionStatus)
	}
}

func TestSentinelOneService_InvestigateAlertMissingAlertID(t *testing.T) {
	svc, _ := newSentinelOneServiceForTest(t, func(w http.ResponseWriter, r *http.Request) {
		t.Fatalf("no request should be sent without an alertId")
	})

	if _, err := svc.InvestigateAlert(&InvestigateAlertRequest{}); err == nil {
		t.Fatal("expected an error for a missing alertId")
	}
	if _, err := svc.InvestigateAlert(nil); err == nil {
		t.Fatal("expected an error for a nil request")
	}
}

// TestSentinelOneService_InvestigateAlertGridAccount locks in the routing
// decision: the tenant is named by the gridaccount query parameter, never by a
// field in kargs.
func TestSentinelOneService_InvestigateAlertGridAccount(t *testing.T) {
	svc, flClient := newSentinelOneServiceForTest(t, func(w http.ResponseWriter, r *http.Request) {
		if got := r.URL.Query().Get("gridaccount"); got != "acme" {
			t.Fatalf("unexpected gridaccount %q", got)
		}
		kargs := investigateKargs(t, r)
		if _, ok := kargs["account"]; ok {
			t.Fatalf("the account must not be spliced into kargs: %+v", kargs)
		}

		writeInvestigation(t, w, &model.AlertInvestigation{Alert: &model.AlertSummary{ID: "alert-1"}})
	})
	flClient.SetGridAccount("acme")

	if _, err := svc.InvestigateAlertByID("alert-1", ""); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}
