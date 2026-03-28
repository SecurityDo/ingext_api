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
