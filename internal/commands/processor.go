package commands

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"

	ingextAPI "github.com/SecurityDo/ingext_api/api"
	"github.com/spf13/cobra"
)

var (
	procName    string
	procContent string
	procType    string // Default to "fpl_processor" if not specified
	procDesc    string // Optional description for the processor

	procTestScript    string
	procTestEvent     string
	procTestRawSource bool

	// The update flags keep their own vars: they default to "" so that an
	// unset --type or --desc keeps the stored one, and sharing procType with
	// 'add' would leave whichever command registered last deciding its default.
	procUpdateType string
	procUpdateDesc string
)

var processorCmd = &cobra.Command{
	Use:   "processor",
	Short: "Manage processors",
}

var listProcessorCmd = &cobra.Command{
	Use:   "list",
	Short: "List all processors",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("Listing processors...")

		entries, err := AppAPI.ListProcessor()
		if err != nil {
			return err
		}
		if len(entries) == 0 {
			cmd.PrintErrln("No processor found.")
			return nil
		}
		for _, entry := range entries {
			cmd.PrintErrf("Name: %s, Type: %s\n", entry.Name, entry.Type)
		}
		return nil
	},
}

var processorAddCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a processor",
	// Example usage:
	// 1. ingext processor add --name my-proc --content "@./my-script.js"
	// 2. ingext processor add --name my-proc --content "function process() { ... }"
	// 3. cat my-script.js | ingext processor add --name my-proc --content -
	RunE: func(cmd *cobra.Command, args []string) error {
		content, err := readProcessorContent(procContent, cmd.InOrStdin())
		if err != nil {
			return err
		}

		if len(content) == 0 {
			return fmt.Errorf("processor content is empty")
		}

		// Now you have the content in 'content' variable
		//  cmd.PrintErrln()
		//  cmd.Printf( )  for output data/result
		cmd.PrintErrf("Deploying processor '%s' (%d bytes)...\n", procName, len(content))
		err = AppAPI.AddProcessor(procName, content, procType, procDesc)
		if err != nil {
			return err
		}
		cmd.PrintErrln("Processor added successfully")
		return nil
	},
}

var processorUpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update the script of an existing processor",
	Long: `Replace the script of a processor that already exists.

--content takes the script inline, '@path' to read a file, or '-' to read stdin,
the same as 'processor add'. The stored entry is read back and patched, so the
processor keeps its id, group, tags and the repository it was imported from, and
an unset --type or --desc keeps the stored one.

  ingext processor update --name filter-logic --content @./scripts/filter.js

Updating a processor that does not exist is an error rather than an add.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		content, err := readProcessorContent(procContent, cmd.InOrStdin())
		if err != nil {
			return err
		}
		if len(content) == 0 {
			return fmt.Errorf("processor content is empty")
		}

		cmd.PrintErrf("Updating processor '%s' (%d bytes)...\n", procName, len(content))
		if err := AppAPI.UpdateProcessor(procName, content, procUpdateType, procUpdateDesc); err != nil {
			return err
		}
		cmd.PrintErrln("Processor updated successfully")
		return nil
	},
}

var processorDelCmd = &cobra.Command{
	Use:   "del",
	Short: "Delete a processor",
	RunE: func(cmd *cobra.Command, args []string) error {

		// Now you have the content in 'content' variable
		//  cmd.PrintErrln()
		//  cmd.Printf( )  for output data/result
		cmd.PrintErrf("Deleting processor '%s'...\n", procName)
		err := AppAPI.DeleteProcessor(procName)
		if err != nil {
			return err
		}
		cmd.PrintErrln("Processor deleted successfully")
		return nil
	},
}

var processorTestCmd = &cobra.Command{
	Use:   "test",
	Short: "Run a processor script against a sample event",
	Long: `Run an FPL script against one sample event without deploying it.

--script is the path to the script and --event the path to a JSON file holding
one event object. The object is wrapped in the document envelope the runtime
hands the script as main({obj, size}) -- {"obj": <event>, "props": {}, "size":
<compact byte length of the event>, "source": ""} -- so the file holds the event
itself, not the envelope.

  ingext processor test --script ./cloudtrail.fpl --event ./sample.json

The transformed envelope is printed to stdout, and the script's status, console
output and error to stderr. A script that aborts or drops the event is a result,
not a CLI failure; only a script error exits non-zero.

