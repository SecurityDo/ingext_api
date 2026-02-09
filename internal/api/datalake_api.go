package api

import (
	"fmt"

	ingextAPI "github.com/SecurityDo/ingext_api/api"
	"github.com/SecurityDo/ingext_api/model"
)

func (c *Client) ListDatalakes() (entries []*model.Datalake, err error) {

	datalakeService := ingextAPI.NewDatalakeService(c.ingextClient)

	resp, err := datalakeService.ListDatalake()

	if err != nil {
		c.Logger.Error("failed to list datalakes", "error", err)
		return nil, fmt.Errorf("failed to list datalakes: %w", err)
	}
	return resp, nil
}

func (c *Client) AddDatalake(name string, managed bool, integrationID string) (err error) {

	datalakeService := ingextAPI.NewDatalakeService(c.ingextClient)

	err = datalakeService.AddDatalake(name, managed, integrationID)

	if err != nil {
		c.Logger.Error("failed to add datalakes", "error", err)
		return fmt.Errorf("failed to add datalakes: %w", err)
	}
	return nil
}

func (c *Client) AddDatalakeIndex(lake, index string, schema string) (err error) {

	datalakeService := ingextAPI.NewDatalakeService(c.ingextClient)

	err = datalakeService.AddDatalakeIndex(lake, index, schema)

	if err != nil {
		c.Logger.Error("failed to add datalake index", "error", err)
		return fmt.Errorf("failed to add datalake index: %w", err)
	}
	return nil
}

func (c *Client) ListDatalakeIndex(lake string) (entries []*model.DatalakeIndex, err error) {

	datalakeService := ingextAPI.NewDatalakeService(c.ingextClient)

	entries, err = datalakeService.ListDatalakeIndex(lake)

	if err != nil {
		c.Logger.Error("failed to list datalake index", "error", err)
		return nil, fmt.Errorf("failed to list datalake index: %w", err)
	}
	return entries, nil
}

func (c *Client) ListSchemas() (entries []*model.SchemaEntry, err error) {
	datalakeService := ingextAPI.NewDatalakeService(c.ingextClient)

	entries, err = datalakeService.ListSchema()
	if err != nil {
		c.Logger.Error("failed to list schemas", "error", err)
		return nil, fmt.Errorf("failed to list schemas: %w", err)
	}
	return entries, nil
}

func (c *Client) UpdateSchema(name, description, content string) error {
	datalakeService := ingextAPI.NewDatalakeService(c.ingextClient)

	err := datalakeService.UpdateSchema(name, description, content)
	if err != nil {
		c.Logger.Error("failed to update schema", "error", err)
		return fmt.Errorf("failed to update schema: %w", err)
	}
	return nil
}

func (c *Client) DeleteSchema(name string) error {
	datalakeService := ingextAPI.NewDatalakeService(c.ingextClient)

	err := datalakeService.DeleteSchema(name)
	if err != nil {
		c.Logger.Error("failed to delete schema", "error", err)
		return fmt.Errorf("failed to delete schema: %w", err)
	}
	return nil
}

func (c *Client) AddSchema(name, description, content string) error {
	datalakeService := ingextAPI.NewDatalakeService(c.ingextClient)

	err := datalakeService.AddSchema(name, description, content)
	if err != nil {
		c.Logger.Error("failed to add schema", "error", err)
		return fmt.Errorf("failed to add schema: %w", err)
	}
	return nil
}

func (c *Client) DeleteDatalakeIndex(lake, index string) (err error) {

	datalakeService := ingextAPI.NewDatalakeService(c.ingextClient)

	err = datalakeService.DeleteDatalakeIndex(lake, index)

	if err != nil {
		c.Logger.Error("failed to delete datalake index", "error", err)
		return fmt.Errorf("failed to delete datalake index: %w", err)
	}
	return nil
}
