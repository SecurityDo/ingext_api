package api

import (
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
