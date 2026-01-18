package model

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
