package commands

import (
	"fmt"
	"io"
	"os"

	"github.com/spf13/cobra"
)

var (
	procName string
	procFile string
)

var processorCmd = &cobra.Command{
	Use:   "processor",
	Short: "Manage processors",
}

var processorAddCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a processor",
	//Run: func(cmd *cobra.Command, args []string) {
	//	fmt.Printf("Adding processor %s from file %s\n", procName, procFile)
	//},
	// Example usage:
	// 1. ingext processor add --name my-proc --file ./my-script.js
	// 2. cat my-script.js | ingext processor add --name my-proc --file -
	RunE: func(cmd *cobra.Command, args []string) error {
		var content []byte
		var err error

		// CHECK: Is the user asking to read from Stdin?
		if procFile == "-" {
			// Read from the pipe
			content, err = io.ReadAll(cmd.InOrStdin())
			if err != nil {
				return fmt.Errorf("failed to read from stdin: %w", err)
			}
		} else {
			// Read from the file path provided
			content, err = os.ReadFile(procFile)
			if err != nil {
				return fmt.Errorf("failed to read file '%s': %w", procFile, err)
			}
		}

		if len(content) == 0 {
			return fmt.Errorf("processor content is empty")
		}

		// Now you have the content in 'content' variable
		//  cmd.PrintErrln()
		//  cmd.Printf( )  for output data/result
		cmd.PrintErrf("Deploying processor '%s' (%d bytes)...\n", procName, len(content))

		// Call your global API
		// return AppAPI.AddProcessor(procName, content)
		return nil
	},
}

// ingext processor add --name filter --file ./scripts/filter.js
// echo "function process() { ... }" | ingext processor add --name filter --file -
func init() {
	RootCmd.AddCommand(processorCmd)
	processorCmd.AddCommand(processorAddCmd) // Add del similarly

	//processorAddCmd.Flags().StringVar(&procName, "name", "", "Processor name")
	//processorAddCmd.Flags().StringVar(&procFile, "file", "", "Processor file path")

	processorAddCmd.Flags().StringVar(&procName, "name", "", "Processor name")
	processorAddCmd.Flags().StringVar(&procFile, "file", "", "Processor file path (use '-' for stdin)")
	_ = processorAddCmd.MarkFlagRequired("file")
}
