package commands

import (
	"fmt"
	"os"

	"github.com/SecurityDo/ingext_api/model"
	"github.com/spf13/cobra"
)

var (
	appName             string
	instanceName        string
	instanceDisplayName string
	inputParams         map[string]string
)

var applicationCmd = &cobra.Command{
	Use:   "application",
	Short: "Manage Ingext Application",
}

var appTemplateListCmd = &cobra.Command{
	Use:   "list-template",
	Short: "List installed application templates",
	RunE: func(cmd *cobra.Command, args []string) error {
		// Just call the global interface
		templates, err := AppAPI.ListAppTemplates()
		if err != nil {
			return err
		}

		if len(templates) == 0 {
			cmd.PrintErrf("No application template found.")
			return nil
		}

		for _, entry := range templates {
			cmd.PrintErrf("Name: %s, Description: %s\n", entry.Name, entry.Description)
		}
		return nil
	},
}

var appTemplateInstallCmd = &cobra.Command{
	Use:   "install",
	Short: "install application",
	RunE: func(cmd *cobra.Command, args []string) error {
		// Just call the global interface

		//config := make(map[string]string)
		var paras []*model.InputParameter
		for key, value := range configParams {
			if len(value) > 1 && value[0] == '@' {
				filePath := value[1:] // Remove '@'

				content, err := os.ReadFile(filePath)
				if err != nil {
					return fmt.Errorf("failed to read file for key '%s': %w", key, err)
				}
				paras = append(paras, &model.InputParameter{
					Name:  key,
					Value: string(content),
				})

				// Replace the file path with the actual file content
				//config[key] = string(content)
			} else {
				//config[key] = value
				paras = append(paras, &model.InputParameter{
					Name:  key,
					Value: string(value),
				})
			}
		}

		err := AppAPI.InstallAppInstance(appName, instanceName, instanceDisplayName, paras)
		if err != nil {
			return err
		}

		cmd.PrintErrf("Application %s:%s added successfully: ", appName, instanceName)
		//cmd.Println(id)
		return nil
	},
}

var appTemplateUnInstallCmd = &cobra.Command{
	Use:   "uninstall",
	Short: "uninstall application instance",
	RunE: func(cmd *cobra.Command, args []string) error {
		// Just call the global interface
		err := AppAPI.UnInstallAppInstance(appName, instanceName)
		if err != nil {
			return err
		}

		cmd.PrintErrln("Application instance added successfully")
		return nil
	},
}

func init() {
	RootCmd.AddCommand(applicationCmd)
	applicationCmd.AddCommand(appTemplateListCmd, appTemplateInstallCmd, appTemplateUnInstallCmd) // Add del/update similarly

	appTemplateInstallCmd.Flags().StringVar(&instanceName, "instance", "", "Instance name")
	appTemplateInstallCmd.Flags().StringVar(&appName, "app", "", "Application name")
	appTemplateInstallCmd.Flags().StringVar(&instanceDisplayName, "displayName", "", "Display name")

	appTemplateInstallCmd.Flags().StringToStringVarP(&configParams, "set", "", nil, "Configuration string type parameters")

	// Mark required
	_ = appTemplateInstallCmd.MarkFlagRequired("instance")
	_ = appTemplateInstallCmd.MarkFlagRequired("app")

	appTemplateUnInstallCmd.Flags().StringVar(&instanceName, "instance", "", "Instance name")
	appTemplateUnInstallCmd.Flags().StringVar(&appName, "app", "", "Application name")

	// Mark required
	_ = appTemplateUnInstallCmd.MarkFlagRequired("instance")
	_ = appTemplateUnInstallCmd.MarkFlagRequired("app")

	//streamAddCmd.AddCommand(streamAddSourceCmd)
	// Add other leaf commands: sink, router, connection
}
