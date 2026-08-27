package api

import (
	"fmt"

	ingextAPI "github.com/SecurityDo/ingext_api/api"
	model "github.com/SecurityDo/ingext_api/model"
)

func (c *Client) AddProcessor(name, content, processorType, description string) (err error) {

	platformService := ingextAPI.NewPlatformService(c.ingextClient)

	if processorType == "" {
		processorType = "fpl_processor" // Default to JavaScript if not specified
	}

	entry := &model.FPLScript{
		Name:        name,
		ScriptText:  content,
		Type:        processorType,
		Description: description,
	}

	err = platformService.AddProcessor(entry)

	if err != nil {
		c.Logger.Error("failed to add processor", "error", err)
		return fmt.Errorf("failed to add processor: %s", err.Error())
	}
	return nil
}

func (c *Client) DeleteProcessor(name string) (err error) {

	platformService := ingextAPI.NewPlatformService(c.ingextClient)

	err = platformService.DeleteProcessor(name)

	if err != nil {
		c.Logger.Error("failed to delete processor", "name", name, "error", err)
		return fmt.Errorf("failed to delete processor %s: %s", name, err.Error())
	}
	return nil
}

func (c *Client) ListProcessor() (entries []*model.FPLScript, err error) {

	platformService := ingextAPI.NewPlatformService(c.ingextClient)

	entries, err = platformService.ListProcessors()

	if err != nil {
		c.Logger.Error("failed to list processor", "error", err)
		return nil, fmt.Errorf("failed to list processor: %s", err.Error())
	}
	return entries, nil
}

// TestProcessor runs an FPL script against one sample document without
// deploying it. source is what the runtime named by processorType reads: an
// fpl_processor takes the JSON document envelope built by
// ingextAPI.NewFPLTestDocument, other types take their raw payload.
//
// A script that errors on the document is reported in the result, not as an
// error: only a transport or endpoint failure returns one.
func (c *Client) TestProcessor(script, source, processorType string) (*ingextAPI.FPLProcessorTestResult, error) {

	platformService := ingextAPI.NewPlatformService(c.ingextClient)

	resp, err := platformService.TestProcessor(&ingextAPI.FPLProcessorTestRequest{
		Script: script,
		Source: source,
		Type:   processorType,
	})

	if err != nil {
		c.Logger.Error("failed to test processor", "type", processorType, "error", err)
		return nil, fmt.Errorf("failed to test processor: %s", err.Error())
	}
	return resp, nil
}
