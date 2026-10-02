package api

import (
	"fmt"

	ingextAPI "github.com/SecurityDo/ingext_api/api"
	"github.com/SecurityDo/ingext_api/model"
)

func (c *Client) FacetList() ([]*model.FacetConfigEntry, error) {
	service := ingextAPI.NewFacetService(c.ingextClient)

	entries, err := service.List()
	if err != nil {
		c.Logger.Error("list facets error", "error", err)
		return nil, fmt.Errorf("list facets error: %w", err)
	}
	return entries, nil
}

func (c *Client) FacetGet(name string) (*model.FacetConfigEntry, error) {
	service := ingextAPI.NewFacetService(c.ingextClient)

	entry, err := service.GetByName(name)
	if err != nil {
		c.Logger.Error("get facet error", "error", err)
		return nil, fmt.Errorf("get facet error: %w", err)
	}
	if entry == nil {
		return nil, fmt.Errorf("facet %q not found", name)
	}
	return entry, nil
}

func (c *Client) FacetApply(entry *model.FacetConfigEntry) (bool, error) {
	service := ingextAPI.NewFacetService(c.ingextClient)

	created, err := service.Apply(entry)
	if err != nil {
		c.Logger.Error("apply facet error", "error", err)
		return false, fmt.Errorf("apply facet error: %w", err)
	}
	return created, nil
}

func (c *Client) FacetDelete(name string) error {
	service := ingextAPI.NewFacetService(c.ingextClient)

	entry, err := service.GetByName(name)
	if err != nil {
		c.Logger.Error("delete facet error", "error", err)
		return fmt.Errorf("delete facet error: %w", err)
	}
	if entry == nil {
		return fmt.Errorf("facet %q not found", name)
	}
	if err := service.Delete(entry.Id); err != nil {
		c.Logger.Error("delete facet error", "error", err)
		return fmt.Errorf("delete facet error: %w", err)
	}
	return nil
}
