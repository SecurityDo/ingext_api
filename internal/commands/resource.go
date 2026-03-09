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

func init() {
	resourceCmd.Flags().StringVar(&resourceType, "resource-type", "", "specify the resource type")
	resourceCmd.Flags().StringVar(&customer, "customer", "_all_", "specify the customer")

	RootCmd.AddCommand(resourceCmd)
}
