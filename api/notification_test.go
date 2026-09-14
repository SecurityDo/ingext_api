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

// newNotificationRecorder serves one canned response per function and records
// the kargs each call sent.
func newNotificationRecorder(t *testing.T, responses map[string]interface{}) (*NotificationService, *[]json.RawMessage) {
	t.Helper()
	var kargs []json.RawMessage

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req fsb.CallRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatalf("failed to decode request: %v", err)
		}
		if req.Kargs == nil {
			kargs = append(kargs, nil)
		} else {
			kargs = append(kargs, json.RawMessage(req.Kargs.GetBytes()))
		}

		response, ok := responses[req.Function]
		if !ok {
			t.Fatalf("unexpected function %s", req.Function)
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

	return NewNotificationService(client.NewIngextClient(ts.URL, "", false, nil)), &kargs
}

// newNotificationErrorService serves an ERROR verdict, the way the dao reports a
// name it does not hold or one it already has.
func newNotificationErrorService(t *testing.T, message string) *NotificationService {
	t.Helper()
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(fsb.CallResponse{Verdict: "ERROR", Error: message}); err != nil {
			t.Fatalf("failed to write response: %v", err)
		}
	}))
	t.Cleanup(ts.Close)
	return NewNotificationService(client.NewIngextClient(ts.URL, "", false, nil))
}

func decodeDAO(t *testing.T, raw json.RawMessage) GenericDAORequest[model.EndpointConfig] {
	t.Helper()
	var dao GenericDAORequest[model.EndpointConfig]
	if err := json.Unmarshal(raw, &dao); err != nil {
		t.Fatalf("failed to decode kargs: %v", err)
	}
	return dao
}

// TestNotificationService_Get locks in that a "get" addresses the endpoint by
// name through args.id and that the reply is unwrapped from its "entry" key.
func TestNotificationService_Get(t *testing.T) {
	stored := &model.EndpointConfig{
		Name:        "ops",
		Integration: "Email",
		Action:      "Generic_Email_Action",
		Email:       &model.EndpointEmailConfig{To: []string{"a@x"}},
	}
	svc, kargs := newNotificationRecorder(t, map[string]interface{}{
		"platform_notification_endpoint_dao": GenericDaoEntryResponse[model.EndpointConfig]{Entry: stored},
	})

	entry, err := svc.Get("ops")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if entry == nil || entry.Name != "ops" || entry.Action != "Generic_Email_Action" {
		t.Fatalf("unexpected entry %+v", entry)
	}

	dao := decodeDAO(t, (*kargs)[0])
	if dao.Action != "get" {
		t.Errorf("expected action get, got %q", dao.Action)
	}
	if dao.Args == nil || dao.Args.Id != "ops" {
		t.Errorf("expected args.id ops, got %+v", dao.Args)
	}
	if dao.Args.Entry != nil {
		t.Errorf("get must not send an entry, got %+v", dao.Args.Entry)
	}
}

// TestNotificationService_GetMissing covers what the server actually does with
// a name it does not hold: an ERROR verdict, not an empty result. Code that
// treats "not found" as a nil entry never runs.
func TestNotificationService_GetMissing(t *testing.T) {
	svc := newNotificationErrorService(t, "export not found: nope")

	entry, err := svc.Get("nope")
	if err == nil {
		t.Fatalf("expected an error for a missing endpoint, got entry %+v", entry)
	}
	if entry != nil {
		t.Errorf("expected nil entry alongside the error, got %+v", entry)
	}
}

// TestNotificationService_GetNullEntry covers the defensive path: a dao that
// answers with a null entry yields a nil result rather than a zero-valued one.
func TestNotificationService_GetNullEntry(t *testing.T) {
	svc, _ := newNotificationRecorder(t, map[string]interface{}{
		"platform_notification_endpoint_dao": GenericDaoEntryResponse[model.EndpointConfig]{},
	})

	entry, err := svc.Get("nope")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if entry != nil {
		t.Fatalf("expected nil entry, got %+v", entry)
	}
}

// TestNotificationService_DeleteMissing covers delete of a name the dao does not
// hold: it fails rather than succeeding quietly.
func TestNotificationService_DeleteMissing(t *testing.T) {
	svc := newNotificationErrorService(t, "unknown export")
	if err := svc.Delete("nope"); err == nil {
		t.Fatal("expected an error deleting a missing endpoint")
	}
}

// TestNotificationService_AddDuplicate covers add refusing a name already
// stored: it is not an upsert.
func TestNotificationService_AddDuplicate(t *testing.T) {
	svc := newNotificationErrorService(t, "duplicate endpoint")
	if _, err := svc.AddEmail("ops", "Generic_Email_Action", []string{"a@x"}, nil); err == nil {
		t.Fatal("expected an error adding a duplicate endpoint")
	}
}

// TestNotificationService_AddReturnsNoID locks in that the dao answers an add
// with "{}": the id comes back empty, and that is not an error.
func TestNotificationService_AddReturnsNoID(t *testing.T) {
	svc, _ := newNotificationRecorder(t, map[string]interface{}{
		"platform_notification_endpoint_dao": struct{}{},
	})

	id, err := svc.AddEmail("ops", "Generic_Email_Action", []string{"a@x"}, nil)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if id != "" {
		t.Errorf("expected an empty id, got %q", id)
	}
}

