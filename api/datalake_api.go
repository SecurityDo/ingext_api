package api

import (
	"github.com/SecurityDo/ingext_api/client"
	"github.com/SecurityDo/ingext_api/model"
)

// PlatformService provides helpers for calling platform_* endpoints.
type DatalakeService struct {
	client *client.IngextClient
}

// NewPlatformService constructs a PlatformService instance backed by the provided client.
func NewDatalakeService(client *client.IngextClient) *DatalakeService {
	return &DatalakeService{client: client}
}

func (s *DatalakeService) call(function string, payload interface{}, out interface{}) error {
	return ApiCall(s.client, function, payload, out)
}

func (s *DatalakeService) ListDatalake() (entries []*model.Datalake, err error) {
	request := &GenericDAORequest[model.Datalake]{
		Action: "list",
	}
	var resp GenericDaoListResponse[model.Datalake]
	if err := s.call("ingext_datalake_dao", request, &resp); err != nil {
		return nil, err
	}
	return resp.Entries, nil

}

func (s *DatalakeService) AddDatalake(name string, managed bool, integrationID string) (err error) {
	if managed {
		name = "managed"
	}
	request := &GenericDAORequest[model.Datalake]{
		Action: "add",
		Args: &GenericDAORequestArgs[model.Datalake]{
			Entry: &model.Datalake{
				Name:          name,
				Managed:       managed,
				IntegrationID: integrationID,
			},
		},
	}
	//var resp GenericDaoListResponse[model.Datalake]
	if err := s.call("ingext_datalake_dao", request, nil); err != nil {
		return err
	}
	return nil

}

func (s *DatalakeService) AddDatalakeIndex(datalake string, index string, schema string) (err error) {

	request := &model.DatalakeIndexAddRequest{
		Entry: &model.DatalakeIndex{
			Datalake:      datalake,
			DatalakeIndex: index,
			SchemaName:    schema,
		},
	}
	if err := s.call("ingext_datalake_index_add", request, nil); err != nil {
		return err
	}
	return nil

}

func (s *DatalakeService) DeleteDatalakeIndex(datalake string, index string) (err error) {

	request := &model.DatalakeIndexDeleteRequest{
		Lake:  datalake,
		Index: index,
	}
	if err := s.call("ingext_datalake_index_delete", request, nil); err != nil {
		return err
	}
	return nil

}

func (s *DatalakeService) ListDatalakeIndex(datalake string) (entries []*model.DatalakeIndex, err error) {

	request := &model.DatalakeIndexListRequest{
		Lake: datalake,
	}
	var resp model.DatalakeIndexListResponse
	if err := s.call("ingext_datalake_index_list", request, &resp); err != nil {
		return nil, err
	}
	return resp.Entries, nil

}

func (s *DatalakeService) ListSchema() (entries []*model.SchemaEntry, err error) {
	request := &GenericDAORequest[model.SchemaEntry]{
		Action: "list",
	}
	var resp GenericDaoListResponse[model.SchemaEntry]
	if err := s.call("ingext_datalake_schema_dao", request, &resp); err != nil {
		return nil, err
	}
	return resp.Entries, nil

}

func (s *DatalakeService) UpdateSchema(name, description string, content string) error {
	request := &GenericDAORequest[model.SchemaEntry]{
		Action: "update",
		Args: &GenericDAORequestArgs[model.SchemaEntry]{
			Entry: &model.SchemaEntry{
				Name:        name,
				Description: description,
				Content:     content,
			},
		},
	}
	if err := s.call("ingext_datalake_schema_dao", request, nil); err != nil {
		return err
	}
	return nil
}

func (s *DatalakeService) DeleteSchema(name string) error {
	request := &GenericDAORequest[model.SchemaEntry]{
		Action: "delete",
		Args: &GenericDAORequestArgs[model.SchemaEntry]{
			Id: name,
		},
	}
	if err := s.call("ingext_datalake_schema_dao", request, nil); err != nil {
		return err
	}
	return nil
}

func (s *DatalakeService) AddSchema(name, description string, content string) (err error) {
	request := &GenericDAORequest[model.SchemaEntry]{
		Action: "add",
		Args: &GenericDAORequestArgs[model.SchemaEntry]{
			Entry: &model.SchemaEntry{
				Name:        name,
				Description: description,
				Content:     content,
			},
		},
	}
	var resp GenericDaoListResponse[model.SchemaEntry]
	if err := s.call("ingext_datalake_schema_dao", request, &resp); err != nil {
		return err
	}
	return nil

}
