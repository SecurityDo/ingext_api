package api

import (
	"fmt"

	ingextAPI "github.com/SecurityDo/ingext_api/api"
	model "github.com/SecurityDo/ingext_api/model"
)

func (c *Client) AddIntegration(entry *model.Integration) (id string, err error) {

	platformService := ingextAPI.NewPlatformService(c.ingextClient)

	id, err = platformService.AddIntegration(entry)

	if err != nil {
		c.Logger.Error("failed to add integration", "error", err)
		return "", fmt.Errorf("failed to add integration: %s", err.Error())
	}
	return id, nil
}

func (c *Client) DeleteIntegration(id string) (err error) {

	platformService := ingextAPI.NewPlatformService(c.ingextClient)

	err = platformService.DeleteIntegration(id)

	if err != nil {
		c.Logger.Error("failed to delete integration", "error", err)
		return fmt.Errorf("failed to delete integration: %s", err.Error())
	}
	return nil
}

func (c *Client) ListIntegration() (entries []*model.Integration, err error) {

	platformService := ingextAPI.NewPlatformService(c.ingextClient)

	entries, err = platformService.ListIntegrations()

	if err != nil {
		c.Logger.Error("failed to list integration", "error", err)
		return nil, fmt.Errorf("failed to list integration: %s", err.Error())
	}
	return entries, nil
}
