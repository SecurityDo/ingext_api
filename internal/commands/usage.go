package commands

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"text/tabwriter"
	"time"

	fluencyAPI "github.com/SecurityDo/ingext_api/api"
	"github.com/SecurityDo/ingext_api/model"
	"github.com/spf13/cobra"
)

var (
	usageFrom        string
	usageTo          string
	usageTenantKey   string
	usageIncludeOpen bool
	usageJSON        bool
	usageDate        string
	usageAccounts    []string
	usageLimit       int
)

var usageCmd = &cobra.Command{
	Use:   "usage",
	Short: "Billing usage (metering) ledger",
	Long: "Read the daily billing usage ledger an account measures for itself.\n\n" +
		"Quantities and evidence only -- no prices, no SKUs, nothing from Stripe.\n" +
		"A meter that was not measured prints as '-', which is NOT zero: it means\n" +
		"the collection did not resolve, and totals over such a day are lower bounds.",
}

func usageOpts() *fluencyAPI.UsageOptions {
	return &fluencyAPI.UsageOptions{TenantKey: usageTenantKey, IncludeOpen: usageIncludeOpen}
}

// defaultRange is the last 30 whole UTC days, ending yesterday. Today is
// excluded because it has not finished and is never billable.
func defaultRange() (string, string) {
	end := time.Now().UTC().AddDate(0, 0, -1)
	return end.AddDate(0, 0, -29).Format("2006-01-02"), end.Format("2006-01-02")
}

func resolveRange() (string, string) {
	if usageFrom == "" || usageTo == "" {
		f, t := defaultRange()
		if usageFrom == "" {
			usageFrom = f
		}
		if usageTo == "" {
			usageTo = t
		}
	}
	return usageFrom, usageTo
}

func emitJSON(v interface{}) error {
	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return err
	}
	fmt.Println(string(b))
	return nil
}

// gb renders bytes as decimal GB, or "-" when the meter was not measured.
//
// The dash is deliberate and load-bearing: printing 0 for an unmeasured meter
// is how a reader concludes a tenant used nothing when in fact nothing was
// known.
func gb(d *model.UsageDay, meter string) string {
	v, ok := d.Bytes(meter)
	if !ok {
		return "-"
	}
	return fmt.Sprintf("%.3f", float64(v)/1e9)
}

var usageListCmd = &cobra.Command{
	Use:   "list",
	Short: "Show the daily usage ledger",
	RunE: func(cmd *cobra.Command, args []string) error {
		from, to := resolveRange()
		resp, err := AppAPI.DailyUsage(from, to, usageOpts())
		if err != nil {
			return err
		}
		if usageJSON {
			return emitJSON(resp)
		}
		if !resp.StoreAvailable {
			// Not a warning to bury: an unreachable ledger and a tenant with no
			// usage must never read the same.
			fmt.Fprintf(os.Stderr, "usage ledger UNAVAILABLE for %s: %s\n", resp.Account, resp.Reason)
			fmt.Fprintln(os.Stderr, "no conclusion can be drawn about this account's usage")
			return fmt.Errorf("usage ledger unavailable")
		}
		fmt.Printf("account %s   %s .. %s   (GB, decimal)\n", resp.Account, resp.From, resp.To)
		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "DATE\tSTATE\tSTATUS\tEVENTWATCH\tPROCESSED\tDELETED\tINPUT\tLAKE-IN\tUSERS")
		for _, d := range resp.Days {
			if d.State == model.DayStateMissing || d.State == model.DayStateInProgress {
				fmt.Fprintf(w, "%s\t%s\t\t\t\t\t\t\t\n", d.BillingDate, d.State)
				continue
			}
			var users []string
			for _, u := range d.PaidUsers {
				users = append(users, fmt.Sprintf("%s=%d", u.Provider, u.Quantity))
			}
			fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\n",
				d.BillingDate, d.State, d.Status,
				gb(d, model.MeterEventwatch), gb(d, model.MeterProcessed),
				gb(d, model.MeterDeleted), gb(d, model.MeterInput),
				gb(d, model.MeterLakeIngress), strings.Join(users, ","))
		}
		w.Flush()

		total, days, skipped := resp.MeterTotal(model.MeterEventwatch)
		fmt.Printf("\neventwatch over %d billable day(s): %.3f GB\n", days, float64(total)/1e9)
		if skipped > 0 {
			fmt.Printf("WARNING: %d billable day(s) had no eventwatch measurement; this total is a LOWER BOUND\n", skipped)
		}
		return nil
	},
}

