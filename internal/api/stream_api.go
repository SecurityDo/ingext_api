package api

import (
	"fmt"

	ingextAPI "github.com/SecurityDo/ingext_api/api"
	model "github.com/SecurityDo/ingext_api/model"
)

func (c *Client) AddDataSource(source *model.DataSourceConfig) (resp *ingextAPI.AddDataSourceResponse, err error) {

	platformService := ingextAPI.NewPlatformService(c.ingextClient)

	resp, err = platformService.AddDataSource(source)

	if err != nil {
		c.Logger.Error("failed to add data source", "error", err)
		return nil, fmt.Errorf("failed to add data source: %s", err.Error())
	}
	return resp, nil
}

func (c *Client) DeleteDataSource(id string) (err error) {

	platformService := ingextAPI.NewPlatformService(c.ingextClient)

	err = platformService.DeleteDataSource(id)

	if err != nil {
		c.Logger.Error("failed to delete data source", "error", err)
		return fmt.Errorf("failed to delete data source: %s", err.Error())
	}
	return nil
}

func (c *Client) ListDataSource() (entries []*model.DataSourceConfig, err error) {

	platformService := ingextAPI.NewPlatformService(c.ingextClient)

	entries, err = platformService.ListDataSource()

	if err != nil {
		c.Logger.Error("failed to list data source", "error", err)
		return nil, fmt.Errorf("failed to list data source: %s", err.Error())
	}
	return entries, nil
}

func (c *Client) AddDataSink(sink *model.DataSinkConfig) (resp *ingextAPI.AddDataSinkResponse, err error) {

	platformService := ingextAPI.NewPlatformService(c.ingextClient)

	resp, err = platformService.AddDataSink(sink)

	if err != nil {
		c.Logger.Error("failed to add data sink", "error", err)
		return nil, fmt.Errorf("failed to add data sink: %s", err.Error())
	}
	return resp, nil
}

func (c *Client) DeleteDataSink(id string) (err error) {

	platformService := ingextAPI.NewPlatformService(c.ingextClient)

	err = platformService.DeleteDataSink(id)

	if err != nil {
		c.Logger.Error("failed to delete data sink", "error", err)
		return fmt.Errorf("failed to delete data sink: %s", err.Error())
	}
	return nil
}

func (c *Client) ListDataSink() (entries []*model.DataSinkConfig, err error) {

	platformService := ingextAPI.NewPlatformService(c.ingextClient)

	entries, err = platformService.ListDataSink()

	if err != nil {
		c.Logger.Error("failed to list data sink", "error", err)
		return nil, fmt.Errorf("failed to list data sink: %s", err.Error())
	}
	return entries, nil
}

func (c *Client) AddRouter(routerConfig *model.RouterConfig) (id string, err error) {

	platformService := ingextAPI.NewPlatformService(c.ingextClient)

	resp, err := platformService.AddRouter(routerConfig)

	if err != nil {
		c.Logger.Error("failed to add router", "error", err)
		return "", fmt.Errorf("failed to add router: %s", err.Error())
	}
	return resp.ID, nil
}
