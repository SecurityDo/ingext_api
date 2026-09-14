package commands

import (
	"encoding/json"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"
)

var (
	notifName            string
	notifAction          string
	notifTo              []string
	notifCc              []string
	notifChannels        []string
	notifIntegrationName string
	notifActionsFilter   string
	notifActionsJSON     bool
)

var notificationCmd = &cobra.Command{
	Use:   "notification",
	Short: "Notification endpoint management commands",
}

// printEndpointAdded reports an add. The dao answers with no id today, so the
// id is only mentioned when there is one rather than printed as "(id: )".
func printEndpointAdded(kind, name, id string) {
	if id == "" {
		fmt.Printf("%s notification endpoint %q added successfully.\n", kind, name)
		return
	}
	fmt.Printf("%s notification endpoint %q added successfully (id: %s).\n", kind, name, id)
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

var notificationGetCmd = &cobra.Command{
	Use:   "get",
	Short: "Get a notification endpoint by name",
	RunE: func(cmd *cobra.Command, args []string) error {
		entry, err := AppAPI.NotificationGet(notifName)
		if err != nil {
			return err
		}
		jsonBytes, err := json.MarshalIndent(entry, "", "  ")
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
		printEndpointAdded("Email", notifName, id)
		return nil
	},
}

var notificationAddSlackCmd = &cobra.Command{
	Use:   "add-slack",
	Short: "Add a Slack notification endpoint",
	RunE: func(cmd *cobra.Command, args []string) error {
		id, err := AppAPI.NotificationAddSlack(notifName, notifAction, notifIntegrationName, notifChannels)
		if err != nil {
			return err
		}
		printEndpointAdded("Slack", notifName, id)
		return nil
	},
}

var notificationActionsCmd = &cobra.Command{
	Use:   "actions",
	Short: "List the actions a notification endpoint can use",
	Long: "List the FPL actions that deliver through a notification endpoint.\n" +
		"An endpoint's --action must name one of these, and its integration must\n" +
		"match the action's.",
	RunE: func(cmd *cobra.Command, args []string) error {
		actions, err := AppAPI.NotificationListActions(notifActionsFilter)
		if err != nil {
			return err
		}
		if len(actions) == 0 {
			cmd.Println("No notification actions found.")
			return nil
		}
		if notifActionsJSON {
			jsonBytes, err := json.MarshalIndent(actions, "", "  ")
			if err != nil {
				return fmt.Errorf("failed to marshal response: %w", err)
			}
			fmt.Println(string(jsonBytes))
			return nil
		}
		// The script body is the bulk of each action; show it only with --json.
		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "NAME\tINTEGRATION\tGROUP\tDESCRIPTION")
		for _, a := range actions {
			fmt.Fprintf(w, "%s\t%s\t%s\t%s\n", a.Name, a.ActionConfig.Integration, a.Group, a.Description)
		}
		return w.Flush()
	},
}

func init() {
	RootCmd.AddCommand(notificationCmd)
	notificationCmd.AddCommand(notificationListCmd, notificationGetCmd, notificationDeleteCmd,
		notificationAddEmailCmd, notificationAddSlackCmd, notificationActionsCmd)

	notificationGetCmd.Flags().StringVar(&notifName, "name", "", "Endpoint name")
	_ = notificationGetCmd.MarkFlagRequired("name")

	notificationDeleteCmd.Flags().StringVar(&notifName, "name", "", "Endpoint name")
	_ = notificationDeleteCmd.MarkFlagRequired("name")

	notificationAddEmailCmd.Flags().StringVar(&notifName, "name", "", "Endpoint name")
	notificationAddEmailCmd.Flags().StringVar(&notifAction, "action", "Generic_Email_Action", "Action name")
	notificationAddEmailCmd.Flags().StringArrayVar(&notifTo, "to", []string{}, "To email addresses")
	notificationAddEmailCmd.Flags().StringArrayVar(&notifCc, "cc", []string{}, "Cc email addresses")
	_ = notificationAddEmailCmd.MarkFlagRequired("name")
	_ = notificationAddEmailCmd.MarkFlagRequired("to")

	notificationAddSlackCmd.Flags().StringVar(&notifName, "name", "", "Endpoint name")
	notificationAddSlackCmd.Flags().StringVar(&notifAction, "action", "", "Action name (see 'ingext notification actions --integration Slack')")
	notificationAddSlackCmd.Flags().StringArrayVar(&notifChannels, "channel", []string{}, "Slack channel to post to (repeatable)")
	notificationAddSlackCmd.Flags().StringVar(&notifIntegrationName, "integration-name", "", "Name of the configured Slack integration holding the token")
	_ = notificationAddSlackCmd.MarkFlagRequired("name")
	_ = notificationAddSlackCmd.MarkFlagRequired("action")
	_ = notificationAddSlackCmd.MarkFlagRequired("channel")

	notificationActionsCmd.Flags().StringVar(&notifActionsFilter, "integration", "", "Only show actions for this integration (Email, Slack)")
	notificationActionsCmd.Flags().BoolVar(&notifActionsJSON, "json", false, "Print the full action entries, including script text")
}
