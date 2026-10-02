package model

import (
	"encoding/json"
	"time"
)

// FacetConfig is one facet (term count) of a search facet profile.
type FacetConfig struct {
	Title string `json:"title,omitempty"`
	Field string `json:"field,omitempty"`
	Size  int    `json:"size,omitempty"`
}

// FacetConfigEntry is a search facet profile as stored by the platform's
// facet_config API (fsbv2 facet.FacetEntry). It is named apart from
// FacetEntry, which is a lake search result.
//
// Id is the storage key ("facet-xxxxxxxx") the platform assigns on add; ID is
// the release registry id carried in the Scripts repo file. The platform keys
// storage on Id, but the release sync matches entries by Name.
type FacetConfigEntry struct {
	Id            string          `json:"_id"`
	ID            int64           `json:"id"`
	Repository    string          `json:"repository,omitempty"`
	Group         string          `json:"group"`
	Name          string          `json:"name,omitempty"`
	Description   string          `json:"description,omitempty"`
	Disabled      bool            `json:"disabled"`
	Indexes       []string        `json:"indexes,omitempty"`
	Target        string          `json:"target,omitempty"`
	Query         string          `json:"query,omitempty"`
	LvdbQueryTree json.RawMessage `json:"lvdbQueryTree,omitempty"`
	Facets        []*FacetConfig  `json:"facets,omitempty"`
	FplObject     json.RawMessage `json:"fplObject,omitempty"`
	CreatedOn     time.Time       `json:"createdOn"`
	UpdatedOn     time.Time       `json:"updatedOn"`
}

type FacetConfigArgs struct {
	Id    string            `json:"id,omitempty"`
	Entry *FacetConfigEntry `json:"entry,omitempty"`
}

type FacetConfigRequest struct {
	Action string           `json:"action"`
	Args   *FacetConfigArgs `json:"args,omitempty"`
}

type FacetConfigListResponse struct {
	Entries []*FacetConfigEntry `json:"entries"`
}

type FacetConfigGetResponse struct {
	Entry *FacetConfigEntry `json:"entry"`
}
