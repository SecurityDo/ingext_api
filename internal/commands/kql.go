package commands

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	kqlModel "github.com/SecurityDo/ingext_api/kql/model"
	"github.com/spf13/cobra"
)

var kqlOutput string

var kqlCmd = &cobra.Command{
	Use:   "kql <query or @file>",
	Short: "Run a KQL query against the datalake",
	Long: `Run a KQL query and display the results as a table.

The argument can be either:
  - Inline KQL:  ingext kql "MyTable | where status == 200 | take 10"
  - A file path: ingext kql @query.kql`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		kql := args[0]

		// If argument starts with @, read KQL from file
		if strings.HasPrefix(kql, "@") {
			filePath := kql[1:]
			data, err := os.ReadFile(filePath)
			if err != nil {
				return fmt.Errorf("read query file %s: %w", filePath, err)
			}
			kql = string(data)
		}

		kql = strings.TrimSpace(kql)
		if kql == "" {
			return fmt.Errorf("empty KQL query")
		}

		resp, err := AppAPI.KQLSearch(kql)
		if err != nil {
			return err
		}

		// Save to file if --output is specified
		if kqlOutput != "" {
			jsonBytes, err := json.MarshalIndent(resp, "", "  ")
			if err != nil {
				return fmt.Errorf("marshal response: %w", err)
			}
			if err := os.WriteFile(kqlOutput, jsonBytes, 0644); err != nil {
				return fmt.Errorf("write output file %s: %w", kqlOutput, err)
			}
			fmt.Fprintf(os.Stderr, "Response saved to %s\n", kqlOutput)
		}

		if resp.Data == nil || len(resp.Data.Tables) == 0 {
			fmt.Printf("(%d rows, %d bytes scanned)\n", resp.Total, resp.TotalBytes)
			return nil
		}

		for i, table := range resp.Data.Tables {
			if i > 0 {
				fmt.Println()
			}
			printKQLTable(table)
		}

		fmt.Printf("\n(%d rows, %d bytes scanned)\n", resp.Total, resp.TotalBytes)
		return nil
	},
}

func init() {
	kqlCmd.Flags().StringVar(&kqlOutput, "output", "", "save the full JSON response to a file")
	RootCmd.AddCommand(kqlCmd)
}

func printKQLTable(table *kqlModel.DataTable) {
	if len(table.Columns) == 0 {
		fmt.Printf("Table: %s (empty)\n", table.Name)
		return
	}

	// Calculate column widths
	colWidths := make([]int, len(table.Columns))
	for i, col := range table.Columns {
		colWidths[i] = len(col.Name)
	}
	for _, row := range table.Rows {
		for i, val := range row.Values {
			s := formatKQLValue(val)
			if len(s) > colWidths[i] {
				colWidths[i] = len(s)
			}
		}
	}

	const maxColWidth = 60
	for i := range colWidths {
		if colWidths[i] > maxColWidth {
			colWidths[i] = maxColWidth
		}
	}

	// Build format strings
	formats := make([]string, len(colWidths))
	for i, w := range colWidths {
		formats[i] = fmt.Sprintf("%%-%ds", w)
	}

	printKQLSeparator(colWidths)

	// Header
	fmt.Print("|")
	for i, col := range table.Columns {
		fmt.Printf(" "+formats[i]+" |", truncateStr(col.Name, colWidths[i]))
	}
	fmt.Println()

	printKQLSeparator(colWidths)

	// Rows
	for _, row := range table.Rows {
		fmt.Print("|")
		for i, val := range row.Values {
			s := formatKQLValue(val)
			fmt.Printf(" "+formats[i]+" |", truncateStr(s, colWidths[i]))
		}
		fmt.Println()
	}

	printKQLSeparator(colWidths)
	fmt.Printf("(%d rows)\n", len(table.Rows))
}

func printKQLSeparator(colWidths []int) {
	fmt.Print("+")
	for _, w := range colWidths {
		fmt.Print(strings.Repeat("-", w+2) + "+")
	}
	fmt.Println()
}

func formatKQLValue(val kqlModel.KValue) string {
	if val == nil || val.IsNull() {
		return "<null>"
	}
	return val.String()
}

func truncateStr(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	if maxLen <= 3 {
		return s[:maxLen]
	}
	return s[:maxLen-3] + "..."
}
