package api

import (
	"github.com/SecurityDo/ingext_api/client"
	"github.com/SecurityDo/ingext_api/model"
)

// PlatformService provides helpers for calling platform_* endpoints.
type ApplicationService struct {
	client *client.IngextClient
}

// NewPlatformService constructs a PlatformService instance backed by the provided client.
func NewApplicationService(client *client.IngextClient) *ApplicationService {
	return &ApplicationService{client: client}
}

func (s *ApplicationService) call(function string, payload interface{}, out interface{}) error {
	return ApiCall(s.client, function, payload, out)
}

func (s *ApplicationService) ListAppTemplates() (*ListAppTemplateResponse, error) {
	var resp ListAppTemplateResponse
	if err := s.call("platform_list_application_template", nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

type ListAppTemplateResponse struct {
	Entries []*model.ApplicationTemplateConfig `json:"entries"`
}

type ListAppInstanceResponse struct {
	Entries []*model.InstanceState `json:"entries"`
}

func (s *ApplicationService) ListAppInstances() (*ListAppInstanceResponse, error) {
	var resp ListAppInstanceResponse
	if err := s.call("platform_list_application_template", nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

type InstallAppInstanceRequest struct {
	Config *model.InstanceConfig `json:"config"`
}

func (s *ApplicationService) InstallAppTemplate(req *InstallAppInstanceRequest) error {
	//var resp ListAppTemplateResponse
	if err := s.call("platform_install_application_instance", req, nil); err != nil {
		return err
	}
	return nil
}

type UnInstallAppInstanceRequest struct {
	Application string `json:"application,omitempty" yaml:"application,omitempty"`
	Instance    string `json:"instance,omitempty" yaml:"instance,omitempty"`
}

func (s *ApplicationService) UnInstallAppTemplate(req *UnInstallAppInstanceRequest) error {
	//var resp ListAppTemplateResponse
	if err := s.call("platform_uninstall_application_instance", req, nil); err != nil {
		return err
	}
	return nil
}
