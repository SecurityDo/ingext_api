package commands

import "github.com/SecurityDo/ingext_api/internal/api"

// AppAPI is the global reference used by all subcommands.
// We use the interface type so it can be mocked for testing.
var AppAPI *api.Client

func init() {
	// Inject the real implementation by default.
	// In unit tests, you can overwrite this with a mock: commands.AppAPI = &MockAPI{}
	AppAPI = &api.Client{}
}
