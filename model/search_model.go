package model

import "encoding/json"

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
	// either FacetBucketResult or HistogramBucketResult
	Buckets []interface{} `json:"buckets"`
}
