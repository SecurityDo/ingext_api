package api

import (
	"fmt"

	ingextAPI "github.com/SecurityDo/ingext_api/api"
	model "github.com/SecurityDo/ingext_api/model"
)

// AddProcessor deploys a new processor. action is the ActionConfig an
// "fpl_action" needs and must be nil for the other types.
func (c *Client) AddProcessor(name, content, processorType, description, channel string, action *model.FplActionConfig) (err error) {

	platformService := ingextAPI.NewPlatformService(c.ingextClient)

	if processorType == "" {
		processorType = "fpl_processor" // Default to JavaScript if not specified
	}

	entry := &model.FPLScript{
		Name:         name,
		ScriptText:   content,
		Type:         processorType,
		Description:  description,
		Channel:      channel,
		ActionConfig: action,
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
// An empty processorType, description or channel keeps the stored one, and so
// does a nil action or an empty field of it. Changing the type away from
// "fpl_action" drops the stored ActionConfig.
func (c *Client) UpdateProcessor(name, content, processorType, description, channel string, action *model.FplActionConfig) (err error) {

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
	if channel != "" {
		entry.Channel = channel
	}
	if action != nil {
		merged := model.FplActionConfig{}
		if entry.ActionConfig != nil {
			merged = *entry.ActionConfig
		}
		if action.Target != "" {
			merged.Target = action.Target
		}
		if action.Integration != "" {
			merged.Integration = action.Integration
		}
		entry.ActionConfig = &merged
	}
	if entry.Type != "fpl_action" && action == nil {
		entry.ActionConfig = nil
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

// GetProcessor returns one processor's stored entry, script included.
func (c *Client) GetProcessor(name string) (*model.FPLScript, error) {
	platformService := ingextAPI.NewPlatformService(c.ingextClient)
	entry, err := platformService.GetProcessor(name)
	if err != nil {
		c.Logger.Error("failed to get processor", "name", name, "error", err)
		return nil, fmt.Errorf("failed to get processor %s: %s", name, err.Error())
	}
	if entry == nil {
		return nil, fmt.Errorf("processor %s not found", name)
	}
	return entry, nil
}

// ListActions returns the tenant's fpl_action scripts with their ActionConfig.
// The processor DAO's list returns parsers only, so actions come from
// platform_list_actions.
func (c *Client) ListActions() (entries []*model.FPLScript, err error) {

	platformService := ingextAPI.NewPlatformService(c.ingextClient)

	entries, err = platformService.ListActions()

	if err != nil {
		c.Logger.Error("failed to list actions", "error", err)
		return nil, fmt.Errorf("failed to list actions: %s", err.Error())
	}
	return entries, nil
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

// ValidateProcessorScript compiles a script that has not been deployed, without
// running it: a script that would fail on every event still validates, and so
// does one with no main(). A parse error is reported in the result rather than
// as an error.
func (c *Client) ValidateProcessorScript(script string) (*ingextAPI.FPLProcessorValidateResult, error) {

	platformService := ingextAPI.NewPlatformService(c.ingextClient)

	resp, err := platformService.ValidateProcessorScript(script)

	if err != nil {
		c.Logger.Error("failed to validate processor script", "error", err)
		return nil, fmt.Errorf("failed to validate processor script: %s", err.Error())
	}
	return resp, nil
}

// ValidateDeployedProcessor compiles the stored processor of that name. A name
// that is not deployed comes back as a failed validation, not an error.
func (c *Client) ValidateDeployedProcessor(name string) (*ingextAPI.FPLProcessorValidateResult, error) {

	platformService := ingextAPI.NewPlatformService(c.ingextClient)

	resp, err := platformService.ValidateDeployedProcessor(name)

	if err != nil {
		c.Logger.Error("failed to validate processor", "name", name, "error", err)
		return nil, fmt.Errorf("failed to validate processor %s: %s", name, err.Error())
	}
	return resp, nil
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
