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
