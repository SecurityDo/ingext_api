package api

import (
	"fmt"

	"github.com/SecurityDo/ingext_api/client"
	"github.com/SecurityDo/ingext_api/model"
)

// FacetService calls facet_config, which manages the search facet profiles
// the UI offers for an index. With --gridaccount the provider proxies it to
// the tenant like any other call.
type FacetService struct {
	client *client.IngextClient
}

func NewFacetService(client *client.IngextClient) *FacetService {
	return &FacetService{client: client}
}

func (s *FacetService) call(action string, args *model.FacetConfigArgs, out interface{}) error {
	return ApiCall(s.client, "facet_config", &model.FacetConfigRequest{Action: action, Args: args}, out)
}

func (s *FacetService) List() ([]*model.FacetConfigEntry, error) {
	var resp model.FacetConfigListResponse
	if err := s.call("list", nil, &resp); err != nil {
		return nil, err
	}
	return resp.Entries, nil
}

// GetByName finds a profile by name, the key the release sync uses. It
// returns nil, nil when there is none.
func (s *FacetService) GetByName(name string) (*model.FacetConfigEntry, error) {
	entries, err := s.List()
	if err != nil {
		return nil, err
	}
	var found *model.FacetConfigEntry
	for _, e := range entries {
		if e.Name != name {
			continue
		}
		if found != nil {
			return nil, fmt.Errorf("facet %q is stored more than once (%s, %s)", name, found.Id, e.Id)
		}
		found = e
	}
	return found, nil
}

// Apply adds entry, or replaces the stored profile of the same name in place.
// facet_config's add always mints a new storage key, so adding blindly would
// leave a duplicate profile per run.
func (s *FacetService) Apply(entry *model.FacetConfigEntry) (created bool, err error) {
	if entry.Name == "" {
		return false, fmt.Errorf("facet entry has no name")
	}
	existing, err := s.GetByName(entry.Name)
	if err != nil {
		return false, err
	}
	if existing == nil {
		entry.Id = ""
		return true, s.call("add", &model.FacetConfigArgs{Entry: entry}, nil)
	}
	entry.Id = existing.Id
	if entry.CreatedOn.IsZero() {
		entry.CreatedOn = existing.CreatedOn
	}
	return false, s.call("update", &model.FacetConfigArgs{Id: existing.Id, Entry: entry}, nil)
}

func (s *FacetService) Delete(id string) error {
	return s.call("delete", &model.FacetConfigArgs{Id: id}, nil)
}
