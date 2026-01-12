package commands

import (
	"fmt"

	model "github.com/SecurityDo/ingext_api/model"
	"github.com/spf13/cobra"
)

var (
	sourceType      string
	resourceName    string
	dataFormat      string
	dataCompression string
	integrationID   string // For associating with an integration
)

var streamCmd = &cobra.Command{
	Use:   "stream",
	Short: "Manage streams",
}

var addSourceCmd = &cobra.Command{
	Use:   "add-source",
	Short: "Add a stream source",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("Adding stream datasource...")

		source := &model.DataSourceConfig{
			Type: sourceType,
			Name: resourceName,
			// Add other necessary fields for DataSourceConfig
			Format: "json", // Example format, adjust as needed
			//Compression
		}

		if sourceType == "plugin" {
			if integrationID == "" {
				return fmt.Errorf("integration-id is required for plugin source type")
			}
			source.Plugin = &model.PluginSourceConfig{
				ID: integrationID,
			}
		}

		response, err := AppAPI.AddDataSource(source)
		if err != nil {
			return err
		}

		cmd.PrintErrln("Stream source added successfully: ", response.ID)
		if response.URL != "" {
			cmd.PrintErrln("Access URL:", response.URL)
		}
		if len(response.Secret) > 0 {
			b, _ := response.Secret.MarshalJSON()
			cmd.PrintErrln("Secret:", string(b))
		}
		cmd.Println(response.ID)
		return nil
	},
}

// Example leaf command: source
var addSinkCmd = &cobra.Command{
	Use:   "add-sink",
	Short: "Add a stream sink",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("Adding stream source...")
		// Just call the global interface
		_, err := AppAPI.AddDataSink(nil)
		if err != nil {
			return err
		}

		fmt.Println("Stream source added successfully.")
		return nil
	},
}

// ... Repeat for sink, router, connection ...

func init() {
	RootCmd.AddCommand(streamCmd)
	streamCmd.AddCommand(addSourceCmd, addSinkCmd) // Add del/update similarly

	addSourceCmd.Flags().StringVar(&sourceType, "source-type", "", "data source type: plugin, s3, hec, webhook ")
	addSourceCmd.Flags().StringVar(&resourceName, "name", "", "Name")
	addSourceCmd.Flags().StringVar(&dataFormat, "format", "json", "Data Format")
	addSourceCmd.Flags().StringVar(&dataCompression, "compression", "", "Data Compression")

	addSourceCmd.Flags().StringVar(&integrationID, "integration-id", "", "Integration ID")

	_ = addSourceCmd.MarkFlagRequired("source-type")
	_ = addSourceCmd.MarkFlagRequired("name")

	//streamAddCmd.AddCommand(streamAddSourceCmd)
	// Add other leaf commands: sink, router, connection
}