Either path may be '-' to read from stdin (only one of them). --raw-source sends
the --event file verbatim as the source instead of wrapping it, which is what
the non-processor types (fpl_receiver, fpl_packer, fpl_report) read.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if procTestScript == "-" && procTestEvent == "-" {
			return fmt.Errorf("only one of --script and --event can read from stdin")
		}
		script, err := readProcessorInput(procTestScript, "--script", cmd.InOrStdin())
		if err != nil {
			return err
		}
		if len(bytes.TrimSpace(script)) == 0 {
			return fmt.Errorf("script %s is empty", procTestScript)
		}
		event, err := readProcessorInput(procTestEvent, "--event", cmd.InOrStdin())
		if err != nil {
			return err
		}

		source := string(event)
		if procTestRawSource {
			if len(bytes.TrimSpace(event)) == 0 {
				return fmt.Errorf("%s is empty: the endpoint rejects an empty source", procTestEvent)
			}
		} else {
			doc, err := ingextAPI.NewFPLTestDocument(json.RawMessage(event))
			if err != nil {
				return fmt.Errorf("%s: %w", procTestEvent, err)
			}
			if source, err = doc.Encode(); err != nil {
				return err
			}
		}

		resp, err := AppAPI.TestProcessor(string(script), source, procType)
		if err != nil {
			return err
		}

		if resp.Console != "" {
			cmd.PrintErrln(strings.TrimRight(resp.Console, "\n"))
		}
		if resp.Status != "" {
			cmd.PrintErrf("status: %s\n", resp.Status)
		}
		if resp.NewContent != "" {
			var out bytes.Buffer
			if err := json.Indent(&out, []byte(resp.NewContent), "", "  "); err != nil {
				// Not JSON: a non-processor runtime returns its payload as text.
				cmd.Println(resp.NewContent)
			} else {
				cmd.Println(out.String())
			}
		}
		if resp.Error != "" {
			return fmt.Errorf("script error: %s", resp.Error)
		}
		return nil
	},
}

// readProcessorContent reads a --content argument: '-' reads stdin, '@path'
// reads that file, and anything else is the content itself.
func readProcessorContent(content string, stdin io.Reader) (string, error) {
	if content == "-" {
		b, err := io.ReadAll(stdin)
		if err != nil {
			return "", fmt.Errorf("failed to read from stdin: %w", err)
		}
		return string(b), nil
	}
	if len(content) > 1 && content[0] == '@' {
		filePath := content[1:]
		b, err := os.ReadFile(filePath)
		if err != nil {
			return "", fmt.Errorf("failed to read file '%s': %w", filePath, err)
		}
		return string(b), nil
	}
	return content, nil
}

// readProcessorInput reads a command input that is normally a file path: '-'
// reads stdin, a leading '@' is accepted for symmetry with 'processor add'.
func readProcessorInput(path, flag string, stdin io.Reader) ([]byte, error) {
	if path == "" {
		return nil, fmt.Errorf("%s is required", flag)
	}
	if path == "-" {
		b, err := io.ReadAll(stdin)
		if err != nil {
			return nil, fmt.Errorf("failed to read %s from stdin: %w", flag, err)
		}
		return b, nil
	}
	b, err := os.ReadFile(strings.TrimPrefix(path, "@"))
	if err != nil {
		return nil, fmt.Errorf("failed to read %s file: %w", flag, err)
	}
	return b, nil
}

// ingext processor add --name filter --content "@./scripts/filter.js"
// echo "function process() { ... }" | ingext processor add --name filter --content -
func init() {
	RootCmd.AddCommand(processorCmd)
	processorCmd.AddCommand(processorAddCmd, processorUpdateCmd, listProcessorCmd, processorDelCmd, processorTestCmd) // Add del similarly

	//processorAddCmd.Flags().StringVar(&procName, "name", "", "Processor name")
	//processorAddCmd.Flags().StringVar(&procFile, "file", "", "Processor file path")

	processorAddCmd.Flags().StringVar(&procName, "name", "", "Processor name")
	processorAddCmd.Flags().StringVar(&procContent, "content", "", "Processor content or file path (use '-' for stdin)")
	processorAddCmd.Flags().StringVar(&procType, "type", "fpl_processor", "Processor type (fpl_processor|fpl_receiver|fpl_packer|fpl_report)")
	processorAddCmd.Flags().StringVar(&procDesc, "desc", "", "Processor description (optional)")

	_ = processorAddCmd.MarkFlagRequired("name")
	_ = processorAddCmd.MarkFlagRequired("content")

	processorUpdateCmd.Flags().StringVar(&procName, "name", "", "Processor name")
	processorUpdateCmd.Flags().StringVar(&procContent, "content", "", "Processor content or file path (use '-' for stdin)")
	processorUpdateCmd.Flags().StringVar(&procUpdateType, "type", "", "Processor type (fpl_processor|fpl_receiver|fpl_packer|fpl_report); unset keeps the stored type")
	processorUpdateCmd.Flags().StringVar(&procUpdateDesc, "desc", "", "Processor description; unset keeps the stored description")

	_ = processorUpdateCmd.MarkFlagRequired("name")
	_ = processorUpdateCmd.MarkFlagRequired("content")

	processorDelCmd.Flags().StringVar(&procName, "name", "", "Processor name")
	_ = processorDelCmd.MarkFlagRequired("name")

	processorTestCmd.Flags().StringVar(&procTestScript, "script", "", "Path to the FPL script to test (use '-' for stdin)")
	processorTestCmd.Flags().StringVar(&procTestEvent, "event", "", "Path to a JSON file holding one event object (use '-' for stdin)")
	processorTestCmd.Flags().BoolVar(&procTestRawSource, "raw-source", false, "Send the --event file verbatim as the source instead of wrapping it in the {obj, size} envelope")
	processorTestCmd.Flags().StringVar(&procType, "type", "fpl_processor", "Processor type (fpl_processor|fpl_receiver|fpl_packer|fpl_report)")

	_ = processorTestCmd.MarkFlagRequired("script")
	_ = processorTestCmd.MarkFlagRequired("event")

}
