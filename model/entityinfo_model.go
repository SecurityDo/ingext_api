package model

import "time"

// EntityTable is the config entry managed by the entityinfo_table_dao endpoint.
type EntityTable struct {
	ID         int64  `json:"id"`
	Repository string `json:"repository,omitempty"`
	Group      string `json:"group"`
	SyncPolicy string `json:"syncPolicy,omitempty"`

	Name        string   `json:"name,omitempty"`
	Description string   `json:"description,omitempty"`
	Fields      []string `json:"fields"`

	LookupType string `json:"lookupType"`

	FileContent string `json:"fileContent,omitempty"`

	// site or "", append, fork, "local"
	Mode string `json:"mode,omitempty"`

	UpdatedOn time.Time `json:"updatedOn"`
	CreatedOn time.Time `json:"createdOn"`
}

// EntityInfo is a single lookup row managed by the entityinfo_entry_dao endpoint.
type EntityInfo struct {
	Key   string                 `json:"key,omitempty"`
	Value string                 `json:"value,omitempty"`
	Info  map[string]interface{} `json:"info,omitempty"`
}
