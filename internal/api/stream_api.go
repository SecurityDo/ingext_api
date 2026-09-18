package api

import (
	"fmt"

	"github.com/SecurityDo/ingext_api/api"
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

func (c *Client) AddSimpleRouter(processorName string, routerName string) (id string, err error) {

	platformService := ingextAPI.NewPlatformService(c.ingextClient)

	id, err = platformService.AddSimpleRouter(processorName, routerName)

	if err != nil {
		c.Logger.Error("failed to add router", "error", err)
		return "", fmt.Errorf("failed to add router: %s", err.Error())
	}
	return id, nil
}

func (c *Client) SetRouterSink(routerID, sinkID string) (err error) {

	platformService := ingextAPI.NewPlatformService(c.ingextClient)

	result, err := platformService.GetRouter(routerID)

	if err != nil {
		c.Logger.Error("failed to get router by ID", "error", err)
		return fmt.Errorf("failed get router by ID: %s", err.Error())
	}
	if len(result.Pipes) == 0 {
		c.Logger.Error("router has no pipes", "routerID", routerID)
		return fmt.Errorf("router has no pipes: %s", routerID)
	}
	pipeConfig := result.Pipes[len(result.Pipes)-1] // Get the last pipe

	for _, sid := range pipeConfig.SinkIDs {
		if sid == sinkID {
			c.Logger.Error("sink already exists in router", "routerID", routerID, "sinkID", sinkID)
			return fmt.Errorf("sink already exists in router: %s, sink: %s", routerID, sinkID)
		}
	}
	pipeConfig.SinkIDs = append(pipeConfig.SinkIDs, sinkID)

	req := &api.PipeUpdateReq{
		RouterID:   pipeConfig.RouterID,
		PipeConfig: pipeConfig,
	}

	err = platformService.UpdatePipe(req)
	if err != nil {
		c.Logger.Error("failed to update pipe with new sink", "error", err)
		return fmt.Errorf("failed to update pipe with new sink: %s", err.Error())
	}

	return nil
}

func (c *Client) UpdatePipeProcessor(routerName, pipeName, processorName string) (err error) {

	platformService := ingextAPI.NewPlatformService(c.ingextClient)

	err = platformService.UpdatePipeProcessor(&api.PipeProcessorUpdateReq{
		RouterName:    routerName,
		PipeName:      pipeName,
		ProcessorName: processorName,
	})

	if err != nil {
		c.Logger.Error("failed to update pipe processor", "error", err)
		return fmt.Errorf("failed to update pipe processor: %s", err.Error())
	}
	return nil
}

func (c *Client) SetSourceRouter(sourceID, routerID string) (err error) {

	platformService := ingextAPI.NewPlatformService(c.ingextClient)

	err = platformService.SetDataSourceRouter(&api.SourceSetRouterReq{
		DataSourceID: sourceID,
		RouterID:     routerID,
	})

	if err != nil {
		c.Logger.Error("failed to connect source to router", "error", err)
		return fmt.Errorf("failed to connect source to router: %s", err.Error())
	}
	return nil
}

/*
func (c *Client) AddPipe(routerConfig *model.StreamPipeConfig) (id string, err error) {

	platformService := ingextAPI.NewPlatformService(c.ingextClient)

	resp, err := platformService.AddRouter(routerConfig)

	if err != nil {
		c.Logger.Error("failed to add router", "error", err)
		return "", fmt.Errorf("failed to add router: %s", err.Error())
	}
	return resp.ID, nil
}*/

// ListConfigs returns the whole platform topology of the target account:
// sources, sinks, routers, pipes, connections and error states. platform_router_dao
// has no "list" action, so this is the only inventory call.
func (c *Client) ListConfigs() (*ingextAPI.ListConfigsResponse, error) {

	platformService := ingextAPI.NewPlatformService(c.ingextClient)

	resp, err := platformService.ListConfigs()

	if err != nil {
		c.Logger.Error("failed to list platform configs", "error", err)
		return nil, fmt.Errorf("failed to list platform configs: %s", err.Error())
	}
	return resp, nil
}

// AddRouterPipe creates a new pipe under an existing router, leaving the pipes
// already on it untouched.
func (c *Client) AddRouterPipe(routerID string, pipe *model.StreamPipeConfig) (id string, err error) {

	platformService := ingextAPI.NewPlatformService(c.ingextClient)

	resp, err := platformService.AddRouterPipe(&ingextAPI.RouterAddPipeReq{
		RouterID:   routerID,
		PipeConfig: pipe,
	})

	if err != nil {
		c.Logger.Error("failed to add router pipe", "error", err)
		return "", fmt.Errorf("failed to add router pipe: %s", err.Error())
	}
	return resp.ID, nil
}

// DeleteRouterPipe removes one pipe from a router.
func (c *Client) DeleteRouterPipe(routerID, pipeID string) (err error) {

	platformService := ingextAPI.NewPlatformService(c.ingextClient)

	err = platformService.DeleteRouterPipe(&ingextAPI.RouterDeletePipeReq{
		RouterID: routerID,
		PipeID:   pipeID,
	})

	if err != nil {
		c.Logger.Error("failed to delete router pipe", "error", err)
		return fmt.Errorf("failed to delete router pipe: %s", err.Error())
	}
	return nil
}

// DeleteRouter removes a router by id.
func (c *Client) DeleteRouter(routerID string) (err error) {

	platformService := ingextAPI.NewPlatformService(c.ingextClient)

	err = platformService.DeleteRouter(routerID)

	if err != nil {
		c.Logger.Error("failed to delete router", "error", err)
		return fmt.Errorf("failed to delete router: %s", err.Error())
	}
	return nil
}

// ReloadDataSource restarts one data source. For a plugin source this re-forks
// the plugin, which is how a newly published binary is picked up.
func (c *Client) ReloadDataSource(dataSourceID string) (err error) {

	platformService := ingextAPI.NewPlatformService(c.ingextClient)

	err = platformService.SourceReload(dataSourceID)

	if err != nil {
		c.Logger.Error("failed to reload data source", "error", err)
		return fmt.Errorf("failed to reload data source: %s", err.Error())
	}
	return nil
}
