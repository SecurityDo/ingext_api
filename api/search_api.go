package api

import (
	"github.com/SecurityDo/ingext_api/client"
	kqlModel "github.com/SecurityDo/ingext_api/kql/model"
)

// PlatformService provides helpers for calling platform_* endpoints.
type SearchService struct {
	client *client.IngextClient
}

// NewPlatformService constructs a PlatformService instance backed by the provided client.
func NewSearchService(client *client.IngextClient) *SearchService {
	return &SearchService{client: client}
}

func (s *SearchService) call(function string, payload interface{}, out interface{}) error {
	return ApiCall(s.client, function, payload, out)
}

type KQLSearchRequest struct {
	Index     string `json:"index,omitempty"`
	RangeFrom int64  `json:"range_from"`
	RangeTo   int64  `json:"range_to"`
	KQL       string `json:"kql,omitempty"`
}

func (s *SearchService) KQLSearch(kql string) (resp *kqlModel.KQLSearchResponse, err error) {
	request := &KQLSearchRequest{
		KQL: kql,
	}
	//var resp kqlModel.KQLSearchResponse
	if err := s.call("kql_search", request, &resp); err != nil {
		return nil, err
	}
	return resp, nil
}

// KQLValidateRequest is the payload for the kql_validate endpoint.
type KQLValidateRequest struct {
	KQL string `json:"kql"`
}

// KQLValidateResponse mirrors lakeAPI.KQLValidateResponse.
type KQLValidateResponse struct {
	OK             bool   `json:"ok"`
	Error          string `json:"error,omitempty"`
	Table          string `json:"table,omitempty"`
	TimeRangeFound bool   `json:"timeRangeFound,omitempty"`
}

// KQLValidate parses a KQL query on the search service without executing it.
// Returns a structured response; the CLI treats OK=false as a parser error.
func (s *SearchService) KQLValidate(kql string) (*KQLValidateResponse, error) {
	request := &KQLValidateRequest{KQL: kql}
	var resp KQLValidateResponse
	if err := s.call("kql_validate", request, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}
