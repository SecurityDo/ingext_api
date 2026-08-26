package model

import (
	"encoding/json"
	"fmt"
)

type LakeFacetSearchOption struct {
	SearchStr string `json:"searchStr,omitempty"`
	//LVDBQuery json.RawMessage `json:"lvdbQuery,omitempty"`
	//FplText string               `json:"fplText,omitempty"`
	RangeFrom int64 `json:"range_from"`
	RangeTo   int64 `json:"range_to"`

	// for top N search
	FetchLimit  int    `json:"fetchLimit,omitempty"`
	FetchOffset int    `json:"fetchOffset"`
	SortField   string `json:"sortField,omitempty"`
	SortOrder   string `json:"sortOrder,omitempty"`

	/*
		DateFacetName  string `json:"dateFacetName"`
		DateFacetField string `json:"dateFacetField"`
		DateFacetKey   string `json:"dateFacetKey"`
		DateFacetValue string `json:"dateFacetValue"`
	*/

	Facets *FacetsOption `json:"facets,omitempty"`

	//Aggregates     map[string]*AggregateRequest `json:"aggs,omitempty"`
	//APIVersion     uint64                       `json:"APIVersion"`
	//Fields         []string                     `json:"fields,omitempty"`
}

type FacetsOption struct {
	DateFacets     []*DateFacetEntry `json:"dateFacets"`
	Facets         []*FacetEntry     `json:"facets"`
	MustFilters    []*FilterEntry    `json:"mustFilters"`
	MustNotFilters []*FilterEntry    `json:"mustNotFilters"`
}
type DateFacetEntry struct {
	Name string `json:"name,omitempty"`
	//Interval    string      `json:"interval,omitempty"`
	Key   string `json:"key,omitempty"`
	Value string `json:"value,omitempty"`
	//FilterField string      `json:"filterField,omitempty"`
	//FilterTerm  interface{} `json:"filterTerm,omitempty"`
}
type FacetEntry struct {
	Title string `json:"title"`
	Field string `json:"field"`
	Order string `json:"order"`
	Size  int    `json:"size"`
}
type FilterEntry struct {
	Field string   `json:"field"`
	Terms []string `json:"terms"`
	//FilterType string        `json:"filterType"`
}

type LakeSearchResponse struct {
	Took         int64                       `json:"took"`
	Hits         *LegacySearchHits           `json:"hits,omitempty"`
	Aggregations map[string]*SearchAggregate `json:"aggregations,omitempty"`
	Query        json.RawMessage             `json:"query"`
	Terms        []string                    `json:"terms,omitempty"`
	////GroupBys     *GroupByAggregations        `json:"groupBys,omitempty"`
	// APIVersion == 2
	////Aggs map[string]*AggregateResponse `json:"aggs,omitempty"`
	// APIVersion == 3
	TaskID string `json:"taskID,omitempty"`
	////Table     *fplmodel.FplTable `json:"table,omitempty"`
	RangeFrom int64   `json:"range_from,omitempty"`
	RangeTo   int64   `json:"range_to,omitempty"`
	Cost      float64 `json:"cost,omitempty"`
	Total     int64   `json:"total,omitempty"`
	Filtered  int64   `json:"filtered,omitempty"`
}

type LegacySearchHits struct {
	Total         uint64             `json:"total"`
	SortFieldType string             `json:"sortFieldType"`
	Hits          []*LegacySearchHit `json:"hits"`

	Ascend bool `json:"ascend"` // true for ascending, false for descending
	Limit  int  `json:"limit"`  // the limit of hits, 0 means no limit
}
type LegacySearchHit struct {
	Source json.RawMessage `json:"_source"`
	ID     string          `json:"_id,omitempty"` // for legacy search, we use _id to store the doc ID
}
type SearchAggregate struct {
	TokenEntity string
	name        string
	//
	//Buckets     []map[string]interface{} `json:"buckets"`
	// Buckets are kept as raw JSON because the two bucket shapes are not
	// interchangeable: a term facet keys on a string, the date histogram keys on
	// an epoch-millisecond int64. Decoding into interface{} would turn that
	// timestamp into a float64 and lose its exact value. Use FacetBuckets or
	// HistogramBuckets to read them.
	Buckets []json.RawMessage `json:"buckets"`
}

// DateHistogramAggregation is the aggregation name the platform reserves for the
// date histogram over the search range. Every other key in
// LakeSearchResponse.Aggregations is a term facet, named after the facet field.
const DateHistogramAggregation = "dateHistogram"

// FacetBucket is one term of a facet aggregation.
type FacetBucket struct {
	Key      string `json:"key"`
	DocCount int64  `json:"doc_count"`
}

// HistogramBucket is one slot of the date histogram. Key is the epoch
// millisecond the slot starts at; the platform always returns a fixed number of
// slots spanning the search range, including the empty ones.
type HistogramBucket struct {
	Key      int64           `json:"key"`
	DocCount int64           `json:"doc_count"`
	SubAggs  []*SubAggResult `json:"subAggs,omitempty"`
}

// SubAggResult is a named value computed inside one histogram slot.
type SubAggResult struct {
	Name  string `json:"name"`
	Value int64  `json:"value"`
}

// FacetBuckets decodes the aggregation as a term facet. It fails on the
// dateHistogram aggregation, whose keys are numbers - read that one with
// HistogramBuckets.
func (a *SearchAggregate) FacetBuckets() ([]*FacetBucket, error) {
	if a == nil || len(a.Buckets) == 0 {
		return nil, nil
	}
	buckets := make([]*FacetBucket, 0, len(a.Buckets))
	for i, raw := range a.Buckets {
		var bucket FacetBucket
		if err := json.Unmarshal(raw, &bucket); err != nil {
			return nil, fmt.Errorf("facet bucket %d: %w", i, err)
		}
		buckets = append(buckets, &bucket)
	}
	return buckets, nil
}

// HistogramBuckets decodes the aggregation as the date histogram.
func (a *SearchAggregate) HistogramBuckets() ([]*HistogramBucket, error) {
	if a == nil || len(a.Buckets) == 0 {
		return nil, nil
	}
	buckets := make([]*HistogramBucket, 0, len(a.Buckets))
	for i, raw := range a.Buckets {
		var bucket HistogramBucket
		if err := json.Unmarshal(raw, &bucket); err != nil {
			return nil, fmt.Errorf("histogram bucket %d: %w", i, err)
		}
		buckets = append(buckets, &bucket)
	}
	return buckets, nil
}
