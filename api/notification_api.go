package api

import (
	"fmt"

	"github.com/SecurityDo/ingext_api/client"
	"github.com/SecurityDo/ingext_api/model"
)

// NotificationService wraps platform_notification_endpoint_dao, the dao that
// holds the notification endpoints an FPL action delivers through.
//
// An endpoint pairs an Integration ("Email", "Slack") with the name of an
// fpl_action script in its Action field, plus the per-integration config that
// script reads as its `config` argument. PlatformService.ListActions reports
// which action names exist and which integration each one expects.
type NotificationService struct {
	client *client.IngextClient
}

func NewNotificationService(client *client.IngextClient) *NotificationService {
	return &NotificationService{client: client}
}

func (s *NotificationService) call(function string, payload interface{}, out interface{}) error {
	return ApiCall(s.client, function, payload, out)
}

const notificationEndpointDAO = "platform_notification_endpoint_dao"

// Add stores a new notification endpoint. The dao keys endpoints by name and
// rejects a name it already holds with "duplicate endpoint" -- add is not an
// upsert; use Update for that.
//
// The returned id is empty in practice: the dao answers an add with "{}" and no
// id field. It is kept because the response shape is the generic dao one and may
// start carrying an id. Address the endpoint by its name, not by this value.
func (s *NotificationService) Add(entry *model.EndpointConfig) (id string, err error) {
	if entry == nil {
		return "", fmt.Errorf("notification endpoint entry is required")
	}
	request := &GenericDAORequest[model.EndpointConfig]{
		Action: "add",
		Args: &GenericDAORequestArgs[model.EndpointConfig]{
			Entry: entry,
		},
	}
	var resp GenericDaoAddResponse
	if err := s.call(notificationEndpointDAO, request, &resp); err != nil {
		return "", err
	}
	return resp.ID, nil
}

// Update replaces the stored endpoint of entry.Name with entry.
//
// Two things to know. The dao takes the whole entry, not a patch: a field left
// empty is stored empty, so read the current endpoint with Get and modify that
// rather than building one from scratch. And update is an UPSERT -- a name the
// dao does not hold is created rather than refused, so a typo in entry.Name
// silently adds a second endpoint instead of editing the one meant.
func (s *NotificationService) Update(entry *model.EndpointConfig) error {
	if entry == nil {
		return fmt.Errorf("notification endpoint entry is required")
	}
	request := &GenericDAORequest[model.EndpointConfig]{
		Action: "update",
		Args: &GenericDAORequestArgs[model.EndpointConfig]{
			Id:    entry.Name,
			Entry: entry,
		},
	}
	return s.call(notificationEndpointDAO, request, nil)
}

// AddEmail stores an Email endpoint delivering to the given recipients. action
// names the fpl_action script that renders the mail, e.g. "Generic_Email_Action".
// See Add on the returned id.
func (s *NotificationService) AddEmail(name string, action string, to []string, cc []string) (id string, err error) {
	return s.Add(&model.EndpointConfig{
		Name:        name,
		Integration: "Email",
		Action:      action,
		Email: &model.EndpointEmailConfig{
			To: to,
			Cc: cc,
		},
	})
}

// AddSlack stores a Slack endpoint posting to the given channels.
// integrationName is the name of the configured Slack integration
// (platform_integration_dao) that holds the token to post with.
func (s *NotificationService) AddSlack(name string, action string, integrationName string, channels []string) (id string, err error) {
	slack := &model.EndpointSlackConfig{
		IntegrationName: integrationName,
		Channels:        channels,
	}
	if len(channels) == 1 {
		slack.Channel = channels[0]
	}
	return s.Add(&model.EndpointConfig{
		Name:        name,
		Integration: "Slack",
		Action:      action,
		Slack:       slack,
	})
}

// Get returns one notification endpoint by name.
//
// A name the dao does not hold is an ERROR, not an empty result: the call fails
// with "export not found: <name>". The nil return is only for a dao that answers
// with a null entry, which the server is not observed to do.
func (s *NotificationService) Get(name string) (*model.EndpointConfig, error) {
	request := &GenericDAORequest[model.EndpointConfig]{
		Action: "get",
		Args: &GenericDAORequestArgs[model.EndpointConfig]{
			Id: name,
		},
	}
	var resp GenericDaoEntryResponse[model.EndpointConfig]
	if err := s.call(notificationEndpointDAO, request, &resp); err != nil {
		return nil, err
	}
	return resp.Entry, nil
}

// Delete removes the endpoint of that name. Deleting a name the dao does not
// hold fails with "unknown export" rather than succeeding quietly.
func (s *NotificationService) Delete(name string) (err error) {
	request := &GenericDAORequest[model.EndpointConfig]{
		Action: "delete",
		Args: &GenericDAORequestArgs[model.EndpointConfig]{
			Id: name,
		},
	}
	return s.call(notificationEndpointDAO, request, nil)
}

func (s *NotificationService) List() (endpoints []*model.EndpointConfig, err error) {
	request := &GenericDAORequest[model.EndpointConfig]{
		Action: "list",
	}
	var resp GenericDaoListResponse[model.EndpointConfig]
	if err := s.call(notificationEndpointDAO, request, &resp); err != nil {
		return nil, err
	}
	return resp.Entries, nil
}

// ListActions returns the actions a notification endpoint can name: the
// fpl_action scripts whose ActionConfig.Target is "Platform Notification". If
// integration is non-empty, only the actions expecting that endpoint kind
// ("Email", "Slack") are returned.
func (s *NotificationService) ListActions(integration string) ([]*model.FPLScript, error) {
	actions, err := NewPlatformService(s.client).ListActions()
	if err != nil {
		return nil, err
	}
	filtered := make([]*model.FPLScript, 0, len(actions))
	for _, a := range actions {
		if a == nil || a.ActionConfig == nil {
			continue
		}
		if a.ActionConfig.Target != "Platform Notification" {
			continue
		}
		if integration != "" && a.ActionConfig.Integration != integration {
			continue
		}
		filtered = append(filtered, a)
	}
	return filtered, nil
}
