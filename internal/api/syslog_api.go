package api

import (
	"fmt"

	ingextAPI "github.com/SecurityDo/ingext_api/api"
	"github.com/SecurityDo/ingext_api/model"
)

func (c *Client) RegisterSyslogConfig(ports []string) (*model.GetSyslogConfigResponse, error) {
	service := ingextAPI.NewSyslogService(c.ingextClient)

	resp, err := service.Register(ports)
	if err != nil {
		c.Logger.Error("register syslog config error", "error", err)
		return nil, fmt.Errorf("register syslog config error: %w", err)
	}
	return resp, nil
}

func (c *Client) UpdateSyslogConfig(ports []string) (*model.GetSyslogConfigResponse, error) {
	service := ingextAPI.NewSyslogService(c.ingextClient)

	resp, err := service.Update(ports)
	if err != nil {
		c.Logger.Error("update syslog config error", "error", err)
		return nil, fmt.Errorf("update syslog config error: %w", err)
	}
	return resp, nil
}

func (c *Client) GetSyslogConfig() (*model.GetSyslogConfigResponse, error) {
	service := ingextAPI.NewSyslogService(c.ingextClient)

	resp, err := service.Get()
	if err != nil {
		c.Logger.Error("get syslog config error", "error", err)
		return nil, fmt.Errorf("get syslog config error: %w", err)
	}
	return resp, nil
}

func (c *Client) DeleteSyslogConfig() error {
	service := ingextAPI.NewSyslogService(c.ingextClient)

	if err := service.Delete(); err != nil {
		c.Logger.Error("delete syslog config error", "error", err)
		return fmt.Errorf("delete syslog config error: %w", err)
	}
	return nil
}
