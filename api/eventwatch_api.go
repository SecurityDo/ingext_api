package api

import (
	"fmt"
	"strings"

	"github.com/SecurityDo/ingext_api/client"
	"github.com/SecurityDo/ingext_api/model"
)

// EventWatchService provides helpers for calling eventwatch/overview endpoints.
type EventWatchService struct {
	client *client.IngextClient
}

// NewEventWatchService constructs an EventWatchService instance backed by the provided client.
func NewEventWatchService(client *client.IngextClient) *EventWatchService {
	return &EventWatchService{client: client}
}

func (s *EventWatchService) call(function string, payload interface{}, out interface{}) error {
	return ApiCall(s.client, function, payload, out)
}

// SummarySearch calls /api/ds/overview_summary_search with the given search string and time range.
func (s *EventWatchService) SummarySearch(searchString string, rangeFrom, rangeTo int64) (*model.ElasticSearchResult, error) {
	req := &ElasticSearchRequest{
		Options: &SimpleSearchOption{
			SearchStr:   searchString,
			RangeFrom:   rangeFrom,
			RangeTo:     rangeTo,
			RangeField:  "from",
			FetchLimit:  100,
			FetchOffset: 0,
			SortField:   "to",
			SortOrder:   "desc",
			Facets: &FacetsOption{
				Facets:         []*FacetEntry{},
				MustFilters:    []*FilterEntry{},
				MustNotFilters: []*FilterEntry{},
			},
		},
	}
	var resp model.ElasticSearchResult
	if err := s.call("behavior_summary_search", req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// TimelineSearch calls /api/ds/fsm_behavior_search with the given search string and time range.
func (s *EventWatchService) TimelineSearch(searchString string, rangeFrom, rangeTo int64) (*model.ElasticSearchResult, error) {
	req := &ElasticSearchRequest{
		Options: &SimpleSearchOption{
			SearchStr:   searchString,
			RangeFrom:   rangeFrom,
			RangeTo:     rangeTo,
			RangeField:  "timestamp",
			FetchLimit:  100,
			FetchOffset: 0,
			SortField:   "timestamp",
			SortOrder:   "desc",
			Facets: &FacetsOption{
				Facets:         []*FacetEntry{},
				MustFilters:    []*FilterEntry{},
				MustNotFilters: []*FilterEntry{},
			},
		},
	}
	var resp model.ElasticSearchResult
	if err := s.call("fsm_behavior_search", req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// RuleSearch calls /api/ds/eventwatch_bucket_search with the given search string (no time range).
func (s *EventWatchService) RuleSearch(searchString string) (*model.ElasticSearchResult, error) {
	req := &ElasticSearchRequest{
		Options: &SimpleSearchOption{
			SearchStr:   searchString,
			FetchLimit:  100,
			FetchOffset: 0,
			SortField:   "name",
			SortOrder:   "asc",
			Facets: &FacetsOption{
				Facets:         []*FacetEntry{},
				MustFilters:    []*FilterEntry{},
				MustNotFilters: []*FilterEntry{},
			},
		},
	}
	var resp model.ElasticSearchResult
	if err := s.call("eventwatch_bucket_search", req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// --- EventWatch rule DAO (backed by the eventwatch_bucket_dao endpoint) ---

// EventWatchRuleGetResponse wraps the single-entry response returned by the
// eventwatch_bucket_dao "get" action.
type EventWatchRuleGetResponse struct {
	Entry *model.EventWatchBucket `json:"entry"`
}

// EventWatchRuleListResponse wraps the response returned by the
// eventwatch_bucket_dao "list" action.
type EventWatchRuleListResponse struct {
	Entries []*model.EventWatchBucket `json:"entries"`
	Tags    []string                  `json:"tags"`
	Groups  []string                  `json:"groups"`
}

// ListRule returns all eventwatch rules along with the distinct tags and groups.
func (s *EventWatchService) ListRule() (*EventWatchRuleListResponse, error) {
	req := &GenericDAORequest[model.EventWatchBucket]{
		Action: "list",
	}
	var resp EventWatchRuleListResponse
	if err := s.call("eventwatch_bucket_dao", req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// GetRule fetches a single eventwatch rule by name. When name is empty the first
// available rule is returned.
func (s *EventWatchService) GetRule(name string) (*model.EventWatchBucket, error) {
	req := &GenericDAORequest[model.EventWatchBucket]{
		Action: "get",
		Args: &GenericDAORequestArgs[model.EventWatchBucket]{
			Id: name,
		},
	}
	var resp EventWatchRuleGetResponse
	if err := s.call("eventwatch_bucket_dao", req, &resp); err != nil {
		return nil, err
	}
	return resp.Entry, nil
}

// AddRule creates a new eventwatch rule.
func (s *EventWatchService) AddRule(entry *model.EventWatchBucket) error {
	req := &GenericDAORequest[model.EventWatchBucket]{
		Action: "add",
		Args: &GenericDAORequestArgs[model.EventWatchBucket]{
			Entry: entry,
		},
	}
	return s.call("eventwatch_bucket_dao", req, nil)
}

// UpdateRule updates an existing eventwatch rule.
func (s *EventWatchService) UpdateRule(entry *model.EventWatchBucket) error {
	req := &GenericDAORequest[model.EventWatchBucket]{
		Action: "update",
		Args: &GenericDAORequestArgs[model.EventWatchBucket]{
			Entry: entry,
		},
	}
	return s.call("eventwatch_bucket_dao", req, nil)
}

// ToggleRule flips the disabled state of an eventwatch rule identified by name.
func (s *EventWatchService) ToggleRule(name string) error {
	req := &GenericDAORequest[model.EventWatchBucket]{
		Action: "toggle",
		Args: &GenericDAORequestArgs[model.EventWatchBucket]{
			Id: name,
		},
	}
	return s.call("eventwatch_bucket_dao", req, nil)
}

// DeleteRule removes an eventwatch rule identified by name.
func (s *EventWatchService) DeleteRule(name string) error {
	req := &GenericDAORequest[model.EventWatchBucket]{
		Action: "delete",
		Args: &GenericDAORequestArgs[model.EventWatchBucket]{
			Id: name,
		},
	}
	return s.call("eventwatch_bucket_dao", req, nil)
}

// --- Behavior filter DAO (behavior_filter_dao) ---

// behaviorFilterDAO is the /api/ds function backing the behavior filter CRUD calls.
const behaviorFilterDAO = "behavior_filter_dao"

// BehaviorFilterID builds the DAO id of a behavior filter. A filter is keyed by
// its owning behavior rule plus its name, and the two travel to the server joined
// by a slash in the single "id" arg -- the server splits on the first slash, so a
// name may contain slashes but a behavior rule may not.
func BehaviorFilterID(behaviorRule, name string) string {
	return behaviorRule + "/" + name
}

// behaviorFilterArgs builds the args for the id-addressed actions, rejecting the
// inputs the server cannot resolve back into a key.
func behaviorFilterArgs(behaviorRule, name string) (*GenericDAORequestArgs[model.BehaviorEventFilterT], error) {
	if behaviorRule == "" {
		return nil, fmt.Errorf("behavior filter behaviorRule is required")
	}
	if name == "" {
		return nil, fmt.Errorf("behavior filter name is required")
	}
	if strings.Contains(behaviorRule, "/") {
		return nil, fmt.Errorf("behavior filter behaviorRule %q must not contain a slash", behaviorRule)
	}
	return &GenericDAORequestArgs[model.BehaviorEventFilterT]{
		Id: BehaviorFilterID(behaviorRule, name),
	}, nil
}

// BehaviorFilterGetResponse wraps the single-entry response returned by the
// behavior_filter_dao "get" action.
type BehaviorFilterGetResponse struct {
	Entry *model.BehaviorEventFilterT `json:"entry"`
}

// BehaviorFilterListResponse wraps the response returned by the
// behavior_filter_dao "list" action.
type BehaviorFilterListResponse struct {
	Entries []*model.BehaviorEventFilterT `json:"entries"`
}

// ListBehaviorFilter returns every behavior filter of the account, across all
// behavior rules.
func (s *EventWatchService) ListBehaviorFilter() ([]*model.BehaviorEventFilterT, error) {
	req := &GenericDAORequest[model.BehaviorEventFilterT]{Action: "list"}
	var resp BehaviorFilterListResponse
	if err := s.call(behaviorFilterDAO, req, &resp); err != nil {
		return nil, err
	}
	return resp.Entries, nil
}

// GetBehaviorFilter fetches a single behavior filter. Both the owning behavior
// rule and the filter name are required; the server reports a missing filter as
// an error rather than an empty result.
func (s *EventWatchService) GetBehaviorFilter(behaviorRule, name string) (*model.BehaviorEventFilterT, error) {
	args, err := behaviorFilterArgs(behaviorRule, name)
	if err != nil {
		return nil, err
	}
	req := &GenericDAORequest[model.BehaviorEventFilterT]{Action: "get", Args: args}
	var resp BehaviorFilterGetResponse
	if err := s.call(behaviorFilterDAO, req, &resp); err != nil {
		return nil, err
	}
	return resp.Entry, nil
}

// AddBehaviorFilter creates a new behavior filter. The entry needs Name,
// BehaviorRule, and at least one entry in Filters; the server stamps CreatedOn
// and UpdatedOn and rejects a name already used under the same behavior rule.
func (s *EventWatchService) AddBehaviorFilter(entry *model.BehaviorEventFilterT) error {
	req := &GenericDAORequest[model.BehaviorEventFilterT]{
		Action: "add",
		Args:   &GenericDAORequestArgs[model.BehaviorEventFilterT]{Entry: entry},
	}
	return s.call(behaviorFilterDAO, req, nil)
}

// UpdateBehaviorFilter replaces an existing behavior filter, identified by the
// entry's BehaviorRule and Name rather than by an id. Changing either field
// addresses a different filter, so the update fails as "behavior filter does not
// exist" instead of renaming: add the new filter and delete the old one.
func (s *EventWatchService) UpdateBehaviorFilter(entry *model.BehaviorEventFilterT) error {
	req := &GenericDAORequest[model.BehaviorEventFilterT]{
		Action: "update",
		Args:   &GenericDAORequestArgs[model.BehaviorEventFilterT]{Entry: entry},
	}
	return s.call(behaviorFilterDAO, req, nil)
}

// ToggleBehaviorFilter flips the disabled state of a behavior filter.
func (s *EventWatchService) ToggleBehaviorFilter(behaviorRule, name string) error {
	args, err := behaviorFilterArgs(behaviorRule, name)
	if err != nil {
		return err
	}
	req := &GenericDAORequest[model.BehaviorEventFilterT]{Action: "toggle", Args: args}
	return s.call(behaviorFilterDAO, req, nil)
}

// DeleteBehaviorFilter removes a behavior filter.
func (s *EventWatchService) DeleteBehaviorFilter(behaviorRule, name string) error {
	args, err := behaviorFilterArgs(behaviorRule, name)
	if err != nil {
		return err
	}
	req := &GenericDAORequest[model.BehaviorEventFilterT]{Action: "delete", Args: args}
	return s.call(behaviorFilterDAO, req, nil)
}

// EventWatchDeleteGroupRequest is the payload for the eventwatch_bucket_delete_group endpoint.
type EventWatchDeleteGroupRequest struct {
	Group string `json:"group"`
}

// EventWatchDeleteGroupResponse is the response returned by the
// eventwatch_bucket_delete_group endpoint, reporting how many rules were removed.
type EventWatchDeleteGroupResponse struct {
	Count int `json:"count"`
}

// DeleteRuleGroup removes an entire eventwatch rule group via the
// eventwatch_bucket_delete_group endpoint and returns the number of rules deleted.
func (s *EventWatchService) DeleteRuleGroup(group string) (int, error) {
	req := &EventWatchDeleteGroupRequest{Group: group}
	var resp EventWatchDeleteGroupResponse
	if err := s.call("eventwatch_bucket_delete_group", req, &resp); err != nil {
		return 0, err
	}
	return resp.Count, nil
}

type ElasticSearchRequest struct {
	Options *SimpleSearchOption `json:"options,omitempty"`
}
