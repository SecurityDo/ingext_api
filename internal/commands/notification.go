package commands

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
)

var (
	notifName   string
	notifAction string
	notifTo     []string
	notifCc     []string
)

var notificationCmd = &cobra.Command{
	Use:   "notification",
	Short: "Notification endpoint management commands",
}

var notificationListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all notification endpoints",
	RunE: func(cmd *cobra.Command, args []string) error {
		endpoints, err := AppAPI.NotificationList()
		if err != nil {
			return err
		}
		if len(endpoints) == 0 {
			cmd.Println("No notification endpoints found.")
			return nil
		}
		jsonBytes, err := json.MarshalIndent(endpoints, "", "  ")
		if err != nil {
			return fmt.Errorf("failed to marshal response: %w", err)
		}
		fmt.Println(string(jsonBytes))
		return nil
	},
}

var notificationDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete a notification endpoint",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := AppAPI.NotificationDelete(notifName); err != nil {
			return err
		}
		fmt.Printf("Notification endpoint %q deleted successfully.\n", notifName)
		return nil
	},
}

var notificationAddEmailCmd = &cobra.Command{
	Use:   "add-email",
	Short: "Add an email notification endpoint",
	RunE: func(cmd *cobra.Command, args []string) error {
		id, err := AppAPI.NotificationAddEmail(notifName, notifAction, notifTo, notifCc)
		if err != nil {
			return err
		}
		fmt.Printf("Email notification endpoint %q added successfully (id: %s).\n", notifName, id)
		return nil
	},
}

func init() {
	RootCmd.AddCommand(notificationCmd)
	notificationCmd.AddCommand(notificationListCmd, notificationDeleteCmd, notificationAddEmailCmd)

	notificationDeleteCmd.Flags().StringVar(&notifName, "name", "", "Endpoint name")
	_ = notificationDeleteCmd.MarkFlagRequired("name")

	notificationAddEmailCmd.Flags().StringVar(&notifName, "name", "", "Endpoint name")
	notificationAddEmailCmd.Flags().StringVar(&notifAction, "action", "Generic_Email_Action", "Action name")
	notificationAddEmailCmd.Flags().StringArrayVar(&notifTo, "to", []string{}, "To email addresses")
	notificationAddEmailCmd.Flags().StringArrayVar(&notifCc, "cc", []string{}, "Cc email addresses")
	_ = notificationAddEmailCmd.MarkFlagRequired("name")
	_ = notificationAddEmailCmd.MarkFlagRequired("to")
}
