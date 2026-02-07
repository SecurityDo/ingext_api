package model

import "encoding/json"

// KQLSubSearchRequest is the request sent from coordinator to pods.
// Pods re-parse the KQL string (cheap) instead of receiving serialized AST.
type KQLSubSearchRequest struct {
	KQL            string `json:"kql"`
	Strategy       int    `json:"strategy"`
	TimeRangeFrom  int64  `json:"timeRangeFrom,omitempty"`
	TimeRangeTo    int64  `json:"timeRangeTo,omitempty"`
	TimeRangeFound bool   `json:"timeRangeFound,omitempty"`
	WhereOpIndices []int  `json:"whereOpIndices,omitempty"`
}

// KQLSubSearchJobResult is the result returned by a pod.
// GroupStates and FacetResults use json.RawMessage to avoid model→eval dependency.
type KQLSubSearchJobResult struct {
	TotalBytes   int64           `json:"totalBytes,omitempty"`
	TotalRows    int64           `json:"totalRows,omitempty"`
	Data         *DataSet        `json:"data,omitempty"`
	GroupStates  json.RawMessage `json:"groupStates,omitempty"`
	FacetResults json.RawMessage `json:"facetResults,omitempty"`
	Strategy     int             `json:"strategy"`
	SchemaNames  []string        `json:"schemaNames,omitempty"`
}

type KQLSearchResponse struct {
	Total      int64 `json:"total"`
	TotalBytes int64 `json:"totalBytes,omitempty"`
	//TaskID     string `json:"taskID,omitempty"`

	Data *DataSet `json:"data,omitempty"`

	RangeFrom int64   `json:"range_from,omitempty"`
	RangeTo   int64   `json:"range_to,omitempty"`
	Cost      float64 `json:"cost,omitempty"`
}