var usageAttemptsCmd = &cobra.Command{
	Use:   "attempts",
	Short: "Show the collection attempt ledger",
	Long: "Every collection try, including the ones that failed and the ones that\n" +
		"declined to run. This is what separates 'this day has no data' from\n" +
		"'nobody ever looked'.",
	RunE: func(cmd *cobra.Command, args []string) error {
		from, to := resolveRange()
		resp, err := AppAPI.UsageAttempts(from, to, usageLimit, usageOpts())
		if err != nil {
			return err
		}
		if usageJSON {
			return emitJSON(resp)
		}
		if !resp.StoreAvailable {
			return fmt.Errorf("usage ledger unavailable")
		}
		if len(resp.Attempts) == 0 {
			cmd.Println("No attempts recorded in that range.")
			return nil
		}
		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "DATE\tFAMILY\tPROVIDER\tSTATUS\tMETHOD\tERROR\tACTOR")
		for _, a := range resp.Attempts {
			fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\t%s\t%s\n",
				a.BillingDate, a.MeterFamily, a.Provider, a.Status, a.Method, a.ErrorCode, a.Actor)
		}
		return w.Flush()
	},
}

var usageCollectCmd = &cobra.Command{
	Use:   "collect",
	Short: "Collect and close one past day on demand",
	Long: "Re-derive a day whose scheduled collection failed. It cannot rewrite a\n" +
		"day that already closed -- the ledger refuses that -- so the worst it can\n" +
		"do is fill a gap.\n\n" +
		"closed=false is a normal outcome: some meter did not resolve, so nothing\n" +
		"was written. Run 'usage attempts' to see which meter and why.",
	RunE: func(cmd *cobra.Command, args []string) error {
		if usageDate == "" {
			return fmt.Errorf("--date is required (YYYY-MM-DD, UTC, in the past)")
		}
		resp, err := AppAPI.CollectUsageDay(usageDate, usageOpts())
		if err != nil {
			return err
		}
		if usageJSON {
			return emitJSON(resp)
		}
		if resp.Closed {
			fmt.Printf("%s closed (actor %s)\n", resp.BillingDate, resp.Actor)
			return nil
		}
		fmt.Printf("%s NOT closed -- a meter did not resolve; see 'usage attempts --from %s --to %s'\n",
			resp.BillingDate, resp.BillingDate, resp.BillingDate)
		return nil
	},
}

var usageSinksCmd = &cobra.Command{
	Use:   "sinks",
	Short: "Show how datasinks are classified into meters",
	Long: "processed_bytes is the one meter whose value depends on a judgement about\n" +
		"the topology rather than a metric label, so this is the first thing to\n" +
		"check when that number looks wrong.",
	RunE: func(cmd *cobra.Command, args []string) error {
		resp, err := AppAPI.UsageSinkClassification()
		if err != nil {
			return err
		}
		if usageJSON {
			return emitJSON(resp)
		}
		fmt.Printf("account %s\n", resp.Account)
		fmt.Printf("  eventwatch (%d): %s\n", len(resp.Eventwatch), strings.Join(resp.Eventwatch, " "))
		fmt.Printf("  datalake   (%d): %s\n", len(resp.Datalake), strings.Join(resp.Datalake, " "))
		fmt.Printf("  processed  (%d): %s\n", len(resp.Processed), strings.Join(resp.Processed, " "))
		if len(resp.Ambiguous) > 0 {
			fmt.Printf("  ambiguous  (%d): %s\n", len(resp.Ambiguous), strings.Join(resp.Ambiguous, " "))
			fmt.Println("    these are counted as eventwatch; the platform's resident and job")
			fmt.Println("    modes disagree about them, and eventwatch under-counts rather than")
			fmt.Println("    double-bills")
		}
		if len(resp.Rejected) > 0 {
			fmt.Printf("  rejected   (%d): %s\n", len(resp.Rejected), strings.Join(resp.Rejected, " "))
			fmt.Println("    ids unsafe to place in a metric selector; excluded from processed")
		}
		return nil
	},
}

