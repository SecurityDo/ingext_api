package commands

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/SecurityDo/ingext_api/model"
	"github.com/spf13/cobra"
)

var (
	updID            string
	updName          string
	updType          string
	updNewName       string
	updDesc          string
	updUnset         []string
	updReplaceConfig bool
	updDryRun        bool
)

// integration update: read the entry, merge the given config keys into its
// config (or replace it with --replace-config), and save it. The platform's
// update never changes the secret; --secret replaces the whole secret through
// platform_integration_set_secret.
var integrationUpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update an integration's config, name or description (and optionally replace its secret)",
	Long: `Update an existing integration, picked by --id or by --name (with --integration when the
name is not unique). Config keys given with --config/--config-bool/--config-int/--config-json
are merged into the stored config; --unset-config removes keys; --replace-config starts from an
empty config instead. --secret replaces the WHOLE secret (every key it needs), and is skipped
when not given: an update never changes the stored secret otherwise.

  ingext integration update --name fluidia --integration EntraIDResponse \
    --config secretExpiry=2027-11-23T20:00:00Z`,
	RunE: func(cmd *cobra.Command, args []string) error {
		entry, err := findIntegration(updID, updName, updType)
		if err != nil {
			return err
		}

		config := map[string]any{}
		if !updReplaceConfig && len(entry.Config) > 0 && string(entry.Config) != "null" {
			if err := json.Unmarshal(entry.Config, &config); err != nil {
				return fmt.Errorf("stored config of %s is not a JSON object: %w", entry.ID, err)
			}
		}
		given, err := getConfigRaw()
		if err != nil {
			return fmt.Errorf("failed to get config: %w", err)
		}
		var changes map[string]any
		_ = json.Unmarshal(given, &changes)
		for k, v := range changes {
			config[k] = v
		}
		for _, k := range updUnset {
			delete(config, k)
		}
		changedKeys := make([]string, 0, len(changes)+len(updUnset))
		for k := range changes {
			changedKeys = append(changedKeys, k)
		}
		for _, k := range updUnset {
			changedKeys = append(changedKeys, "-"+k)
		}

		raw, _ := json.Marshal(config)
		entry.Config = raw
		entry.Secret = nil
		if updNewName != "" {
			entry.Name = updNewName
		}
		if cmd.Flags().Changed("description") {
			entry.Description = updDesc
		}
		configOrNothing := len(changedKeys) > 0 || updReplaceConfig || updNewName != "" || cmd.Flags().Changed("description")
		if !configOrNothing && len(secretParams) == 0 {
			return fmt.Errorf("nothing to update: give --config*, --unset-config, --new-name, --description or --secret")
		}

		cmd.PrintErrf("Integration %s (%s %s): keys %s\n", entry.ID, entry.Integration, entry.Name, strings.Join(changedKeys, ", "))
		if updDryRun {
			cmd.PrintErrln("dry run: new config follows; nothing saved")
			fmt.Println(string(raw))
			return nil
		}
		if configOrNothing {
			if err := AppAPI.UpdateIntegration(entry); err != nil {
				return err
			}
			cmd.PrintErrln("Integration updated: ", entry.ID)
		}
		if len(secretParams) > 0 {
			secret, err := getSecretRaw()
			if err != nil {
				return fmt.Errorf("failed to get secret: %w", err)
			}
			if err := AppAPI.SetIntegrationSecret(entry.ID, secret); err != nil {
				return err
			}
			cmd.PrintErrln("Integration secret replaced: ", entry.ID)
		}
		fmt.Println(entry.ID)
		return nil
	},
}

// findIntegration picks one integration by id, or by name (and type).
func findIntegration(id, name, typ string) (*model.Integration, error) {
	if id != "" {
		return AppAPI.GetIntegration(id)
	}
	if name == "" {
		return nil, fmt.Errorf("give --id, or --name (with --integration when the name is not unique)")
	}
	entries, err := AppAPI.ListIntegration()
	if err != nil {
		return nil, err
	}
	var hits []*model.Integration
	for _, e := range entries {
		if e.Name == name && (typ == "" || strings.EqualFold(e.Integration, typ)) {
			hits = append(hits, e)
		}
	}
	switch len(hits) {
	case 0:
		return nil, fmt.Errorf("no integration named %q%s", name, map[bool]string{true: " of type " + typ}[typ != ""])
	case 1:
		// The list may omit fields; read the entry itself.
		return AppAPI.GetIntegration(hits[0].ID)
	}
	ids := make([]string, 0, len(hits))
	for _, h := range hits {
		ids = append(ids, h.ID+" ("+h.Integration+")")
	}
	return nil, fmt.Errorf("%d integrations are named %q: %s; give --integration or --id", len(hits), name, strings.Join(ids, ", "))
}

func init() {
	integrationCmd.AddCommand(integrationUpdateCmd)
	f := integrationUpdateCmd.Flags()
	f.StringVar(&updID, "id", "", "Integration id")
	f.StringVar(&updName, "name", "", "Integration name (when --id is not given)")
	f.StringVar(&updType, "integration", "", "Integration type, to pick among integrations of the same name")
	f.StringVar(&updNewName, "new-name", "", "Rename the integration")
	f.StringVar(&updDesc, "description", "", "New description")
	f.StringToStringVarP(&configParams, "config", "", nil, "Config string keys to set (merged)")
	f.StringToStringVarP(&configBoolParams, "config-bool", "", nil, "Config boolean keys to set")
	f.StringToInt64VarP(&configIntParams, "config-int", "", nil, "Config int keys to set")
	f.StringArrayVar(&configJsonFlags, "config-json", []string{}, "Config JSON values to set (e.g. key=[1,2])")
	f.StringSliceVar(&updUnset, "unset-config", nil, "Config keys to remove")
	f.BoolVar(&updReplaceConfig, "replace-config", false, "Start from an empty config instead of merging")
	f.StringToStringVarP(&secretParams, "secret", "", nil, "Replace the whole secret with these keys (value @file reads a file)")
	f.BoolVar(&updDryRun, "dry-run", false, "Print the new config and save nothing")
}
