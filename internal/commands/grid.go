package commands

import (
	"fmt"

	"github.com/spf13/cobra"
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

func init() {
	RootCmd.AddCommand(gridCmd)
	gridCmd.AddCommand(listAccountCmd)
}
