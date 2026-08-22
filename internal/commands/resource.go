package commands

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
)

var resourceType string
var customer string

func PrettyPrintJSON(x interface{}) {
	pretty, _ := json.MarshalIndent(x, "", "   ")
	fmt.Printf("%s\n", pretty)
}

var resourceCmd = &cobra.Command{
	Use:   "resource ",
	Short: "Resource search",
	Long:  `Run a resource search`,
	//Args:  cobra.ExactArgs(0),

	RunE: func(cmd *cobra.Command, args []string) error {

		resp, err := AppAPI.ResourceSearch(resourceType, customer)
		if err != nil {
			return err
		}
		PrettyPrintJSON(resp)

		return nil
	},
}

var dumpCustomer string

var resourceDumpDeleteCmd = &cobra.Command{
	Use:   "dump_delete",
	Short: "Delete every resource dump of one customer",
	Long: `Purge the resource dumps of one customer across all resource types.

The customer is matched literally, so "_all_" is a customer name here, not a
wildcard. Deleting collected data this way has no undo, and an unknown customer
is not an error - it simply deletes nothing.`,
	RunE: func(cmd *cobra.Command, args []string) error {

		deleted, err := AppAPI.DeleteResourceDump(dumpCustomer)
		if err != nil {
			return err
		}
		if deleted == 0 {
			cmd.PrintErrf("No resource dump found for customer '%s'\n", dumpCustomer)
			return nil
		}
		cmd.PrintErrf("Deleted %d resource dump(s) for customer '%s'\n", deleted, dumpCustomer)
		return nil
	},
}

func init() {
	resourceCmd.Flags().StringVar(&resourceType, "resource-type", "", "specify the resource type")
	resourceCmd.Flags().StringVar(&customer, "customer", "_all_", "specify the customer")

	resourceDumpDeleteCmd.Flags().StringVar(&dumpCustomer, "customer", "", "customer whose resource dumps are deleted")
	_ = resourceDumpDeleteCmd.MarkFlagRequired("customer")
	resourceCmd.AddCommand(resourceDumpDeleteCmd)

	RootCmd.AddCommand(resourceCmd)
}
