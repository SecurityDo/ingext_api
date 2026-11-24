package api

import (
	"encoding/json"
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

func (c *Client) GetIntegration(id string) (*model.Integration, error) {
	entry, err := ingextAPI.NewPlatformService(c.ingextClient).GetIntegration(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get integration %s: %s", id, err.Error())
	}
	return entry, nil
}

// UpdateIntegration saves the entry (name, description, config, tags). The
// platform keeps the stored secret whatever the entry carries.
func (c *Client) UpdateIntegration(entry *model.Integration) error {
	if err := ingextAPI.NewPlatformService(c.ingextClient).UpdateIntegration(entry); err != nil {
		c.Logger.Error("failed to update integration", "error", err)
		return fmt.Errorf("failed to update integration: %s", err.Error())
	}
	return nil
}

// SetIntegrationSecret replaces the integration's whole secret.
func (c *Client) SetIntegrationSecret(id string, secret json.RawMessage) error {
	if err := ingextAPI.NewPlatformService(c.ingextClient).SetIntegrationSecret(id, secret); err != nil {
		return fmt.Errorf("failed to set integration secret: %s", err.Error())
	}
	return nil
}
