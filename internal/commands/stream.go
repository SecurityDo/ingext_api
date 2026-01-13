package commands

import (
	"fmt"

	model "github.com/SecurityDo/ingext_api/model"
	"github.com/spf13/cobra"
)

var (
	sourceType      string
	sinkType        string
	resourceID      string
	resourceName    string
	dataFormat      string
	dataCompression string
	integrationID   string // For associating with an integration
	url             string // For HEC sink
	token           string // For HEC sink
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

var delSourceCmd = &cobra.Command{
	Use:   "del-source",
	Short: "Delete a stream source",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("Deleting stream datasource...")

		err := AppAPI.DeleteDataSource(resourceID)
		if err != nil {
			return err
		}
		cmd.PrintErrln("Stream source deleted successfully: ", resourceID)
		return nil
	},
}

var listSourceCmd = &cobra.Command{
	Use:   "list-source",
	Short: "List all stream sources",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("Listing stream datasource...")

		entries, err := AppAPI.ListDataSource()
		if err != nil {
			return err
		}
		if len(entries) == 0 {
			cmd.Println("No stream source found.")
			return nil
		}
		for _, entry := range entries {
			cmd.PrintErrf("ID: %s, Name: %s, Type: %s\n", entry.ID, entry.Name, entry.Type)
		}
		return nil
	},
}

var delSinkCmd = &cobra.Command{
	Use:   "del-sink",
	Short: "Delete a stream sink",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("Deleting stream datasink...")

		err := AppAPI.DeleteDataSink(resourceID)
		if err != nil {
			return err
		}
		cmd.PrintErrln("Stream sink deleted successfully: ", resourceID)
		return nil
	},
}

var listSinkCmd = &cobra.Command{
	Use:   "list-sink",
	Short: "List all stream sinks",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("Listing stream datasink...")

		entries, err := AppAPI.ListDataSink()
		if err != nil {
			return err
		}
		if len(entries) == 0 {
			cmd.Println("No stream sink found.")
			return nil
		}
		for _, entry := range entries {
			cmd.PrintErrf("ID: %s, Name: %s, Type: %s\n", entry.ID, entry.Name, entry.Type)
		}
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
		sink := &model.DataSinkConfig{
			Type: sinkType,
			Name: resourceName,
			// Add other necessary fields for DataSourceConfig
			//Format: "json", // Example format, adjust as needed
			//Compression
		}

		if sinkType == "datalake" {
			if integrationID == "" {
				return fmt.Errorf("integration-id is required for datalake sink type")
			}
			sink.DataLake = &model.DataLakeSinkConfig{
				IntegrationID: integrationID,
			}
		} else if sinkType == "hec" {
			sink.Hec = &model.HecSinkConfig{
				URL:   url,
				Token: token,
				// Add necessary fields for HEC sink
			}
		} else if sinkType == "webhook" {
			sink.Webhook = &model.WebhookSinkConfig{
				URL: url,
				//Token: token,
				// Add necessary fields for HEC sink
			}
		}

		response, err := AppAPI.AddDataSink(sink)
		if err != nil {
			return err
		}

		cmd.PrintErrln("Stream sink added successfully: ", response.ID)
		cmd.Println(response.ID)

		return nil
	},
}

// ... Repeat for sink, router, connection ...

func init() {
	RootCmd.AddCommand(streamCmd)
	streamCmd.AddCommand(addSourceCmd, delSourceCmd, listSourceCmd, addSinkCmd, delSinkCmd, listSinkCmd) // Add del/update similarly

	addSourceCmd.Flags().StringVar(&sourceType, "source-type", "", "data source type: plugin, s3, hec, webhook ")
	addSourceCmd.Flags().StringVar(&resourceName, "name", "", "Name")
	addSourceCmd.Flags().StringVar(&dataFormat, "format", "json", "Data Format")
	addSourceCmd.Flags().StringVar(&dataCompression, "compression", "", "Data Compression")

	addSourceCmd.Flags().StringVar(&integrationID, "integration-id", "", "Integration ID")

	_ = addSourceCmd.MarkFlagRequired("source-type")
	_ = addSourceCmd.MarkFlagRequired("name")

	addSinkCmd.Flags().StringVar(&sourceType, "sink-type", "", "data sink type: datalake, hec, webhook")
	addSinkCmd.Flags().StringVar(&resourceName, "name", "", "Name")
	addSinkCmd.Flags().StringVar(&resourceName, "url", "", "URL for HEC sink")
	addSinkCmd.Flags().StringVar(&resourceName, "token", "", "Token")

	_ = addSinkCmd.MarkFlagRequired("sink-type")
	_ = addSinkCmd.MarkFlagRequired("name")

	delSourceCmd.Flags().StringVar(&resourceID, "id", "", "data source ID")
	_ = delSourceCmd.MarkFlagRequired("id")

	delSinkCmd.Flags().StringVar(&resourceID, "id", "", "data sink ID")
	_ = delSinkCmd.MarkFlagRequired("id")

}
