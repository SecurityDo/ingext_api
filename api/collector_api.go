package api

import (
	"github.com/SecurityDo/ingext_api/client"
	"github.com/SecurityDo/ingext_api/model"
)

// CollectorService provides helpers for calling collector endpoints.
type CollectorService struct {
	client *client.IngextClient
}

// NewCollectorService constructs a CollectorService instance backed by the provided client.
func NewCollectorService(c *client.IngextClient) *CollectorService {
	return &CollectorService{client: c}
}

func (s *CollectorService) call(function string, payload interface{}, out interface{}) error {
	return ApiCall(s.client, function, payload, out)
}

type CollectorListResponse struct {
	Entries []*model.CollectorForWeb `json:"entries"`
}

// CollectorList calls the collector_list API with kargs {} and returns an array of CollectorForWeb.
func (s *CollectorService) CollectorList() ([]*model.CollectorForWeb, error) {
	var out CollectorListResponse
	if err := s.call("collector_list", map[string]interface{}{}, &out); err != nil {
		return nil, err
	}
	return out.Entries, nil
}

// CollectorImportKargs is the payload for collector_import (kargs: collector).
type CollectorImportKargs struct {
	Collector *model.CollectorT `json:"collector"`
}

// CollectorImport calls the collector_import API to recreate an existing collector
// on this site with its token preserved, which is how a collector is migrated from
// one site to another without re-enrolling the agent. The collector name must not
// already exist on the target site; createdOn is reset by the server.
func (s *CollectorService) CollectorImport(collector *model.CollectorT) error {
	payload := CollectorImportKargs{Collector: collector}
	return s.call("collector_import", payload, nil)
}

// CollectorStatusKargs is the payload for collector_status (kargs: collector + cargs).
type CollectorStatusKargs struct {
	Collector string                 `json:"collector"`
	Cargs     map[string]interface{} `json:"cargs,omitempty"`
}

// CollectorStatus calls the collector_status API with kargs {"collector":"name","cargs":{...}} and returns the response as map[string]interface{}.
func (s *CollectorService) CollectorStatus(collector string, cargs map[string]interface{}) (map[string]interface{}, error) {
	payload := CollectorStatusKargs{Collector: collector, Cargs: cargs}
	var out map[string]interface{}
	if err := s.call("get_system_status", payload, &out); err != nil {
		return nil, err
	}
	return out, nil
}
