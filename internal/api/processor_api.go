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

// UpdateProcessor replaces the script of an existing processor, leaving
// everything the caller is not changing alone. The DAO update takes a whole
// entry, and a processor carries more than its script -- the id it is keyed by,
// the repository and gitPath it was imported from, its group and tags -- so the
// current entry is read back and patched rather than rebuilt from the flags.
//
// An empty processorType or description keeps the stored one.
func (c *Client) UpdateProcessor(name, content, processorType, description string) (err error) {

	platformService := ingextAPI.NewPlatformService(c.ingextClient)

	entry, err := platformService.GetProcessor(name)
	if err != nil {
		c.Logger.Error("failed to get processor", "name", name, "error", err)
		return fmt.Errorf("failed to get processor %s: %s", name, err.Error())
	}
	// The DAO answers a missing name with an ERROR verdict rather than a null
	// entry, so this is the belt to that suspenders -- either way update never
	// falls through to writing a fresh entry.
	if entry == nil {
		return fmt.Errorf("processor %s not found: add it first", name)
	}

	entry.ScriptText = content
	if processorType != "" {
		entry.Type = processorType
	}
	if description != "" {
		entry.Description = description
	}

	if err = platformService.UpdateProcessor(entry); err != nil {
		c.Logger.Error("failed to update processor", "name", name, "error", err)
		return fmt.Errorf("failed to update processor %s: %s", name, err.Error())
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
