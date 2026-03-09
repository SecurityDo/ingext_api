package commands

import (
	//"fmt"

	"fmt"
	"strings"

	"github.com/spf13/cobra"
	//"github.com/SecurityDo/ingext_api/model"
)

var (
	authName        string
	authDisplayName string
	authRole        string
	authOrg         string
	authOAuth       string
	authDescription string
	authUsername     string
	authPolicy      string
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
		err := AppAPI.AddUser(authName, authDisplayName, authRole, authOrg, authOAuth)
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

var tokenAddCmd = &cobra.Command{
	Use:   "add-token",
	Short: "Add a token",
	RunE: func(cmd *cobra.Command, args []string) error {
		//cmd.PrintErrf("Adding user: %s (Role: %s)\n", authName, authRole)
		token, err := AppAPI.AddToken(authName, authDescription, authRole)
		if err != nil {
			cmd.PrintErrf("Error adding token: %s %v\n", authName, err)
			return err
		}
		fmt.Println(token)
		return nil
	},
}

var tokenDelCmd = &cobra.Command{
	Use:   "del-token",
	Short: "Delete a token",
	RunE: func(cmd *cobra.Command, args []string) error {
		//cmd.PrintErrf("Adding user: %s (Role: %s)\n", authName, authRole)
		err := AppAPI.DeleteToken(authName)
		if err != nil {
			cmd.PrintErrf("Error deleting token: %s %v\n", authName, err)
			return err
		}
		return nil
	},
}

var tokenListCmd = &cobra.Command{
	Use:   "list-token",
	Short: "List tokens",
	RunE: func(cmd *cobra.Command, args []string) error {
		//cmd.PrintErrf("Adding user: %s (Role: %s)\n", authName, authRole)
		tokens, err := AppAPI.ListToken()
		if err != nil {
			cmd.PrintErrf("Error listing tokens: %v\n", err)
			return err
		}
		for _, token := range tokens {
			cmd.PrintErrf("Token: %s, Display Name: %s, Role: %s\n", token.Name, token.Description, strings.Join(token.Roles, ","))
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

var setUserSitePolicyCmd = &cobra.Command{
	Use:   "set-user-site-policy",
	Short: "Set site policy for a user",
	RunE: func(cmd *cobra.Command, args []string) error {
		err := AppAPI.SetUserSitePolicy(authUsername, authPolicy)
		if err != nil {
			cmd.PrintErrf("Error setting site policy for user %s: %v\n", authUsername, err)
			return err
		}
		return nil
	},
}

func init() {
	RootCmd.AddCommand(authCmd)
	authCmd.AddCommand(userAddCmd, userDelCmd, userListCmd, tokenAddCmd, tokenDelCmd, tokenListCmd, setUserSitePolicyCmd)

	// Add 'user' and 'token' to 'add'
	//authAddCmd.AddCommand(authAddUserCmd, authDelUserCmd)

	// Add flags to the leaf commands (or persistent flags to the verbs)
	userAddCmd.Flags().StringVar(&authName, "name", "", "Name of the user")
	userAddCmd.Flags().StringVar(&authDisplayName, "displayName", "", "Display name")
	userAddCmd.Flags().StringVar(&authRole, "role", "", "Role (admin|analyst)")
	userAddCmd.Flags().StringVar(&authOrg, "org", "ingext", "Organization")
	userAddCmd.Flags().StringVar(&authOAuth, "oauth", "", "OAuth provider (Azure|Google)")

	// Mark required
	_ = userAddCmd.MarkFlagRequired("name")
	_ = userAddCmd.MarkFlagRequired("role")
	//_ = authAddUserCmd.MarkFlagRequired("org")

	userDelCmd.Flags().StringVar(&authName, "name", "", "Name of the user")
	_ = userDelCmd.MarkFlagRequired("name")

	tokenAddCmd.Flags().StringVar(&authName, "name", "", "Name of the token")
	tokenAddCmd.Flags().StringVar(&authDescription, "description", "", "Description of the token")
	tokenAddCmd.Flags().StringVar(&authRole, "role", "", "Role (admin|analyst)")

	// Mark required
	_ = tokenAddCmd.MarkFlagRequired("name")
	_ = tokenAddCmd.MarkFlagRequired("role")

	tokenDelCmd.Flags().StringVar(&authName, "name", "", "Name of the token")

	// Mark required
	_ = tokenDelCmd.MarkFlagRequired("name")

	setUserSitePolicyCmd.Flags().StringVar(&authUsername, "username", "", "Username")
	setUserSitePolicyCmd.Flags().StringVar(&authPolicy, "policy", "", "Policy name")
	_ = setUserSitePolicyCmd.MarkFlagRequired("username")
	_ = setUserSitePolicyCmd.MarkFlagRequired("policy")
}
