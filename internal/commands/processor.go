package commands

import (
	"fmt"
	"io"
	"os"

	"github.com/spf13/cobra"
)

var (
	procName    string
	procContent string
	procType    string // Default to "parser" if not specified
	procDesc    string // Optional description for the processor
)

var processorCmd = &cobra.Command{
	Use:   "processor",
	Short: "Manage processors",
}

var processorAddCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a processor",
	// Example usage:
	// 1. ingext processor add --name my-proc --content "@./my-script.js"
	// 2. ingext processor add --name my-proc --content "function process() { ... }"
	// 3. cat my-script.js | ingext processor add --name my-proc --content -
	RunE: func(cmd *cobra.Command, args []string) error {
		var content string
		//var err error

		// CHECK: Is the user asking to read from Stdin?
		if procContent == "-" {
			// Read from the pipe
			b, err := io.ReadAll(cmd.InOrStdin())
			if err != nil {
				return fmt.Errorf("failed to read from stdin: %w", err)
			}
			content = string(b)
		} else if len(procContent) > 1 && procContent[0] == '@' {
			filePath := procContent[1:]
			// Read from the file path provided
			b, err := os.ReadFile(filePath)
			if err != nil {
				return fmt.Errorf("failed to read file '%s': %w", filePath, err)
			}
			content = string(b)
		} else {
			content = procContent
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

// ingext processor add --name filter --content "@./scripts/filter.js"
// echo "function process() { ... }" | ingext processor add --name filter --content -
func init() {
	RootCmd.AddCommand(processorCmd)
	processorCmd.AddCommand(processorAddCmd) // Add del similarly

	//processorAddCmd.Flags().StringVar(&procName, "name", "", "Processor name")
	//processorAddCmd.Flags().StringVar(&procFile, "file", "", "Processor file path")

	processorAddCmd.Flags().StringVar(&procName, "name", "", "Processor name")
	processorAddCmd.Flags().StringVar(&procContent, "content", "", "Processor content or file path (use '-' for stdin)")
	processorAddCmd.Flags().StringVar(&procType, "type", "parser", "Processor type (parser|receiver|packer|report)")
	processorAddCmd.Flags().StringVar(&procDesc, "desc", "", "Processor description (optional)")

	_ = processorAddCmd.MarkFlagRequired("name")
	_ = processorAddCmd.MarkFlagRequired("content")
}
