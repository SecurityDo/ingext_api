package commands

import (
	"fmt"
	"github.com/spf13/cobra"
)

var (
	lakeStorage   string
	lakeBucket    string
	lakePrefix    string
	lakeAccount   string
	lakeContainer string
)

var lakeCmd = &cobra.Command{
	Use:   "lake",
	Short: "Manage data lake",
}

var lakeAddCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a lake resource",
}

var lakeAddIndexCmd = &cobra.Command{
	Use:   "index",
	Short: "Add an index to the lake",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Adding lake index to storage %s (Bucket: %s)\n", lakeStorage, lakeBucket)
	},
}

func init() {
	RootCmd.AddCommand(lakeCmd)
	lakeCmd.AddCommand(lakeAddCmd)
	lakeAddCmd.AddCommand(lakeAddIndexCmd)

	// Flags for 'lake add index'
	lakeAddIndexCmd.Flags().StringVar(&lakeStorage, "storage", "", "Storage type (s3|blob|gcs)")
	lakeAddIndexCmd.Flags().StringVar(&lakeBucket, "bucket", "", "Bucket name")
	lakeAddIndexCmd.Flags().StringVar(&lakePrefix, "prefix", "", "Path prefix")
	lakeAddIndexCmd.Flags().StringVar(&lakeAccount, "storageaccount", "", "Storage account")
	lakeAddIndexCmd.Flags().StringVar(&lakeContainer, "container", "", "Container name")
}
