package commands

import (
	"encoding/json"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/SecurityDo/ingext_api/model"
	"github.com/spf13/cobra"
)

var (
	datalake         string
	index            string
	managed          bool
	integration      string
	schema           string
	schemaName       string
	schemaDesc       string
	schemaFile       string
)

var lakeCmd = &cobra.Command{
	Use:   "datalake",
	Short: "Manage datalake",
}

var lakeListCmd = &cobra.Command{
	Use:   "list",
	Short: "List datalake",
	RunE: func(cmd *cobra.Command, args []string) error {
		//	cmd.PrintErrf("Adding datalake %s\n", datalake)
		lakes, err := AppAPI.ListDatalakes()
		if err != nil {
			return err
		}

		if len(lakes) == 0 {
			cmd.PrintErrf("No datalake found.")
			return nil
		}

		for _, entry := range lakes {
			cmd.PrintErrf("Name: %s, Description: %s\n", entry.Name, entry.StorageDescription)
		}

		//cmd.PrintErrln("Datalake added successfully")
		return nil
	},
}

var lakeAddCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a datalake",
	RunE: func(cmd *cobra.Command, args []string) error {
		cmd.PrintErrf("Adding datalake %s  %t\n", datalake, managed)
		err := AppAPI.AddDatalake(datalake, managed, integration)
		if err != nil {
			return err
		}
		cmd.PrintErrln("Datalake added successfully")
		return nil
	},
}

var lakeAddIndexCmd = &cobra.Command{
	Use:   "add-index",
	Short: "Add an index to the datalake",
	RunE: func(cmd *cobra.Command, args []string) error {
		cmd.PrintErrf("Adding index %s to datalake %s\n", index, datalake)
		err := AppAPI.AddDatalakeIndex(datalake, index, schema)
		if err != nil {
			return err
		}
		cmd.PrintErrln("Index added successfully")
		return nil
	},
}

var lakeListIndexCmd = &cobra.Command{
	Use:   "list-index",
	Short: "list all index ine one datalake",
	RunE: func(cmd *cobra.Command, args []string) error {
		//cmd.PrintErrf("Adding index % to datalake %s\n", index, datalake)
		entries, err := AppAPI.ListDatalakeIndex(datalake)
		if err != nil {
			return err
		}
		if len(entries) == 0 {
			cmd.Println("No datalake index found.")
			return nil
		}

		for _, entry := range entries {
			cmd.PrintErrf("Datalake: %s, Index: %s, Description: %s\n", entry.Datalake, entry.DatalakeIndex, entry.StorageDescription)
		}
		return nil

	},
}

var lakeDeleteIndexCmd = &cobra.Command{
	Use:   "del-index",
	Short: "delete index from one datalake",
	RunE: func(cmd *cobra.Command, args []string) error {
		//cmd.PrintErrf("Adding index % to datalake %s\n", index, datalake)
		err := AppAPI.DeleteDatalakeIndex(datalake, index)
		if err != nil {
			return err
		}
		cmd.PrintErrln("Index deleted successfully")
		return nil

	},
}

var lakeAddSchemaCmd = &cobra.Command{
	Use:   "add-schema",
	Short: "Add a schema to the datalake",
	RunE: func(cmd *cobra.Command, args []string) error {
		content, err := os.ReadFile(schemaFile)
		if err != nil {
			return fmt.Errorf("failed to read schema file %s: %w", schemaFile, err)
		}
		cmd.PrintErrf("Adding schema %s\n", schemaName)
		err = AppAPI.AddSchema(schemaName, schemaDesc, string(content))
		if err != nil {
			return err
		}
		cmd.PrintErrln("Schema added successfully")
		return nil
	},
}

var lakeListSchemaCmd = &cobra.Command{
	Use:   "list-schema",
	Short: "List all schemas",
	RunE: func(cmd *cobra.Command, args []string) error {
		schemas, err := AppAPI.ListSchemas()
		if err != nil {
			return err
		}
		if len(schemas) == 0 {
			cmd.PrintErrln("No schemas found.")
			return nil
		}

		for _, entry := range schemas {
			var table model.Table
			if err := json.Unmarshal([]byte(entry.Content), &table); err != nil {
				cmd.PrintErrf("Schema: %s (failed to decode: %v)\n", entry.Name, err)
				continue
			}

			cmd.PrintErrf("Schema: %s\n", entry.Name)
			if entry.Description != "" {
				cmd.PrintErrf("Description: %s\n", entry.Description)
			}

			w := tabwriter.NewWriter(os.Stderr, 0, 0, 2, ' ', 0)
			fmt.Fprintln(w, "NAME\tPHYSICAL TYPE\tLOGICAL TYPE\tNULLABLE")
			fmt.Fprintln(w, "----\t-------------\t------------\t--------")
			printFields(w, table.Fields, "")
			w.Flush()
			cmd.PrintErrln()
		}
		return nil
	},
}

