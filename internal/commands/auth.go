package commands

import (
	//"fmt"

	"strings"

	"github.com/spf13/cobra"
	//"github.com/SecurityDo/ingext_api/model"
)

var (
	authName        string
	authDisplayName string
	authRole        string
	authOrg         string
)

// Parent command
var authCmd = &cobra.Command{
	Use:   "auth",
	Short: "Manage users",
}

// Verbs
var userAddCmd = &cobra.Command{
	Use:   "add-user",
	Short: "Add a user",
	RunE: func(cmd *cobra.Command, args []string) error {
		//cmd.PrintErrf("Adding user: %s (Role: %s)\n", authName, authRole)
		err := AppAPI.AddUser(authName, authDisplayName, authRole, authOrg)
		if err != nil {
			cmd.PrintErrf("Error adding user: %s %v\n", authName, err)
			return err
		}
		return nil
	},
}

var userDelCmd = &cobra.Command{
	Use:   "del-user",
	Short: "Delete a uer",
	RunE: func(cmd *cobra.Command, args []string) error {
		//cmd.PrintErrf("Adding user: %s (Role: %s)\n", authName, authRole)
		err := AppAPI.DeleteUser(authName)
		if err != nil {
			cmd.PrintErrf("Error deleting user: %s %v\n", authName, err)
			return err
		}
		return nil
	},
}

var userListCmd = &cobra.Command{
	Use:   "list-user",
	Short: "List users",
	RunE: func(cmd *cobra.Command, args []string) error {
		//cmd.PrintErrf("Adding user: %s (Role: %s)\n", authName, authRole)
		users, err := AppAPI.ListUser()
		if err != nil {
			cmd.PrintErrf("Error listing user: %v\n", err)
			return err
		}
		for _, user := range users {
			cmd.Printf("User: %s, Display Name: %s, Role: %s, Org: %s\n", user.Username, user.FirstName, strings.Join(user.Roles, ","), user.Organization)
		}
		return nil
	},
}

/*
// Nouns (Token)
var authAddTokenCmd = &cobra.Command{
	Use:   "token",
	Short: "Add a new token",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Adding token for: %s\n", authName)
	},
}*/

// ... Repeat similar logic for del/update user/token ...

func init() {
	RootCmd.AddCommand(authCmd)
	authCmd.AddCommand(userAddCmd, userDelCmd, userListCmd)

	// Add 'user' and 'token' to 'add'
	//authAddCmd.AddCommand(authAddUserCmd, authDelUserCmd)

	// Add flags to the leaf commands (or persistent flags to the verbs)
	userAddCmd.Flags().StringVar(&authName, "name", "", "Name of the user")
	userAddCmd.Flags().StringVar(&authDisplayName, "displayName", "", "Display name")
	userAddCmd.Flags().StringVar(&authRole, "role", "", "Role (admin|analyst)")
	userAddCmd.Flags().StringVar(&authOrg, "org", "ingext", "Organization")

	// Mark required
	_ = userAddCmd.MarkFlagRequired("name")
	_ = userAddCmd.MarkFlagRequired("role")
	//_ = authAddUserCmd.MarkFlagRequired("org")

	userDelCmd.Flags().StringVar(&authName, "name", "", "Name of the user")
	_ = userDelCmd.MarkFlagRequired("name")

}
