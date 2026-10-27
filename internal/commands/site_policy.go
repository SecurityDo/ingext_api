package commands

import (
	"fmt"
	"strings"

	"github.com/SecurityDo/ingext_api/model"
	"github.com/spf13/cobra"
)

// Site policies decide which grid accounts a user or API token on a provider site
// can see: get_grid_accounts returns only the accounts named in a policy that
// lists the caller (or "*"). A user or token in no policy sees an empty grid.

var (
	spName        string
	spDescription string
	spAccounts    []string
	spUsers       []string
	spTokens      []string
)

var sitePolicyListCmd = &cobra.Command{
	Use:   "list-site-policy",
	Short: "List site policies (provider sites)",
	RunE: func(cmd *cobra.Command, args []string) error {
		policies, err := AppAPI.ListSitePolicy()
		if err != nil {
			cmd.PrintErrf("Error listing site policies: %v\n", err)
			return err
		}
		for _, p := range policies {
			fmt.Printf("SitePolicy: %s, Accounts: [%s], Users: [%s], Tokens: [%s]\n", p.Name,
				strings.Join(p.Accounts, ","), strings.Join(p.Users, ","), strings.Join(p.Tokens, ","))
		}
		return nil
	},
}

var sitePolicyAddCmd = &cobra.Command{
	Use:   "add-site-policy",
	Short: "Create a site policy (needs at least one user or token)",
	RunE: func(cmd *cobra.Command, args []string) error {
		entry := &model.SitePolicy{
			Name:        spName,
			Description: spDescription,
			Accounts:    spAccounts,
			Users:       spUsers,
			Tokens:      spTokens,
		}
		if err := AppAPI.AddSitePolicy(entry); err != nil {
			cmd.PrintErrf("Error adding site policy %s: %v\n", spName, err)
			return err
		}
		return nil
	},
}

var sitePolicyUpdateCmd = &cobra.Command{
	Use:   "update-site-policy",
	Short: "Add accounts, users and/or tokens to an existing site policy",
	RunE: func(cmd *cobra.Command, args []string) error {
		policies, err := AppAPI.ListSitePolicy()
		if err != nil {
			cmd.PrintErrf("Error listing site policies: %v\n", err)
			return err
		}
		var entry *model.SitePolicy
		for _, p := range policies {
			if p != nil && p.Name == spName {
				entry = p
				break
			}
		}
		if entry == nil {
			return fmt.Errorf("site policy %q not found", spName)
		}
		// The server replaces the stored policy with what we send, so merge into
		// the current one rather than sending only the additions.
		entry.Accounts = appendMissing(entry.Accounts, spAccounts)
		entry.Users = appendMissing(entry.Users, spUsers)
		entry.Tokens = appendMissing(entry.Tokens, spTokens)
		if cmd.Flags().Changed("description") {
			entry.Description = spDescription
		}
		if err := AppAPI.UpdateSitePolicy(entry); err != nil {
			cmd.PrintErrf("Error updating site policy %s: %v\n", spName, err)
			return err
		}
		return nil
	},
}

var sitePolicyDelCmd = &cobra.Command{
	Use:   "del-site-policy",
	Short: "Delete a site policy",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := AppAPI.DeleteSitePolicy(spName); err != nil {
			cmd.PrintErrf("Error deleting site policy %s: %v\n", spName, err)
			return err
		}
		return nil
	},
}

func appendMissing(have, add []string) []string {
	seen := make(map[string]bool, len(have))
	for _, s := range have {
		seen[s] = true
	}
	for _, s := range add {
		if !seen[s] {
			have = append(have, s)
			seen[s] = true
		}
	}
	return have
}

func init() {
	authCmd.AddCommand(sitePolicyListCmd, sitePolicyAddCmd, sitePolicyUpdateCmd, sitePolicyDelCmd)

	for _, c := range []*cobra.Command{sitePolicyAddCmd, sitePolicyUpdateCmd} {
		c.Flags().StringVar(&spName, "name", "", "Site policy name")
		c.Flags().StringVar(&spDescription, "description", "", "Description")
		c.Flags().StringSliceVar(&spAccounts, "accounts", nil, "Grid account names, comma separated (\"*\" = all)")
		c.Flags().StringSliceVar(&spUsers, "users", nil, "Usernames, comma separated")
		c.Flags().StringSliceVar(&spTokens, "tokens", nil, "API token names, comma separated")
		_ = c.MarkFlagRequired("name")
	}
	sitePolicyDelCmd.Flags().StringVar(&spName, "name", "", "Site policy name")
	_ = sitePolicyDelCmd.MarkFlagRequired("name")
}
