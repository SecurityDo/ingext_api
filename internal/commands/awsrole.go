package commands

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	roleName        string
	roleDisplayName string
	roleExternalID  string
	roleARN         string
	roleID          string
)

var eksCmd = &cobra.Command{
	Use:   "eks",
	Short: "Manage EKS pod identity assumed roles",
}

var getPodRoleCmd = &cobra.Command{
	Use:   "get-pod-role",
	Short: "Get iam role for Pod Identity Agent",
	RunE: func(cmd *cobra.Command, args []string) error {
		cmd.PrintErrln("Getting Pod Role...")
		// Just call the global interface
		role, arn, err := AppAPI.GetPodRole()
		if err != nil {
			return err
		}

		cmd.PrintErrln("Get Pod Role: ", role, " ARN: ", arn)
		fmt.Println(role)
		return nil
	},
}

var testAssumedRoleCmd = &cobra.Command{
	Use:   "test-assumed-role",
	Short: "Test Assumed Roles for Pod Identity Agent",
	RunE: func(cmd *cobra.Command, args []string) error {
		//cmd.PrintErrln("Testing AWS Role...")
		// Just call the global interface
		err := AppAPI.TestAssumedRole(roleARN, roleExternalID)
		if err != nil {
			fmt.Println("Error testing assumed role:", err)
			return nil
		}

		cmd.PrintErrln("Role tested successfully: ")
		fmt.Println("OK")
		return nil
	},
}

var addAssumedRoleCmd = &cobra.Command{
	Use:   "add-assumed-role",
	Short: "Add Assumed Roles for Pod Identity Agent",
	RunE: func(cmd *cobra.Command, args []string) error {
		cmd.PrintErrln("Adding AWS Role...")
		// Just call the global interface
		id, err := AppAPI.AddAssumedRole(roleName, roleARN, roleExternalID)
		if err != nil {
			return err
		}

		cmd.PrintErrln("Role added successfully: ", id)
		fmt.Println(id)
		return nil
	},
}

var delAssumedRoleCmd = &cobra.Command{
	Use:   "del-assumed-role",
	Short: "Delete Assumed Roles for Pod Identity Agent",
	RunE: func(cmd *cobra.Command, args []string) error {
		cmd.PrintErrln("Deleting AWS Role...")
		// Just call the global interface
		err := AppAPI.DeleteAssumedRole(roleID)
		if err != nil {
			return err
		}

		cmd.PrintErrln("Role deleted successfully: ")
		//cmd.Println(id)
		return nil
	},
}

var listAssumedRoleCmd = &cobra.Command{
	Use:   "list-assumed-role",
	Short: "List Assumed Roles for Pod Identity Agent",
	RunE: func(cmd *cobra.Command, args []string) error {
		//cmd.PrintErrln("Deleting AWS Role...")
		// Just call the global interface
		roles, err := AppAPI.ListAssumedRole()
		if err != nil {
			return err
		}
		cmd.PrintErrln("Listing AWS Roles...")
		if len(roles) == 0 {
			cmd.Println("No roles found.")
			return nil
		}

		for _, role := range roles {
			cmd.PrintErrf("Role ID: %s, Name: %s, ARN: %s, External ID: %s\n", role.ID, role.DisplayName, role.RoleARN, role.ExternalID)
		}
		return nil
	},
}

func init() {
	RootCmd.AddCommand(eksCmd)
	eksCmd.AddCommand(addAssumedRoleCmd, delAssumedRoleCmd, listAssumedRoleCmd, getPodRoleCmd, testAssumedRoleCmd) // Add del/update similarly

	addAssumedRoleCmd.Flags().StringVar(&roleName, "name", "", "displayName of the role")
	//addAssumedRoleCmd.Flags().StringVar(&roleDisplayName, "displayName", "", "Display name")
	addAssumedRoleCmd.Flags().StringVar(&roleARN, "roleArn", "", "Role ARN to assume")
	addAssumedRoleCmd.Flags().StringVar(&roleExternalID, "externalId", "", "External ID (optional)")

	// Mark required
	_ = addAssumedRoleCmd.MarkFlagRequired("name")
	_ = addAssumedRoleCmd.MarkFlagRequired("roleArn")

	//streamAddCmd.AddCommand(streamAddSourceCmd)
	// Add other leaf commands: sink, router, connection

	testAssumedRoleCmd.Flags().StringVar(&roleARN, "roleArn", "", "Role ARN to assume")
	testAssumedRoleCmd.Flags().StringVar(&roleExternalID, "externalId", "", "External ID (optional)")

	// Mark required
	_ = testAssumedRoleCmd.MarkFlagRequired("roleArn")

}
