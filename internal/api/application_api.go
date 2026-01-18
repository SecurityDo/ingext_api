package api

import (
	"fmt"

	"github.com/SecurityDo/ingext_api/api"
	ingextAPI "github.com/SecurityDo/ingext_api/api"
	"github.com/SecurityDo/ingext_api/model"
)

func (c *Client) ListAppTemplates() (templates []*model.ApplicationTemplateConfig, err error) {

	applicationService := ingextAPI.NewApplicationService(c.ingextClient)

	resp, err := applicationService.ListAppTemplates()

	if err != nil {
		c.Logger.Error("failed to list application templates", "error", err)
		return nil, fmt.Errorf("failed to list application templates: %w", err)
	}
	return resp.Entries, nil
}

func (c *Client) InstallAppInstance(application, instance string, displayName string, parameters []*model.InputParameter) (err error) {

	applicationService := ingextAPI.NewApplicationService(c.ingextClient)

	req := &api.InstallAppInstanceRequest{
		Config: &model.InstanceConfig{
			Application:     application,
			Instance:        instance,
			DisplayName:     displayName,
			InputParameters: parameters,
		},
	}

	err = applicationService.InstallAppTemplate(req)

	if err != nil {
		c.Logger.Error("failed to install application instance", "error", err)
		return fmt.Errorf("failed to install application instance: %w", err)
	}
	return nil
}

func (c *Client) UnInstallAppInstance(application, instance string) (err error) {

	applicationService := ingextAPI.NewApplicationService(c.ingextClient)

	req := &api.UnInstallAppInstanceRequest{
		Application: application,
		Instance:    instance,
	}

	err = applicationService.UnInstallAppTemplate(req)

	if err != nil {
		c.Logger.Error("failed to uninstall application instance", "error", err)
		return fmt.Errorf("failed to uninstall application instance: %w", err)
	}
	return nil
}
