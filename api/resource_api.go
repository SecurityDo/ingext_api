package api

import (
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
