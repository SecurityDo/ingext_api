package api

import (
	"errors"
	"fmt"

	"github.com/SecurityDo/ingext_api/client"
	kqlModel "github.com/SecurityDo/ingext_api/kql/model"
	"github.com/SecurityDo/ingext_api/model"
)

// PlatformService provides helpers for calling platform_* endpoints.
type SearchService struct {
	client *client.IngextClient
}

// NewPlatformService constructs a PlatformService instance backed by the provided client.
func NewSearchService(client *client.IngextClient) *SearchService {
	return &SearchService{client: client}
}

func (s *SearchService) call(function string, payload interface{}, out interface{}) error {
	return ApiCall(s.client, function, payload, out)
}

type KQLSearchRequest struct {
	Index     string `json:"index,omitempty"`
	RangeFrom int64  `json:"rangeFrom"`
	RangeTo   int64  `json:"rangeTo"`
	KQL       string `json:"kql,omitempty"`
}

func (s *SearchService) KQLSearch(kql string) (resp *kqlModel.KQLSearchResponse, err error) {
	request := &KQLSearchRequest{
		KQL: kql,
	}
	//var resp kqlModel.KQLSearchResponse
	if err := s.call("kql_search", request, &resp); err != nil {
		return nil, err
	}
	return resp, nil
}

// KQLValidateRequest is the payload for the kql_validate endpoint.
type KQLValidateRequest struct {
	KQL string `json:"kql"`
}

// KQLValidateResponse mirrors lakeAPI.KQLValidateResponse.
type KQLValidateResponse struct {
	OK             bool   `json:"ok"`
	Error          string `json:"error,omitempty"`
	Table          string `json:"table,omitempty"`
	TimeRangeFound bool   `json:"timeRangeFound,omitempty"`
}

