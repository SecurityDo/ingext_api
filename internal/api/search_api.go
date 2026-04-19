package api

import (
	"fmt"

	ingextAPI "github.com/SecurityDo/ingext_api/api"
	kqlModel "github.com/SecurityDo/ingext_api/kql/model"
)

func (c *Client) KQLSearch(kql string) (resp *kqlModel.KQLSearchResponse, err error) {

	service := ingextAPI.NewSearchService(c.ingextClient)

	resp, err = service.KQLSearch(kql)

	if err != nil {
		c.Logger.Error("kql search error", "error", err)
		return nil, fmt.Errorf("kql search error: %w", err)
	}
	return resp, nil
}

// KQLValidate parses a KQL query on the search service without executing it.
func (c *Client) KQLValidate(kql string) (*ingextAPI.KQLValidateResponse, error) {
	service := ingextAPI.NewSearchService(c.ingextClient)
	resp, err := service.KQLValidate(kql)
	if err != nil {
		c.Logger.Error("kql validate error", "error", err)
		return nil, fmt.Errorf("kql validate error: %w", err)
	}
	return resp, nil
}
