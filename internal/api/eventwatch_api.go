package api

import (
	"encoding/json"
	"fmt"

	fluencyAPI "github.com/SecurityDo/ingext_api/api"
	"github.com/SecurityDo/ingext_api/model"
)

// SummarySearch calls the overview summary search API with the given search string and time range.
func (c *Client) SummarySearch(searchString string, rangeFrom, rangeTo int64) (*model.ElasticSearchResult, error) {
	svc := fluencyAPI.NewEventWatchService(c.ingextClient)
	resp, err := svc.SummarySearch(searchString, rangeFrom, rangeTo)
	if err != nil {
		c.Logger.Error("failed to run summary search", "error", err)
		return nil, fmt.Errorf("failed to run summary search: %w", err)
	}
	return resp, nil
}

// TimelineSearch calls the fsm_behavior_search API with the given search string and time range.
func (c *Client) TimelineSearch(searchString string, rangeFrom, rangeTo int64) (*model.ElasticSearchResult, error) {
	svc := fluencyAPI.NewEventWatchService(c.ingextClient)
	resp, err := svc.TimelineSearch(searchString, rangeFrom, rangeTo)
	if err != nil {
		c.Logger.Error("failed to run timeline search", "error", err)
		return nil, fmt.Errorf("failed to run timeline search: %w", err)
	}
	return resp, nil
}

// RuleSearch calls the eventwatch_bucket_search API with the given search string (no time range).
func (c *Client) RuleSearch(searchString string) (*model.ElasticSearchResult, error) {
	svc := fluencyAPI.NewEventWatchService(c.ingextClient)
	resp, err := svc.RuleSearch(searchString)
	if err != nil {
		c.Logger.Error("failed to run rule search", "error", err)
		return nil, fmt.Errorf("failed to run rule search: %w", err)
	}
	return resp, nil
}

// ListRule lists all eventwatch rules via the eventwatch_bucket_dao API.
func (c *Client) ListRule() (*fluencyAPI.EventWatchRuleListResponse, error) {
	svc := fluencyAPI.NewEventWatchService(c.ingextClient)
	resp, err := svc.ListRule()
	if err != nil {
		c.Logger.Error("failed to list eventwatch rules", "error", err)
		return nil, fmt.Errorf("failed to list eventwatch rules: %w", err)
	}
	return resp, nil
}

// GetRule fetches a single eventwatch rule by name via the eventwatch_bucket_dao API.
func (c *Client) GetRule(name string) (*model.EventWatchBucket, error) {
	svc := fluencyAPI.NewEventWatchService(c.ingextClient)
	resp, err := svc.GetRule(name)
	if err != nil {
		c.Logger.Error("failed to get eventwatch rule", "name", name, "error", err)
		return nil, fmt.Errorf("failed to get eventwatch rule %q: %w", name, err)
	}
	return resp, nil
}

// TestRule runs a rule definition against one event via the eventwatch_rule_test
// API, without deploying the rule or storing what it produces.
func (c *Client) TestRule(rule *model.EventWatchBucket, event json.RawMessage) (*fluencyAPI.EventWatchRuleTestResult, error) {
	svc := fluencyAPI.NewEventWatchService(c.ingextClient)
	resp, err := svc.TestRuleEvent(rule, event)
	if err != nil {
		c.Logger.Error("failed to test eventwatch rule", "error", err)
		return nil, fmt.Errorf("failed to test eventwatch rule: %w", err)
	}
	return resp, nil
}

// TestDeployedRule reads the stored rule of that name and runs it against one
// event via the eventwatch_rule_test API.
func (c *Client) TestDeployedRule(name string, event json.RawMessage) (*fluencyAPI.EventWatchRuleTestResult, error) {
	svc := fluencyAPI.NewEventWatchService(c.ingextClient)
	resp, err := svc.TestDeployedRule(name, event)
	if err != nil {
		c.Logger.Error("failed to test eventwatch rule", "name", name, "error", err)
		return nil, fmt.Errorf("failed to test eventwatch rule %q: %w", name, err)
	}
	return resp, nil
}

// AddRule creates a new eventwatch rule via the eventwatch_bucket_dao API.
func (c *Client) AddRule(entry *model.EventWatchBucket) error {
	svc := fluencyAPI.NewEventWatchService(c.ingextClient)
	if err := svc.AddRule(entry); err != nil {
		c.Logger.Error("failed to add eventwatch rule", "error", err)
		return fmt.Errorf("failed to add eventwatch rule: %w", err)
	}
	return nil
}

// UpdateRule updates an existing eventwatch rule via the eventwatch_bucket_dao API.
func (c *Client) UpdateRule(entry *model.EventWatchBucket) error {
	svc := fluencyAPI.NewEventWatchService(c.ingextClient)
	if err := svc.UpdateRule(entry); err != nil {
		c.Logger.Error("failed to update eventwatch rule", "error", err)
		return fmt.Errorf("failed to update eventwatch rule: %w", err)
	}
	return nil
}