// KQLValidate parses a KQL query on the search service without executing it.
// Returns a structured response; the CLI treats OK=false as a parser error.
func (s *SearchService) KQLValidate(kql string) (*KQLValidateResponse, error) {
	request := &KQLValidateRequest{KQL: kql}
	var resp KQLValidateResponse
	if err := s.call("kql_validate", request, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// --- lake_search: Elastic-style facet search over a datalake index ---

const (
	// LakeSearchMaxPagination is the platform's hard cap on FetchOffset+FetchLimit
	// for one lake search. Past it the call is rejected server-side.
	LakeSearchMaxPagination = 5000

	// lakeSearchDefaultFacetSize is applied to a facet that asks for no size. A
	// facet of size 0 is not "unlimited" on the platform - it returns no buckets
	// at all, which looks like a broken search rather than a missing default.
	lakeSearchDefaultFacetSize = 20

	// lakeSearchDefaultFetchLimit is applied to a search that asks for no hits.
	// A fetchLimit of 0 is not a facet-only search: it panics the search node
	// ("index out of range [0] with length 0") because the top-hits heap is
	// built with a capacity of 0 and then indexed. There is no way to ask this
	// endpoint for counts without hits.
	lakeSearchDefaultFetchLimit = 100

	// lakeSearchDefaultSortField / Order: lake_search has no server-side default
	// for the sort field, and an empty one fails the query parser.
	lakeSearchDefaultSortField = "@timestamp"
	lakeSearchDefaultSortOrder = "desc"
)

// LakeSearchRequest is the kargs payload for the lake_search endpoint.
//
// The endpoint also declares dataType, partition and dayIndex, but never reads
// them, so they are left off the wire.
type LakeSearchRequest struct {
	Index   string                       `json:"index,omitempty"`
	Options *model.LakeFacetSearchOption `json:"options"`
}

// LakeSearchOptions describes one lake facet search: a Lucene query and a time
// range, narrowed by term filters, with term facets counted over whatever
// matches.
type LakeSearchOptions struct {
	// Index is the datalake index to search, in "<datalake>-<index>" form:
	// "managed-Office365" is index "Office365" of the "managed" datalake. The
	// bare index name is what list_data_tables reports as a stream table. An
	// empty Index searches the "default" index.
	Index string

	// Query is a Lucene query string. Empty matches every event in range.
	Query string

	// RangeFrom and RangeTo bound the search in epoch milliseconds. Both are
	// required: without a range the platform scans the whole index.
	RangeFrom int64
	RangeTo   int64

	// Facets are the fields to count terms on. Each needs a Field; Size defaults
	// to 20 and Title, which the platform ignores, defaults to the field name.
	// The response is keyed by field. Do not assume it holds exactly the
	// requested set: a console response for this endpoint carried facets no one
	// asked for, so read the aggregations by name rather than by position.
	Facets []*model.FacetEntry

	// MustFilters and MustNotFilters narrow the search by exact terms. Terms
	// within one entry are OR'ed, and the entries are AND'ed.
	MustFilters    []*model.FilterEntry
	MustNotFilters []*model.FilterEntry

	// FetchLimit is how many hits to return; 0 means the default of 100. It
	// cannot be used to ask for counts alone - a fetchLimit of 0 reaches the
	// platform and panics the search node. FetchOffset+FetchLimit may not
	// exceed LakeSearchMaxPagination.
	FetchLimit  int
	FetchOffset int

	// SortField and SortOrder order the returned hits. They default to
	// "@timestamp" and "desc"; SortOrder is "asc" or "desc".
	SortField string
	SortOrder string
}

// LakeSearch runs a facet search against one datalake index and returns the
// matching hits together with the term counts.
//
// It is the search behind the console's event explorer: a Lucene query over a
// time range, must / must-not term filters, and a facet per field to count. The
// response also carries a "dateHistogram" aggregation - a fixed set of slots
// spanning the range, empty ones included.
//
// Total counts the events that matched, Filtered the ones that were read and
// discarded. Both are reported even when FetchLimit is 0.
func (s *SearchService) LakeSearch(opts *LakeSearchOptions) (*model.LakeSearchResponse, error) {
	request, err := opts.toRequest()
	if err != nil {
		return nil, err
	}
	var resp model.LakeSearchResponse
	if err := s.call("lake_search", request, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// toRequest validates the options and fills in the defaults the endpoint has
// none for. The checks are client-side on purpose: each of these fails deep
// inside the search service, where the error text says nothing about which
// field was wrong.
func (opts *LakeSearchOptions) toRequest() (*LakeSearchRequest, error) {
	if opts == nil {
		return nil, errors.New("lake search options are required")
	}
	if opts.RangeFrom <= 0 || opts.RangeTo <= 0 {
		return nil, errors.New("lake search requires rangeFrom and rangeTo in epoch milliseconds")
	}
	if opts.RangeFrom >= opts.RangeTo {
		return nil, fmt.Errorf("lake search range is empty: rangeFrom %d is not before rangeTo %d", opts.RangeFrom, opts.RangeTo)
	}
	if opts.FetchLimit < 0 || opts.FetchOffset < 0 {
		return nil, errors.New("lake search fetchLimit and fetchOffset cannot be negative")
	}
	// A fetchLimit of 0 panics the search node rather than returning counts
	// alone, so an unset limit becomes the default instead of being sent.
	fetchLimit := opts.FetchLimit
	if fetchLimit == 0 {
		fetchLimit = lakeSearchDefaultFetchLimit
	}
	if opts.FetchOffset+fetchLimit > LakeSearchMaxPagination {
		return nil, fmt.Errorf("lake search pagination limit is %d entries: fetchOffset %d + fetchLimit %d",
			LakeSearchMaxPagination, opts.FetchOffset, fetchLimit)
	}

	// Facets, and the slices inside it, are always sent: the endpoint walks them
	// without a nil check. The date facet is named so the histogram lands under
	// DateHistogramAggregation - the platform keys that aggregation by the name
	// of the first date facet, and an unnamed one keys it under "".
	facets := &model.FacetsOption{
		DateFacets:     []*model.DateFacetEntry{{Name: model.DateHistogramAggregation}},
		Facets:         []*model.FacetEntry{},
		MustFilters:    []*model.FilterEntry{},
		MustNotFilters: []*model.FilterEntry{},
	}
	for i, facet := range opts.Facets {
		if facet == nil || facet.Field == "" {
			return nil, fmt.Errorf("lake search facet %d: field is required", i)
		}
		entry := *facet
		if entry.Size <= 0 {
			entry.Size = lakeSearchDefaultFacetSize
		}
		if entry.Title == "" {
			entry.Title = entry.Field
		}
		facets.Facets = append(facets.Facets, &entry)
	}
	mustFilters, err := checkFilters("mustFilter", opts.MustFilters)
	if err != nil {
		return nil, err
	}
	facets.MustFilters = mustFilters
	mustNotFilters, err := checkFilters("mustNotFilter", opts.MustNotFilters)
	if err != nil {
		return nil, err
	}
	facets.MustNotFilters = mustNotFilters

	sortField := opts.SortField
	if sortField == "" {
		sortField = lakeSearchDefaultSortField
	}
	sortOrder := opts.SortOrder
	if sortOrder == "" {
		sortOrder = lakeSearchDefaultSortOrder
	}
	return &LakeSearchRequest{
		Index: opts.Index,
		Options: &model.LakeFacetSearchOption{
			SearchStr:   opts.Query,
			RangeFrom:   opts.RangeFrom,
			RangeTo:     opts.RangeTo,
			FetchLimit:  fetchLimit,
			FetchOffset: opts.FetchOffset,
			SortField:   sortField,
			SortOrder:   sortOrder,
			Facets:      facets,
		},
	}, nil
}

// checkFilters rejects the two filters that silently do nothing useful: one
// with no field, and one with no terms to match.
func checkFilters(kind string, filters []*model.FilterEntry) ([]*model.FilterEntry, error) {
	checked := make([]*model.FilterEntry, 0, len(filters))
	for i, filter := range filters {
		if filter == nil || filter.Field == "" {
			return nil, fmt.Errorf("lake search %s %d: field is required", kind, i)
		}
		if len(filter.Terms) == 0 {
			return nil, fmt.Errorf("lake search %s %d (%s): at least one term is required", kind, i, filter.Field)
		}
		checked = append(checked, filter)
	}
	return checked, nil
}
