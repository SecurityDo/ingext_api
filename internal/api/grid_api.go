package api

import (
	"fmt"
	"sort"
	"strings"

	ingextAPI "github.com/SecurityDo/ingext_api/api"
	"github.com/SecurityDo/ingext_api/model"
)

func (c *Client) ListAccount() (resp *model.ListFluencyAccountsResponse, err error) {
	service := ingextAPI.NewGridService(c.ingextClient)

	resp, err = service.ListAccount()
	if err != nil {
		c.Logger.Error("list account error", "error", err)
		return nil, fmt.Errorf("list account error: %w", err)
	}
	return resp, nil
}

// VerifyGridAccount checks that the connected site can actually proxy requests to
// the given tenant account before any command runs against it.
//
// Only grid manager sites implement the "?gridaccount=" proxy. Plain ingext sites
// parse the query parameter and then ignore it, so the call silently executes
// against the site's own account — a delete aimed at a tenant would hit the wrong
// data. Manager sites serve /api/grid; plain sites return 404 for it, which is how
// the two are told apart here.
func (c *Client) VerifyGridAccount(account string) error {
	if account == "" {
		return nil
	}
	service := ingextAPI.NewGridService(c.ingextClient)
	resp, err := service.ProbeAccounts()
	if err != nil {
		if strings.Contains(err.Error(), "404") {
			return fmt.Errorf("site %s does not support --gridaccount: it is not a grid manager, "+
				"so the tenant would be ignored and the command would run against the site's own account. "+
				"Connect to the tenant directly (INGEXT_SITE_URL/INGEXT_TOKEN or --site) instead",
				c.SiteURL())
		}
		return fmt.Errorf("failed to verify --gridaccount %q against %s: %w", account, c.SiteURL(), err)
	}

	names := make([]string, 0, len(resp.Accounts))
	for _, acct := range resp.Accounts {
		if acct == nil {
			continue
		}
		if acct.Name == account {
			if acct.Disabled {
				return fmt.Errorf("grid account %q is disabled on %s", account, c.SiteURL())
			}
			return nil
		}
		names = append(names, acct.Name)
	}
	sort.Strings(names)
	return fmt.Errorf("grid account %q not found on %s; available accounts: %s",
		account, c.SiteURL(), strings.Join(names, ", "))
}

func (c *Client) AddSaasAccount(req *model.GridAddSaasAccountRequest) error {
	service := ingextAPI.NewGridService(c.ingextClient)

	if err := service.AddSaasAccount(req); err != nil {
		c.Logger.Error("add saas account error", "error", err)
		return fmt.Errorf("add saas account error: %w", err)
	}
	return nil
}

func (c *Client) DeleteSaasAccount(req *model.GridDeleteSaasAccountRequest) error {
	service := ingextAPI.NewGridService(c.ingextClient)

	if err := service.DeleteSaasAccount(req); err != nil {
		c.Logger.Error("delete saas account error", "error", err)
		return fmt.Errorf("delete saas account error: %w", err)
	}
	return nil
}
