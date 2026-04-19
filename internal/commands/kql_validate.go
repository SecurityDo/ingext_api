package commands

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var kqlValidateCmd = &cobra.Command{
	Use:   "validate <query or @file>",
	Short: "Parse a KQL query without executing it",
	Long: `Parse a KQL query and report whether it is syntactically valid.
The query is sent to the search service for validation; no rows are scanned
and no results are returned.

The argument can be either:
  - Inline KQL:  ingext kql validate "MyTable | where status == 200"
  - A file path: ingext kql validate @query.kql

Exit code 0 means the query parsed; any other exit code indicates the parser
rejected the input, with the error printed to stderr.`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		kql := args[0]
		if strings.HasPrefix(kql, "@") {
			data, err := os.ReadFile(kql[1:])
			if err != nil {
				return fmt.Errorf("read query file %s: %w", kql[1:], err)
			}
			kql = string(data)
		}
		kql = strings.TrimSpace(kql)
		if kql == "" {
			return fmt.Errorf("empty KQL query")
		}

		resp, err := AppAPI.KQLValidate(kql)
		if err != nil {
			return err
		}
		if !resp.OK {
			return fmt.Errorf("%s", resp.Error)
		}
		fmt.Fprintln(cmd.OutOrStdout(), "OK")
		return nil
	},
}

func init() {
	kqlCmd.AddCommand(kqlValidateCmd)
}
