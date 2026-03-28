package api

import (
	"github.com/SecurityDo/ingext_api/client"
	"github.com/SecurityDo/ingext_api/model"
)

// PlatformService provides helpers for calling platform_* endpoints.
type NotificationService struct {
	client *client.IngextClient
}

func NewNotificationService(client *client.IngextClient) *NotificationService {
	return &NotificationService{client: client}
}

func (s *NotificationService) call(function string, payload interface{}, out interface{}) error {
	return ApiCall(s.client, function, payload, out)
}

func (s *NotificationService) AddEmail(name string, action string, to []string, cc []string) (id string, err error) {
	request := &GenericDAORequest[model.EndpointConfig]{
		Action: "add",
		Args: &GenericDAORequestArgs[model.EndpointConfig]{
			//Id:    args.Id,
			Entry: &model.EndpointConfig{
				Name:        name,
				Integration: "Email",
				Action:      action,
				Email: &model.EndpointEmailConfig{
					To: to,
					Cc: cc,
				},
			},
		},
	}
	var resp GenericDaoAddResponse
	if err := s.call("platform_notification_endpoint_dao", request, &resp); err != nil {
		return "", err
	}
	return resp.ID, nil
}

func (s *NotificationService) Delete(name string) (err error) {
	request := &GenericDAORequest[model.EndpointConfig]{
		Action: "delete",
		Args: &GenericDAORequestArgs[model.EndpointConfig]{
			Id: name,
		},
	}
	//var resp GenericDaoAddResponse
	if err := s.call("platform_notification_endpoint_dao", request, nil); err != nil {
		return err
	}
	return nil

}

func (s *NotificationService) List() (endpoints []*model.EndpointConfig, err error) {
	request := &GenericDAORequest[model.EndpointConfig]{
		Action: "list",
	}
	var resp GenericDaoListResponse[model.EndpointConfig]
	if err := s.call("platform_notification_endpoint_dao", request, &resp); err != nil {
		return nil, err
	}
	return resp.Entries, nil

}
