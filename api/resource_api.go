package api

import (
	"errors"

	"github.com/SecurityDo/ingext_api/client"
	"github.com/SecurityDo/ingext_api/model"
)

// PlatformService provides helpers for calling platform_* endpoints.
type ResourceService struct {
	client *client.IngextClient
}

func NewResourceService(client *client.IngextClient) *ResourceService {
	return &ResourceService{client: client}
}

func (s *ResourceService) call(function string, payload interface{}, out interface{}) error {
	return ApiCall(s.client, function, payload, out)
}

type ResourceSearchRequest struct {
	Options  *model.LakeFacetSearchOption `json:"options"`
	Customer string                       `json:"customer"`
	Resource string                       `json:"resource"`
}

func (s *ResourceService) Search(resourceType string, customer string) (resp *model.LakeSearchResponse, err error) {
	request := &ResourceSearchRequest{
		Options: &model.LakeFacetSearchOption{
			FetchLimit: 1000,
			Facets: &model.FacetsOption{
				Facets: []*model.FacetEntry{
					{
						Title: "Groups",
						Field: "@office365User.groups",
						Size:  20,
					},
				},
			},
		},
		Customer: customer,
		Resource: resourceType,
	}
	//var resp model.LakeSearchResponse
	if err := s.call("resource_search", request, &resp); err != nil {
		return nil, err
	}
	return resp, nil
}

// ResourceDumpDeleteRequest is the payload for the ingext_resource_dump_delete
// endpoint. Customer is the dump's customer segment — for a plugin datasource
// that is the plugin name the dump was written under.
type ResourceDumpDeleteRequest struct {
	Customer string `json:"customer"`
}

// ResourceDumpDeleteResponse reports how many (resourceType, customer) dump
// folders were purged.
type ResourceDumpDeleteResponse struct {
	Deleted int `json:"deleted"`
}

// DeleteResourceDump purges every resource dump belonging to one customer across
// all resource types, and returns the number of dump folders removed.
//
// This is the same cleanup that deleting the owning plugin datasource performs;
// the endpoint exists for the dumps of datasources that were already removed
// while that purge silently matched nothing. It is destructive and has no undo,
// so it needs data/manage, not data/write.
//
// An unknown customer is not an error — it deletes nothing and returns 0, which
// is the only way to tell a purge from a no-op. A partial failure returns an
// error and no count, even though some dumps were removed: re-running it is
// safe, since deleting an already-deleted dump is the no-op case.
func (s *ResourceService) DeleteResourceDump(customer string) (int, error) {
	if customer == "" {
		return 0, errors.New("resource dump customer is required")
	}
	request := &ResourceDumpDeleteRequest{Customer: customer}
	var resp ResourceDumpDeleteResponse
	if err := s.call("ingext_resource_dump_delete", request, &resp); err != nil {
		return 0, err
	}
	return resp.Deleted, nil
}
