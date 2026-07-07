package api

import (
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

type ElasticSearchRequest struct {
	Options *SimpleSearchOption `json:"options,omitempty"`
}