var lakeUpdateSchemaCmd = &cobra.Command{
	Use:   "update-schema",
	Short: "Update a schema",
	RunE: func(cmd *cobra.Command, args []string) error {
		content, err := os.ReadFile(schemaFile)
		if err != nil {
			return fmt.Errorf("failed to read schema file %s: %w", schemaFile, err)
		}
		cmd.PrintErrf("Updating schema %s\n", schemaName)
		err = AppAPI.UpdateSchema(schemaName, schemaDesc, string(content))
		if err != nil {
			return err
		}
		cmd.PrintErrln("Schema updated successfully")
		return nil
	},
}

var lakeDeleteSchemaCmd = &cobra.Command{
	Use:   "delete-schema",
	Short: "Delete a schema",
	RunE: func(cmd *cobra.Command, args []string) error {
		cmd.PrintErrf("Deleting schema %s\n", schemaName)
		err := AppAPI.DeleteSchema(schemaName)
		if err != nil {
			return err
		}
		cmd.PrintErrln("Schema deleted successfully")
		return nil
	},
}

func printFields(w *tabwriter.Writer, fields []*model.Field, prefix string) {
	for _, f := range fields {
		name := prefix + f.Name
		nullable := ""
		if f.Nullable {
			nullable = "yes"
		}
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\n", name, f.Type, f.ConvertedType, nullable)
		if len(f.Fields) > 0 {
			printFields(w, f.Fields, name+".")
		}
	}
}

func init() {
	RootCmd.AddCommand(lakeCmd)
	lakeCmd.AddCommand(lakeAddCmd, lakeListCmd, lakeAddIndexCmd, lakeListIndexCmd, lakeDeleteIndexCmd, lakeAddSchemaCmd, lakeListSchemaCmd, lakeUpdateSchemaCmd, lakeDeleteSchemaCmd)
	//lakeAddCmd.AddCommand(lakeAddIndexCmd)

	lakeAddCmd.Flags().StringVar(&datalake, "datalake", "", "datalake name")
	lakeAddCmd.Flags().BoolVar(&managed, "managed", true, "Managed datalake (default: true)")
	lakeAddCmd.Flags().StringVar(&integration, "integration", "", "integration id")

	//_ = lakeAddCmd.MarkFlagRequired("datalake")
	// _ = lakeAddCmd.MarkFlagRequired("name")

	// Flags for 'lake add index'
	lakeAddIndexCmd.Flags().StringVar(&datalake, "datalake", "managed", "datalake name")
	lakeAddIndexCmd.Flags().StringVar(&index, "index", "", "datalake index name")
	lakeAddIndexCmd.Flags().StringVar(&schema, "schema", "ingext default", "schema name")

	//_ = lakeAddIndexCmd.MarkFlagRequired("datalake")
	_ = lakeAddIndexCmd.MarkFlagRequired("index")

	lakeListIndexCmd.Flags().StringVar(&datalake, "datalake", "", "datalake name")
	//_ = lakeListIndexCmd.MarkFlagRequired("datalake")

	lakeDeleteIndexCmd.Flags().StringVar(&datalake, "datalake", "", "datalake name")
	lakeDeleteIndexCmd.Flags().StringVar(&index, "index", "", "datalake index name")

	_ = lakeDeleteIndexCmd.MarkFlagRequired("datalake")
	_ = lakeDeleteIndexCmd.MarkFlagRequired("index")

	// Flags for 'datalake add-schema'
	lakeAddSchemaCmd.Flags().StringVar(&schemaName, "name", "", "schema name")
	lakeAddSchemaCmd.Flags().StringVar(&schemaDesc, "description", "", "schema description")
	lakeAddSchemaCmd.Flags().StringVar(&schemaFile, "schema", "", "path to schema file")

	_ = lakeAddSchemaCmd.MarkFlagRequired("name")
	_ = lakeAddSchemaCmd.MarkFlagRequired("schema")

	// Flags for 'datalake update-schema'
	lakeUpdateSchemaCmd.Flags().StringVar(&schemaName, "name", "", "schema name")
	lakeUpdateSchemaCmd.Flags().StringVar(&schemaDesc, "description", "", "schema description")
	lakeUpdateSchemaCmd.Flags().StringVar(&schemaFile, "schema", "", "path to schema file")

	_ = lakeUpdateSchemaCmd.MarkFlagRequired("name")
	_ = lakeUpdateSchemaCmd.MarkFlagRequired("schema")

	// Flags for 'datalake delete-schema'
	lakeDeleteSchemaCmd.Flags().StringVar(&schemaName, "name", "", "schema name")

	_ = lakeDeleteSchemaCmd.MarkFlagRequired("name")
}
