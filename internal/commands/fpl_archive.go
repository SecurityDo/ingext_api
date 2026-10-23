package commands

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/SecurityDo/ingext_api/model"
	"github.com/spf13/cobra"
)

// fpl archive: the scheduled-report archive. A daily scheduled report with
// archive enabled keeps each day's result for 12 months; days are YYYY-MM-DD
// in the archive timezone and ranges include both ends.

var (
	archiveReq     model.ReportArchiveRequest
	archiveObjects string
	archiveOutput  string
	archiveJSON    bool
)

var fplArchiveCmd = &cobra.Command{
	Use:   "archive",
	Short: "Scheduled-report archive: daily runs kept for 12 months",
}

var fplArchiveListCmd = &cobra.Command{
	Use:   "list",
	Short: "List archived reports, or one report's days and coverage (list_fpl_report_archive)",
	RunE: func(cmd *cobra.Command, args []string) error {
		resp, err := AppAPI.ListReportArchive(&archiveReq)
		if err != nil {
			return err
		}
		if archiveJSON {
			return printArchiveJSON(cmd, resp)
		}
		w := tabwriter.NewWriter(cmd.OutOrStdout(), 0, 0, 2, ' ', 0)
		defer w.Flush()
		if archiveReq.Name == "" {
			fmt.Fprintln(w, "NAME\tENABLED\tTIMEZONE\tFIRST DAY\tLAST DAY\tDAYS\tCOMPLETED")
			for _, r := range resp.Reports {
				fmt.Fprintf(w, "%s\t%v\t%s\t%s\t%s\t%d\t%d\n", r.Name, r.Enabled, r.Timezone, r.FirstDay, r.LastDay, r.Days, r.Completed)
			}
			return nil
		}
		fmt.Fprintf(w, "%s (enabled %v, timezone %s)\n", resp.Name, resp.Enabled, resp.Timezone)
		if c := resp.Coverage; c != nil {
			fmt.Fprintf(w, "coverage %s..%s: %d expected, %d completed, %d aborted, %d running, %d missing\n",
				c.FirstDay, c.LastDay, c.Expected, c.Completed, c.Aborted, c.Running, c.Missing)
			if len(c.MissingDays) > 0 {
				fmt.Fprintf(w, "missing: %s\n", strings.Join(c.MissingDays, " "))
			}
			if len(c.AbortedDays) > 0 {
				fmt.Fprintf(w, "aborted: %s\n", strings.Join(c.AbortedDays, " "))
			}
			if len(c.ScriptHashes) > 1 {
				fmt.Fprintf(w, "script changed inside the range: %s\n", strings.Join(c.ScriptHashes, " "))
			}
		}
		if b := resp.Backfill; b != nil {
			fmt.Fprintf(w, "backfill: %d/%d done, %d failed, current %s\n", b.Done, len(b.Days), b.Failed, b.Current)
		}
		fmt.Fprintln(w, "ID\tDAY\tSTATE\tTASK\tROWS\tBYTES\tTOOK(ms)")
		for _, e := range resp.Entries {
			fmt.Fprintf(w, "%d\t%s\t%s\t%d\t%d\t%d\t%d\n", e.ID, e.RangeFrom.Format("2006-01-02"), e.State, e.TaskID, e.Total, e.ResultBytes, e.Took)
		}
		return nil
	},
}

var fplArchiveGetCmd = &cobra.Command{
	Use:   "get",
	Short: "Get one archived day by --id, or --name and --day (get_fpl_report_archive)",
	RunE: func(cmd *cobra.Command, args []string) error {
		if archiveReq.ID == 0 && (archiveReq.Name == "" || archiveReq.Day == "") {
			return fmt.Errorf("--id, or --name with --day, is required")
		}
		entry, err := AppAPI.GetReportArchive(&archiveReq)
		if err != nil {
			return err
		}
		return printArchiveJSON(cmd, entry)
	},
}

