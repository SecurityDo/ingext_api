package api

import (
	"github.com/SecurityDo/ingext_api/client"
	"github.com/SecurityDo/ingext_api/model"
)

// EntityInfoService provides helpers for calling the entityinfo_table_dao and
// entityinfo_entry_dao endpoints.
type EntityInfoService struct {
	client *client.IngextClient
}

// NewEntityInfoService constructs an EntityInfoService backed by the provided client.
func NewEntityInfoService(client *client.IngextClient) *EntityInfoService {
	return &EntityInfoService{client: client}
}

func (s *EntityInfoService) call(function string, payload interface{}, out interface{}) error {
	return ApiCall(s.client, function, payload, out)
}

const (
	entityInfoTableDAO = "entityinfo_table_dao"
	entityInfoEntryDAO = "entityinfo_entry_dao"
)

// --- entityinfo_table_dao ---

// EntityTableGetResponse wraps the single-entry response from the
// entityinfo_table_dao "get" action.
type EntityTableGetResponse struct {
	Entry *model.EntityTable `json:"entry"`
}

// EntityTableListResponse wraps the response from the entityinfo_table_dao
// "list" action.
type EntityTableListResponse struct {
	Entries []*model.EntityTable `json:"entries"`
}

// ListTable returns all entityinfo tables.
func (s *EntityInfoService) ListTable() ([]*model.EntityTable, error) {
	req := &GenericDAORequest[model.EntityTable]{
		Action: "list",
	}
	var resp EntityTableListResponse
	if err := s.call(entityInfoTableDAO, req, &resp); err != nil {
		return nil, err
	}
	return resp.Entries, nil
}

// GetTable fetches a single entityinfo table by name.
func (s *EntityInfoService) GetTable(name string) (*model.EntityTable, error) {
	req := &GenericDAORequest[model.EntityTable]{
		Action: "get",
		Args: &GenericDAORequestArgs[model.EntityTable]{
			Id: name,
		},
	}
	var resp EntityTableGetResponse
	if err := s.call(entityInfoTableDAO, req, &resp); err != nil {
		return nil, err
	}
	return resp.Entry, nil
}

// AddTable creates a new entityinfo table.
func (s *EntityInfoService) AddTable(entry *model.EntityTable) error {
	req := &GenericDAORequest[model.EntityTable]{
		Action: "add",
		Args: &GenericDAORequestArgs[model.EntityTable]{
			Entry: entry,
		},
	}
	return s.call(entityInfoTableDAO, req, nil)
}

// UpdateTable updates an existing entityinfo table.
func (s *EntityInfoService) UpdateTable(entry *model.EntityTable) error {
	req := &GenericDAORequest[model.EntityTable]{
		Action: "update",
		Args: &GenericDAORequestArgs[model.EntityTable]{
			Entry: entry,
		},
	}
	return s.call(entityInfoTableDAO, req, nil)
}

// DeleteTable removes an entityinfo table by name.
func (s *EntityInfoService) DeleteTable(name string) error {
	req := &GenericDAORequest[model.EntityTable]{
		Action: "delete",
		Args: &GenericDAORequestArgs[model.EntityTable]{
			Id: name,
		},
	}
	return s.call(entityInfoTableDAO, req, nil)
}

// --- entityinfo_entry_dao ---

// EntityEntryDAORequestArgs holds the id/entry payload for an entry DAO action.
type EntityEntryDAORequestArgs struct {
	Id    string            `json:"id,omitempty"`
	Entry *model.EntityInfo `json:"entry,omitempty"`
}

// EntityEntryDAORequest is the request envelope for entityinfo_entry_dao. Unlike
// the generic DAO envelope it carries the target entity (table) name.
type EntityEntryDAORequest struct {
	Entity string                     `json:"entity"`
	Action string                     `json:"action"`
	Args   *EntityEntryDAORequestArgs `json:"args,omitempty"`
}

// EntityEntryGetResponse wraps the single-entry response from the
// entityinfo_entry_dao "get" action.
type EntityEntryGetResponse struct {
	Entry *model.EntityInfo `json:"entry"`
}

// EntityEntryListResponse wraps the response from the entityinfo_entry_dao
// "list" action. The backend returns only the info maps of each row.
type EntityEntryListResponse struct {
	Entries []map[string]interface{} `json:"entries"`
}

// ListEntry returns the lookup rows (info maps) of the given entity table.
func (s *EntityInfoService) ListEntry(entity string) ([]map[string]interface{}, error) {
	req := &EntityEntryDAORequest{
		Entity: entity,
		Action: "list",
	}
	var resp EntityEntryListResponse
	if err := s.call(entityInfoEntryDAO, req, &resp); err != nil {
		return nil, err
	}
	return resp.Entries, nil
}

// GetEntry fetches a single lookup row from the entity table by key.
func (s *EntityInfoService) GetEntry(entity, key string) (*model.EntityInfo, error) {
	req := &EntityEntryDAORequest{
		Entity: entity,
		Action: "get",
		Args: &EntityEntryDAORequestArgs{
			Id: key,
		},
	}
	var resp EntityEntryGetResponse
	if err := s.call(entityInfoEntryDAO, req, &resp); err != nil {
		return nil, err
	}
	return resp.Entry, nil
}

// AddEntry creates a new lookup row in the entity table.
func (s *EntityInfoService) AddEntry(entity string, entry *model.EntityInfo) error {
	req := &EntityEntryDAORequest{
		Entity: entity,
		Action: "add",
		Args: &EntityEntryDAORequestArgs{
			Entry: entry,
		},
	}
	return s.call(entityInfoEntryDAO, req, nil)
}

// UpdateEntry updates an existing lookup row in the entity table.
func (s *EntityInfoService) UpdateEntry(entity string, entry *model.EntityInfo) error {
	req := &EntityEntryDAORequest{
		Entity: entity,
		Action: "update",
		Args: &EntityEntryDAORequestArgs{
			Entry: entry,
		},
	}
	return s.call(entityInfoEntryDAO, req, nil)
}

// DeleteEntry removes a lookup row from the entity table by key.
func (s *EntityInfoService) DeleteEntry(entity, key string) error {
	req := &EntityEntryDAORequest{
		Entity: entity,
		Action: "delete",
		Args: &EntityEntryDAORequestArgs{
			Id: key,
		},
	}
	return s.call(entityInfoEntryDAO, req, nil)
}
