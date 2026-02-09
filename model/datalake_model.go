package model

import "time"

type Datalake struct {
	// "managed" indicates if the lake is managed by the platform
	Name               string `json:"name,omitempty"`               // lake name
	Managed            bool   `json:"managed,omitempty"`            // local lake
	StorageDescription string `json:"storageDescription,omitempty"` // storage description
	IntegrationID      string `json:"integrationID,omitempty"`      // integration name
	Description        string `json:"description,omitempty"`        // lake description
	CreatedAt          string `json:"createdAt,omitempty"`          // creation time
}

type DatalakeIndex struct {
	Datalake      string `json:"datalake,omitempty"`      // lake name
	DatalakeIndex string `json:"datalakeIndex,omitempty"` // index name
	//IndexName 	   string                        `json:"indexName,omitempty"`     // index name, e.g. "my_index"
	Description        string `json:"description,omitempty"`        // index description
	SchemaName         string `json:"schemaName,omitempty"`         // schema name
	StorageDescription string `json:"storageDescription,omitempty"` // storage description
}

type SchemaEntry struct {
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt,omitempty"`
	// JSON Schema
	Content string `json:"content"`
}

type Field struct {
	Name string `json:"name"`
	// parquet
	Type           string   `json:"type"`
	ConvertedType  string   `json:"convertedtype,omitempty"`
	RepetitionType string   `json:"repetitiontype,omitempty"`
	Fields         []*Field `json:"fields,omitempty"`
	Nullable       bool     `json:"nullable,omitempty"`
}

type Table struct {
	Name   string   `json:"name"`
	Fields []*Field `json:"fields"`
}

type DatalakeIndexListRequest struct {
	// if lake is not set, return all indexes
	Lake string `json:"lake"`
}

type DatalakeIndexListResponse struct {
	// if lake is not set, return all indexes
	Entries []*DatalakeIndex `json:"entries"`
}

type DatalakeIndexAddRequest struct {
	// if lake is not set, return all indexes
	//Lake  string         `json:"lake"`
	Entry *DatalakeIndex `json:"entry"`
}
type DatalakeIndexDeleteRequest struct {
	// if lake is not set, return all indexes
	Lake  string `json:"lake"`
	Index string `json:"index"`
}