var fplArchiveDownloadCmd = &cobra.Command{
	Use:   "download",
	Short: "Download archived results: one day or a range of up to 92 days (download_fpl_report_archive)",
	RunE: func(cmd *cobra.Command, args []string) error {
		if archiveReq.ID == 0 && (archiveReq.Name == "" || (archiveReq.Day == "" && archiveReq.From == "")) {
			return fmt.Errorf("--id, or --name with --day or --from, is required")
		}
		if archiveObjects != "" {
			archiveReq.Objects = strings.Split(archiveObjects, ",")
		}
		resp, err := AppAPI.DownloadReportArchive(&archiveReq)
		if err != nil {
			return err
		}
		if c := resp.Coverage; c != nil && c.Completed < c.Expected {
			fmt.Fprintf(os.Stderr, "coverage: %d of %d days completed; missing %v aborted %v running %v\n",
				c.Completed, c.Expected, c.MissingDays, c.AbortedDays, c.RunningDays)
		}
		if archiveOutput == "" {
			return printArchiveJSON(cmd, resp)
		}
		b, err := json.MarshalIndent(resp, "", "  ")
		if err != nil {
			return fmt.Errorf("marshal response: %w", err)
		}
		if err := os.WriteFile(archiveOutput, b, 0644); err != nil {
			return fmt.Errorf("write output file %s: %w", archiveOutput, err)
		}
		fmt.Fprintf(os.Stderr, "%d day(s) saved to %s\n", len(resp.Days), archiveOutput)
		return nil
	},
}

var fplArchiveBackfillCmd = &cobra.Command{
	Use:   "backfill",
	Short: "Run the report for past missing days in the background (backfill_fpl_report_archive)",
	RunE: func(cmd *cobra.Command, args []string) error {
		if archiveReq.Name == "" || (archiveReq.Day == "" && archiveReq.From == "") {
			return fmt.Errorf("--name with --day or --from is required")
		}
		resp, err := AppAPI.BackfillReportArchive(&archiveReq)
		if err != nil {
			return err
		}
		cmd.Printf("queued %d day(s), skipped %d already completed (timezone %s)\n", len(resp.Queued), len(resp.Skipped), resp.Timezone)
		if len(resp.Queued) > 0 {
			cmd.Printf("progress: ingext fpl archive list --name %s\n", resp.Name)
		}
		return nil
	},
}

var fplArchiveReindexCmd = &cobra.Command{
	Use:   "reindex",
	Short: "Rebuild the archive index from the stored files, for --name or all reports (reindex_fpl_report_archive)",
	RunE: func(cmd *cobra.Command, args []string) error {
		resp, err := AppAPI.ReindexReportArchive(&archiveReq)
		if err != nil {
			return err
		}
		return printArchiveJSON(cmd, resp)
	},
}

func printArchiveJSON(cmd *cobra.Command, v interface{}) error {
	enc := json.NewEncoder(cmd.OutOrStdout())
	enc.SetIndent("", "  ")
	return enc.Encode(v)
}

func init() {
	fplCmd.AddCommand(fplArchiveCmd)
	fplArchiveCmd.AddCommand(fplArchiveListCmd, fplArchiveGetCmd, fplArchiveDownloadCmd, fplArchiveBackfillCmd, fplArchiveReindexCmd)

	for _, c := range []*cobra.Command{fplArchiveListCmd, fplArchiveGetCmd, fplArchiveDownloadCmd, fplArchiveBackfillCmd, fplArchiveReindexCmd} {
		c.Flags().StringVar(&archiveReq.Name, "name", "", "Scheduled report name")
	}
	for _, c := range []*cobra.Command{fplArchiveListCmd, fplArchiveDownloadCmd, fplArchiveBackfillCmd} {
		c.Flags().StringVar(&archiveReq.From, "from", "", "First day, YYYY-MM-DD in the archive timezone")
		c.Flags().StringVar(&archiveReq.To, "to", "", "Last day, YYYY-MM-DD, included")
	}
	for _, c := range []*cobra.Command{fplArchiveGetCmd, fplArchiveDownloadCmd, fplArchiveBackfillCmd} {
		c.Flags().StringVar(&archiveReq.Day, "day", "", "One day, YYYY-MM-DD in the archive timezone")
	}
	for _, c := range []*cobra.Command{fplArchiveGetCmd, fplArchiveDownloadCmd} {
		c.Flags().Int64Var(&archiveReq.ID, "id", 0, "Archive entry id (from list)")
	}

	fplArchiveListCmd.Flags().StringVar(&archiveReq.State, "state", "", "Only entries in this state: completed, aborted, running")
	fplArchiveListCmd.Flags().BoolVar(&archiveJSON, "json", false, "Print the full JSON response")

	fplArchiveDownloadCmd.Flags().StringVar(&archiveObjects, "objects", "", "Comma-separated result objects to keep (default all)")
	fplArchiveDownloadCmd.Flags().IntVar(&archiveReq.Limit, "limit", 0, "Rows per table to keep (default all)")
	fplArchiveDownloadCmd.Flags().StringVar(&archiveOutput, "output", "", "Save the JSON response to a file")

	fplArchiveBackfillCmd.Flags().BoolVar(&archiveReq.Force, "force", false, "Also rerun days that are already completed")
}
