package api

import (
	"encoding/json"
	"fmt"

	"github.com/SecurityDo/ingext_api/client"
	"github.com/SecurityDo/ingext_api/model"
)

// PlatformService provides helpers for calling platform_* endpoints.
type GridService struct {
	client *client.IngextClient
}

// NewGridService constructs a GridService instance backed by the provided client.
func NewGridService(client *client.IngextClient) *GridService {
	return &GridService{client: client}
}

//func (s *GridService) call(function string, payload interface{}, out interface{}) error {
//	return ApiCall(s.client, function, payload, out)
//}

func (s *GridService) gridCall(function string, payload interface{}, out interface{}) error {
	return ApiCallWithPrefix(s.client, "api/grid", function, payload, out)
}

func (s *GridService) ListAccount() (resp *model.ListFluencyAccountsResponse, err error) {
	//var resp kqlModel.KQLSearchResponse
	if err := s.gridCall("get_grid_accounts", nil, &resp); err != nil {
		return nil, err
	}
	return resp, nil
}

// ProbeAccounts lists grid accounts without logging on failure. A 404 here means
// the site is not a grid manager, which callers use to detect proxy support.
func (s *GridService) ProbeAccounts() (resp *model.ListFluencyAccountsResponse, err error) {
	res, err := s.client.GenericCallQuiet("api/grid", "get_grid_accounts", nil)
	if err != nil {
		return nil, err
	}
	if res == nil {
		return nil, fmt.Errorf("empty response from get_grid_accounts")
	}
	if err := json.Unmarshal(res.GetBytes(), &resp); err != nil {
		return nil, fmt.Errorf("failed to parse get_grid_accounts response: %w", err)
	}
	return resp, nil
}

func (s *GridService) AddSaasAccount(req *model.GridAddSaasAccountRequest) (err error) {
	//var resp kqlModel.KQLSearchResponse
	if err := s.gridCall("add_saas_account", req, nil); err != nil {
		return err
	}
	return nil
}

func (s *GridService) DeleteSaasAccount(req *model.GridDeleteSaasAccountRequest) (err error) {
	//var resp kqlModel.KQLSearchResponse
	if err := s.gridCall("delete_saas_account", req, nil); err != nil {
		return err
	}
	return nil
}