// ToggleRule flips the disabled state of an eventwatch rule via the eventwatch_bucket_dao API.
func (c *Client) ToggleRule(name string) error {
	svc := fluencyAPI.NewEventWatchService(c.ingextClient)
	if err := svc.ToggleRule(name); err != nil {
		c.Logger.Error("failed to toggle eventwatch rule", "name", name, "error", err)
		return fmt.Errorf("failed to toggle eventwatch rule %q: %w", name, err)
	}
	return nil
}

// DeleteRule removes an eventwatch rule via the eventwatch_bucket_dao API.
func (c *Client) DeleteRule(name string) error {
	svc := fluencyAPI.NewEventWatchService(c.ingextClient)
	if err := svc.DeleteRule(name); err != nil {
		c.Logger.Error("failed to delete eventwatch rule", "name", name, "error", err)
		return fmt.Errorf("failed to delete eventwatch rule %q: %w", name, err)
	}
	return nil
}

// DeleteRuleGroup removes an eventwatch rule group via the eventwatch_bucket_delete_group
// API and returns the number of rules deleted.
func (c *Client) DeleteRuleGroup(group string) (int, error) {
	svc := fluencyAPI.NewEventWatchService(c.ingextClient)
	count, err := svc.DeleteRuleGroup(group)
	if err != nil {
		c.Logger.Error("failed to delete eventwatch rule group", "group", group, "error", err)
		return 0, fmt.Errorf("failed to delete eventwatch rule group %q: %w", group, err)
	}
	return count, nil
}

// ListBehaviorFilter lists every behavior filter of the account via the
// behavior_filter_dao API.
func (c *Client) ListBehaviorFilter() ([]*model.BehaviorEventFilterT, error) {
	svc := fluencyAPI.NewEventWatchService(c.ingextClient)
	entries, err := svc.ListBehaviorFilter()
	if err != nil {
		c.Logger.Error("failed to list behavior filters", "error", err)
		return nil, fmt.Errorf("failed to list behavior filters: %w", err)
	}
	return entries, nil
}

// GetBehaviorFilter fetches a single behavior filter via the behavior_filter_dao API.
func (c *Client) GetBehaviorFilter(behaviorRule, name string) (*model.BehaviorEventFilterT, error) {
	svc := fluencyAPI.NewEventWatchService(c.ingextClient)
	entry, err := svc.GetBehaviorFilter(behaviorRule, name)
	if err != nil {
		c.Logger.Error("failed to get behavior filter", "behaviorRule", behaviorRule, "name", name, "error", err)
		return nil, fmt.Errorf("failed to get behavior filter %q/%q: %w", behaviorRule, name, err)
	}
	return entry, nil
}

// AddBehaviorFilter creates a new behavior filter via the behavior_filter_dao API.
func (c *Client) AddBehaviorFilter(entry *model.BehaviorEventFilterT) error {
	svc := fluencyAPI.NewEventWatchService(c.ingextClient)
	if err := svc.AddBehaviorFilter(entry); err != nil {
		c.Logger.Error("failed to add behavior filter", "error", err)
		return fmt.Errorf("failed to add behavior filter: %w", err)
	}
	return nil
}

// UpdateBehaviorFilter updates an existing behavior filter via the behavior_filter_dao API.
func (c *Client) UpdateBehaviorFilter(entry *model.BehaviorEventFilterT) error {
	svc := fluencyAPI.NewEventWatchService(c.ingextClient)
	if err := svc.UpdateBehaviorFilter(entry); err != nil {
		c.Logger.Error("failed to update behavior filter", "error", err)
		return fmt.Errorf("failed to update behavior filter: %w", err)
	}
	return nil
}

// ToggleBehaviorFilter flips the disabled state of a behavior filter via the
// behavior_filter_dao API.
func (c *Client) ToggleBehaviorFilter(behaviorRule, name string) error {
	svc := fluencyAPI.NewEventWatchService(c.ingextClient)
	if err := svc.ToggleBehaviorFilter(behaviorRule, name); err != nil {
		c.Logger.Error("failed to toggle behavior filter", "behaviorRule", behaviorRule, "name", name, "error", err)
		return fmt.Errorf("failed to toggle behavior filter %q/%q: %w", behaviorRule, name, err)
	}
	return nil
}

// DeleteBehaviorFilter removes a behavior filter via the behavior_filter_dao API.
func (c *Client) DeleteBehaviorFilter(behaviorRule, name string) error {
	svc := fluencyAPI.NewEventWatchService(c.ingextClient)
	if err := svc.DeleteBehaviorFilter(behaviorRule, name); err != nil {
		c.Logger.Error("failed to delete behavior filter", "behaviorRule", behaviorRule, "name", name, "error", err)
		return fmt.Errorf("failed to delete behavior filter %q/%q: %w", behaviorRule, name, err)
	}
	return nil
}
