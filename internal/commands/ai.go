package commands

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	aiURL   string
	aiToken string

	aiRegisterAccount     string
	aiRegisterDisplayName string
	aiUnregisterAccount   string
)

var aiCmd = &cobra.Command{
	Use:   "ai",
	Short: "AI token registration on the GPT API",
}

var aiRegisterCmd = &cobra.Command{
	Use:   "register",
	Short: "Add an API token on the GPT API (role admin)",
	RunE: func(cmd *cobra.Command, args []string) error {
		// API expects name and description; CLI uses --account and --display-name.
		token, err := AppAPI.GPTAIAuthService(aiURL, aiToken).AddToken(aiRegisterAccount, aiRegisterDisplayName, "tenant")
		if err != nil {
			cmd.PrintErrf("Error registering GPT AI token: %s %v\n", aiRegisterAccount, err)
			return err
		}
		fmt.Println(token)
		return nil
	},
}

var aiUnregisterCmd = &cobra.Command{
	Use:   "unregister",
	Short: "Delete an API token on the GPT API",
	RunE: func(cmd *cobra.Command, args []string) error {
		err := AppAPI.GPTAIAuthService(aiURL, aiToken).DeleteToken(aiUnregisterAccount)
		if err != nil {
			cmd.PrintErrf("Error unregistering GPT AI token: %s %v\n", aiUnregisterAccount, err)
			return err
		}
		return nil
	},
}

func init() {
	RootCmd.AddCommand(aiCmd)
	aiCmd.AddCommand(aiRegisterCmd, aiUnregisterCmd)

	aiCmd.PersistentFlags().StringVar(&aiURL, "url", "", "GPT API base URL")
	aiCmd.PersistentFlags().StringVar(&aiToken, "token", "", "Bearer token for the GPT API")
	_ = aiCmd.MarkPersistentFlagRequired("url")
	_ = aiCmd.MarkPersistentFlagRequired("token")

	aiRegisterCmd.Flags().StringVar(&aiRegisterAccount, "account", "", "Account name (sent as name to the API)")
	aiRegisterCmd.Flags().StringVar(&aiRegisterDisplayName, "display-name", "", "Display name (sent as description to the API)")
	_ = aiRegisterCmd.MarkFlagRequired("account")

	aiUnregisterCmd.Flags().StringVar(&aiUnregisterAccount, "account", "", "Account name (sent as name to the API)")
	_ = aiUnregisterCmd.MarkFlagRequired("account")
}

// IsAiTokenCommand reports whether cmd is ai register or ai unregister (leaf commands that supply their own URL and token).
func IsAiTokenCommand(cmd *cobra.Command) bool {
	if cmd == nil {
		return false
	}
	if cmd.Name() != "register" && cmd.Name() != "unregister" {
		return false
	}
	if cmd.Parent() == nil || cmd.Parent().Name() != "ai" {
		return false
	}
	return true
}
