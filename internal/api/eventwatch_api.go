package api

import (
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
