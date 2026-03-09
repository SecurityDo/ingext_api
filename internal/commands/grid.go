package commands

import (
	"fmt"

	"github.com/SecurityDo/ingext_api/model"
	"github.com/spf13/cobra"
)

var (
	accountName        string
	accountRegion      string
	accountCluster     string
	accountSiteURL     string
	accountToken       string
	accountDisplayName string
)

var gridCmd = &cobra.Command{
	Use:   "grid",
	Short: "Grid management commands",
}

var listAccountCmd = &cobra.Command{
	Use:   "list-account",
	Short: "List all grid accounts",
	RunE: func(cmd *cobra.Command, args []string) error {
		resp, err := AppAPI.ListAccount()
		if err != nil {
			return err
		}
		if len(resp.Accounts) == 0 {
			cmd.Println("No accounts found.")
			return nil
		}
		for _, acct := range resp.Accounts {
			fmt.Printf("Name: %s, Region: %s, Cluster: %s, URL: %s\n", acct.Name, acct.Region, acct.Cluster, acct.URL)
		}
		return nil
	},
}

var addAccountCmd = &cobra.Command{
	Use:   "add-account",
	Short: "Add a grid SaaS account",
	RunE: func(cmd *cobra.Command, args []string) error {
		req := &model.GridAddSaasAccountRequest{
			Name:        accountName,
			Region:      accountRegion,
			Cluster:     accountCluster,
			SiteURL:     accountSiteURL,
			Token:       accountToken,
			DisplayName: accountDisplayName,
		}
		if err := AppAPI.AddSaasAccount(req); err != nil {
			return err
		}
		fmt.Printf("Account %q added successfully.\n", accountName)
		return nil
	},
}

func init() {
	RootCmd.AddCommand(gridCmd)
	gridCmd.AddCommand(listAccountCmd, addAccountCmd)

	addAccountCmd.Flags().StringVar(&accountName, "name", "", "Account name")
	addAccountCmd.Flags().StringVar(&accountRegion, "region", "", "Region")
	addAccountCmd.Flags().StringVar(&accountCluster, "cluster", "", "Cluster")
	addAccountCmd.Flags().StringVar(&accountSiteURL, "site-url", "", "Site URL")
	addAccountCmd.Flags().StringVar(&accountToken, "token", "", "Token")
	addAccountCmd.Flags().StringVar(&accountDisplayName, "display-name", "", "Display name")

	_ = addAccountCmd.MarkFlagRequired("name")
	_ = addAccountCmd.MarkFlagRequired("region")
	_ = addAccountCmd.MarkFlagRequired("cluster")
	_ = addAccountCmd.MarkFlagRequired("site-url")
	_ = addAccountCmd.MarkFlagRequired("token")
	_ = addAccountCmd.MarkFlagRequired("display-name")
}
