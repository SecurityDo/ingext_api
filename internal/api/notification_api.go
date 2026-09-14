package api

import (
	"fmt"

	ingextAPI "github.com/SecurityDo/ingext_api/api"
	"github.com/SecurityDo/ingext_api/model"
)

func (c *Client) NotificationList() ([]*model.EndpointConfig, error) {
	service := ingextAPI.NewNotificationService(c.ingextClient)

	endpoints, err := service.List()
	if err != nil {
		c.Logger.Error("list notification endpoints error", "error", err)
		return nil, fmt.Errorf("list notification endpoints error: %w", err)
	}
	return endpoints, nil
}

func (c *Client) NotificationDelete(name string) error {
	service := ingextAPI.NewNotificationService(c.ingextClient)

	if err := service.Delete(name); err != nil {
		c.Logger.Error("delete notification endpoint error", "error", err)
		return fmt.Errorf("delete notification endpoint error: %w", err)
	}
	return nil
}

func (c *Client) NotificationAddEmail(name string, action string, to []string, cc []string) (string, error) {
	service := ingextAPI.NewNotificationService(c.ingextClient)

	id, err := service.AddEmail(name, action, to, cc)
	if err != nil {
		c.Logger.Error("add email notification endpoint error", "error", err)
		return "", fmt.Errorf("add email notification endpoint error: %w", err)
	}
	return id, nil
}

func (c *Client) NotificationGet(name string) (*model.EndpointConfig, error) {
	service := ingextAPI.NewNotificationService(c.ingextClient)

	entry, err := service.Get(name)
	if err != nil {
		c.Logger.Error("get notification endpoint error", "error", err)
		return nil, fmt.Errorf("get notification endpoint error: %w", err)
	}
	// A missing name already came back as an error above; this covers a dao
	// that answers with a null entry instead.
	if entry == nil {
		return nil, fmt.Errorf("notification endpoint %q not found", name)
	}
	return entry, nil
}

func (c *Client) NotificationAddSlack(name string, action string, integrationName string, channels []string) (string, error) {
	service := ingextAPI.NewNotificationService(c.ingextClient)

	id, err := service.AddSlack(name, action, integrationName, channels)
	if err != nil {
		c.Logger.Error("add slack notification endpoint error", "error", err)
		return "", fmt.Errorf("add slack notification endpoint error: %w", err)
	}
	return id, nil
}

func (c *Client) NotificationUpdate(entry *model.EndpointConfig) error {
	service := ingextAPI.NewNotificationService(c.ingextClient)

	if err := service.Update(entry); err != nil {
		c.Logger.Error("update notification endpoint error", "error", err)
		return fmt.Errorf("update notification endpoint error: %w", err)
	}
	return nil
}

// NotificationListActions lists the actions a notification endpoint can name,
// optionally narrowed to one integration ("Email", "Slack").
func (c *Client) NotificationListActions(integration string) ([]*model.FPLScript, error) {
	service := ingextAPI.NewNotificationService(c.ingextClient)

	actions, err := service.ListActions(integration)
	if err != nil {
		c.Logger.Error("list notification actions error", "error", err)
		return nil, fmt.Errorf("list notification actions error: %w", err)
	}
	return actions, nil
}