var usageGridCmd = &cobra.Command{
	Use:   "grid",
	Short: "Collect usage for every tenant of a provider",
	Long: "Issue against a PROVIDER site. --account narrows the set and can only\n" +
		"narrow it; anything out of scope is reported, never silently dropped.",
	RunE: func(cmd *cobra.Command, args []string) error {
		from, to := resolveRange()
		resp, err := AppAPI.GridUsage(from, to, usageAccounts, usageOpts())
		if err != nil {
			return err
		}
		if usageJSON {
			return emitJSON(resp)
		}
		fmt.Printf("grid %s   %s .. %s\n", resp.Grid, resp.From, resp.To)
		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "ACCOUNT\tDAYS\tEVENTWATCH-GB\tNOTE")
		for _, a := range resp.Accounts {
			if a.Error != "" {
				fmt.Fprintf(w, "%s\t-\t-\t%s\n", a.Account, a.Error)
				continue
			}
			var total int64
			var n, skipped int
			for _, d := range a.Days {
				if !d.Billable() {
					continue
				}
				if v, ok := d.Bytes(model.MeterEventwatch); ok {
					total += v
					n++
				} else {
					skipped++
				}
			}
			note := ""
			if skipped > 0 {
				note = fmt.Sprintf("%d day(s) unmeasured; lower bound", skipped)
			}
			fmt.Fprintf(w, "%s\t%d\t%.3f\t%s\n", a.Account, n, float64(total)/1e9, note)
		}
		w.Flush()

		fmt.Printf("\nrequested %d, succeeded %d, failed %d\n",
			resp.Summary.Requested, resp.Summary.Succeeded, resp.Summary.Failed)
		if !resp.Complete() {
			fmt.Fprintln(os.Stderr,
				"WARNING: some tenants could not be reached; this document is INCOMPLETE "+
					"and any provider total from it is a LOWER BOUND")
		}
		if len(resp.OutOfScope) > 0 {
			fmt.Fprintf(os.Stderr, "out of scope for this caller: %s\n", strings.Join(resp.OutOfScope, " "))
		}
		return nil
	},
}

func init() {
	RootCmd.AddCommand(usageCmd)
	usageCmd.AddCommand(usageListCmd, usageAttemptsCmd, usageCollectCmd, usageSinksCmd, usageGridCmd)

	for _, c := range []*cobra.Command{usageListCmd, usageAttemptsCmd, usageGridCmd} {
		c.Flags().StringVar(&usageFrom, "from", "", "Start date YYYY-MM-DD UTC (default: 30 days before --to)")
		c.Flags().StringVar(&usageTo, "to", "", "End date YYYY-MM-DD UTC, inclusive (default: yesterday)")
	}
	for _, c := range []*cobra.Command{usageListCmd, usageAttemptsCmd, usageCollectCmd, usageSinksCmd, usageGridCmd} {
		c.Flags().StringVar(&usageTenantKey, "tenant-key", "", "MSSP sub-customer (always empty today)")
		c.Flags().BoolVar(&usageJSON, "json", false, "Emit raw JSON")
	}
	for _, c := range []*cobra.Command{usageListCmd, usageGridCmd} {
		c.Flags().BoolVar(&usageIncludeOpen, "include-open", false,
			"Include today's provisional row (never billable)")
	}
	usageAttemptsCmd.Flags().IntVar(&usageLimit, "limit", 0, "Maximum attempts to return")
	usageCollectCmd.Flags().StringVar(&usageDate, "date", "", "Day to collect, YYYY-MM-DD UTC (required)")
	usageGridCmd.Flags().StringSliceVar(&usageAccounts, "account", nil,
		"Limit to these tenants (repeatable); can only narrow the caller's scope")
}
