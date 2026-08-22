package api

import (
	"fmt"

	ingextAPI "github.com/SecurityDo/ingext_api/api"
	"github.com/SecurityDo/ingext_api/model"
)

func (c *Client) ResourceSearch(resourceType string, customer string) (resp *model.LakeSearchResponse, err error) {

	fmt.Printf("Searching for resource type '%s' and customer '%s'...\n", resourceType, customer)
	service := ingextAPI.NewResourceService(c.ingextClient)

	resp, err = service.Search(resourceType, customer)

	if err != nil {
		c.Logger.Error("resource search error", "error", err)
		return nil, fmt.Errorf("resource search error: %w", err)
	}
	return resp, nil
}

// DeleteResourceDump purges every resource dump of one customer via the
// ingext_resource_dump_delete API and returns the number of dump folders removed.
func (c *Client) DeleteResourceDump(customer string) (int, error) {
	service := ingextAPI.NewResourceService(c.ingextClient)

	deleted, err := service.DeleteResourceDump(customer)
	if err != nil {
		c.Logger.Error("resource dump delete error", "customer", customer, "error", err)
		return 0, fmt.Errorf("resource dump delete error: %w", err)
	}
	return deleted, nil
}