// TestNotificationService_Update locks in that the name goes out as args.id
// alongside the whole replacement entry.
func TestNotificationService_Update(t *testing.T) {
	svc, kargs := newNotificationRecorder(t, map[string]interface{}{
		"platform_notification_endpoint_dao": struct{}{},
	})

	err := svc.Update(&model.EndpointConfig{
		Name:        "ops",
		Integration: "Email",
		Action:      "Generic_Email_Action",
		Email:       &model.EndpointEmailConfig{To: []string{"b@x"}},
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	dao := decodeDAO(t, (*kargs)[0])
	if dao.Action != "update" {
		t.Errorf("expected action update, got %q", dao.Action)
	}
	if dao.Args == nil || dao.Args.Id != "ops" {
		t.Fatalf("expected args.id ops, got %+v", dao.Args)
	}
	if dao.Args.Entry == nil || dao.Args.Entry.Email == nil || dao.Args.Entry.Email.To[0] != "b@x" {
		t.Errorf("unexpected entry %+v", dao.Args.Entry)
	}
}

func TestNotificationService_UpdateNilEntry(t *testing.T) {
	svc, _ := newNotificationRecorder(t, map[string]interface{}{})
	if err := svc.Update(nil); err == nil {
		t.Fatal("expected an error for a nil entry")
	}
}

// TestNotificationService_AddSlack locks in that a single channel is sent in
// both the scalar and the array field, and several channels only in the array.
func TestNotificationService_AddSlack(t *testing.T) {
	for _, tc := range []struct {
		name        string
		channels    []string
		wantChannel string
	}{
		{"one channel", []string{"#soc"}, "#soc"},
		{"several channels", []string{"#soc", "#ops"}, ""},
	} {
		t.Run(tc.name, func(t *testing.T) {
			svc, kargs := newNotificationRecorder(t, map[string]interface{}{
				"platform_notification_endpoint_dao": GenericDaoAddResponse{ID: "ep-1"},
			})

			id, err := svc.AddSlack("soc", "Slack_Action", "corp-slack", tc.channels)
			if err != nil {
				t.Fatalf("expected no error, got %v", err)
			}
			if id != "ep-1" {
				t.Errorf("expected id ep-1, got %q", id)
			}

			dao := decodeDAO(t, (*kargs)[0])
			if dao.Action != "add" {
				t.Errorf("expected action add, got %q", dao.Action)
			}
			entry := dao.Args.Entry
			if entry.Integration != "Slack" || entry.Slack == nil {
				t.Fatalf("unexpected entry %+v", entry)
			}
			if entry.Slack.IntegrationName != "corp-slack" {
				t.Errorf("expected integrationName corp-slack, got %q", entry.Slack.IntegrationName)
			}
			if len(entry.Slack.Channels) != len(tc.channels) {
				t.Errorf("expected %d channels, got %v", len(tc.channels), entry.Slack.Channels)
			}
			if entry.Slack.Channel != tc.wantChannel {
				t.Errorf("expected channel %q, got %q", tc.wantChannel, entry.Slack.Channel)
			}
		})
	}
}

// TestNotificationService_ListActions locks in that the actions come back under
// "actions" rather than "entries", and that the filter keeps only the ones a
// notification endpoint can actually name.
func TestNotificationService_ListActions(t *testing.T) {
	action := func(name, target, integration string) *model.FPLScript {
		return &model.FPLScript{
			Name:         name,
			Type:         "fpl_action",
			ActionConfig: &model.FplActionConfig{Target: target, Integration: integration},
		}
	}
	responses := map[string]interface{}{
		"platform_list_actions": ListActionsResponse{Actions: []*model.FPLScript{
			action("Generic_Email_Action", "Platform Notification", "Email"),
			action("Slack_Action", "Platform Notification", "Slack"),
			action("Ticket_Action", "Case Management", "Jira"),
			{Name: "No_Config", Type: "fpl_action"},
		}},
	}

	svc, kargs := newNotificationRecorder(t, responses)
	all, err := svc.ListActions("")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(all) != 2 {
		t.Fatalf("expected 2 notification actions, got %d", len(all))
	}
	// The endpoint reads no kargs; a nil payload goes out as an empty object.
	if got := string((*kargs)[0]); got != "{}" && got != "null" {
		t.Errorf("platform_list_actions takes no kargs, sent %s", got)
	}

	svc, _ = newNotificationRecorder(t, responses)
	email, err := svc.ListActions("Email")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(email) != 1 || email[0].Name != "Generic_Email_Action" {
		t.Fatalf("unexpected filtered actions %+v", email)
	}
}

// TestPlatformService_ListActions covers the unfiltered endpoint binding,
// including the actionConfig field the processor dao never fills in.
func TestPlatformService_ListActions(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/ds/platform_list_actions" {
			t.Fatalf("unexpected path %s", r.URL.Path)
		}
		body := []byte(`{"actions":[{"id":720024,"name":"Generic_Email_Action","group":"System","type":"fpl_action","actionConfig":{"target":"Platform Notification","integration":"Email"},"description":"Default: HTML Email Action","scriptLang":"","scriptText":"function main(){}","createdOn":"2024-04-18T21:26:01Z","updatedOn":"2024-06-19T00:00:00Z"}]}`)
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(fsb.CallResponse{Verdict: "OK", Response: fsb.NewJNodeByte(body)}); err != nil {
			t.Fatalf("failed to write response: %v", err)
		}
	}))
	t.Cleanup(ts.Close)

	svc := NewPlatformService(client.NewIngextClient(ts.URL, "", false, nil))
	actions, err := svc.ListActions()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(actions) != 1 {
		t.Fatalf("expected 1 action, got %d", len(actions))
	}
	got := actions[0]
	if got.ID != 720024 || got.Name != "Generic_Email_Action" {
		t.Errorf("unexpected action %+v", got)
	}
	if got.ActionConfig == nil || got.ActionConfig.Target != "Platform Notification" || got.ActionConfig.Integration != "Email" {
		t.Errorf("actionConfig not decoded: %+v", got.ActionConfig)
	}
}
