package api

import (
	"fmt"

	ingextAPI "github.com/SecurityDo/ingext_api/api"
	"github.com/SecurityDo/ingext_api/model"
)

func (c *Client) ListAccount() (resp *model.ListFluencyAccountsResponse, err error) {
	service := ingextAPI.NewGridService(c.ingextClient)

	resp, err = service.ListAccount()
	if err != nil {
		c.Logger.Error("list account error", "error", err)
		return nil, fmt.Errorf("list account error: %w", err)
	}
	return resp, nil
}

func (c *Client) AddSaasAccount(req *model.GridAddSaasAccountRequest) error {
	service := ingextAPI.NewGridService(c.ingextClient)

	if err := service.AddSaasAccount(req); err != nil {
		c.Logger.Error("add saas account error", "error", err)
		return fmt.Errorf("add saas account error: %w", err)
	}
	return nil
}
