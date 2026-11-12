package commands

import (
	"fmt"
	"sort"
	"strings"

	ingextAPI "github.com/SecurityDo/ingext_api/api"
	model "github.com/SecurityDo/ingext_api/model"
	"github.com/spf13/cobra"
)

// Removing a data source is three platform calls, in this order:
//
//  1. platform_source_disable     stop it (persists disable=true)
//  2. platform_source_set_router  routerID "" -- drop the routerInput connection
//  3. platform_datasource_dao     action delete
//
// The delete refuses a source that is still "referenced by router <id>", so
// step 2 is not optional. A source an application installed (tag
// application=<App>) is not removed here: uninstall the app instance instead,
// or the app keeps a dangling reference to a source that no longer exists.

var (
	removeSources  []string
	idleDays       int
	allowActive    bool
	removeDryRun   bool
	usageTypeMatch string
)

type sourceView struct {
	src      *model.DataSourceConfig
	routerID string
	router   string
	app      string
	input    float64
	measured bool
}

func appOwner(src *model.DataSourceConfig) string {
	for _, t := range src.Tags {
		if t != nil && t.Name == "application" {
			return t.Value
		}
	}
	return ""
}

// loadSourceViews reads the topology once and resolves each source's router
// and owning application.
func loadSourceViews(configs *ingextAPI.ListConfigsResponse) map[string]*sourceView {
	routerNames := map[string]string{}
	for _, r := range configs.Routers {
		if r != nil {
			routerNames[r.ID] = r.Name
		}
	}
	views := map[string]*sourceView{}
	for _, s := range configs.Sources {
		if s != nil {
			views[s.ID] = &sourceView{src: s, app: appOwner(s)}
		}
	}
	for _, c := range configs.Connections {
		if c == nil {
			continue
		}
		if v, ok := views[c.SourceID]; ok {
			v.routerID = c.RouterID
			v.router = routerNames[c.RouterID]
		}
	}
	return views
}

func findSourceView(views map[string]*sourceView, ref string) (*sourceView, error) {
	if v, ok := views[ref]; ok {
		return v, nil
	}
	var hits []*sourceView
	for _, v := range views {
		if v.src.Name == ref {
			hits = append(hits, v)
		}
	}
	switch len(hits) {
	case 1:
		return hits[0], nil
	case 0:
		return nil, fmt.Errorf("no data source named or numbered %q", ref)
	default:
		return nil, fmt.Errorf("%d data sources are named %q; pass the src_ id", len(hits), ref)
	}
}

func fillInput(v *sourceView, days int) error {
	n, measured, err := AppAPI.SourceInputCount(v.src.ID, days)
	if err != nil {
		return err
	}
	v.input, v.measured = n, measured
	return nil
}

func inputLabel(v *sourceView) string {
	if !v.measured {
		return "unmeasured"
	}
	return fmt.Sprintf("%.0f", v.input)
}

var sourceUsageCmd = &cobra.Command{
	Use:   "source-usage",
	Short: "List data sources with router, owner, state and input over the last N days",
	Long: `List every data source of the account with the router it feeds, the
application that owns it (if any), whether it is disabled, and the number of
events it took in over the last --days days (platform_component_metrics).

"unmeasured" means the metrics store has no input series for the source at all;
that is not the same as zero.

  ingext stream source-usage --days 30 --type hec`,
	RunE: func(cmd *cobra.Command, args []string) error {
		configs, err := AppAPI.ListConfigs()
		if err != nil {
			return err
		}
		views := loadSourceViews(configs)
		list := make([]*sourceView, 0, len(views))
		for _, v := range views {
			if usageTypeMatch != "" && v.src.Type != usageTypeMatch {
				continue
			}
			if err := fillInput(v, idleDays); err != nil {
				return err
			}
			list = append(list, v)
		}
		sort.Slice(list, func(i, j int) bool { return list[i].src.Name < list[j].src.Name })
		fmt.Printf("%-36s %-15s %-9s %-8s %-30s %-24s %s\n", "NAME", "ID", "TYPE", "STATE", "ROUTER", "APP", fmt.Sprintf("INPUT_%dD", idleDays))
		for _, v := range list {
			state := "enabled"
			if v.src.Disable {
				state = "disabled"
			}
			router, app := v.router, v.app
			if router == "" {
				router = "-"
			}
			if app == "" {
				app = "-"
			}
			fmt.Printf("%-36s %-15s %-9s %-8s %-30s %-24s %s\n", v.src.Name, v.src.ID, v.src.Type, state, router, app, inputLabel(v))
		}
		return nil
	},
}

var removeSourceCmd = &cobra.Command{
	Use:   "remove-source",
	Short: "Disable, disconnect and delete idle data sources",
	Long: `Remove standalone data sources the way the platform requires:
disable -> disconnect from its router -> delete, then read the topology back.

Every source is checked first and nothing is changed unless ALL of them pass:
  - it exists (name or src_ id; a duplicated name needs the id)
  - no application owns it (tag application=<App>) -- uninstall the app instead
  - it took in 0 events over the last --idle-days days (default 30); a source
    with no input series at all ("unmeasured") also fails unless --allow-active

  ingext stream remove-source --source ExtraHop --source OracleDB --dry-run
  ingext stream remove-source --source ExtraHop --source OracleDB

Routers and sinks are left alone. A router left with no sources is reported,
not deleted: other sources or a later re-connect may still need it.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(removeSources) == 0 {
			return fmt.Errorf("--source is required (repeatable, name or id)")
		}
		configs, err := AppAPI.ListConfigs()
		if err != nil {
			return err
		}
		views := loadSourceViews(configs)

		// Pre-flight every source before touching any of them.
		var targets []*sourceView
		seen := map[string]bool{}
		var problems []string
		for _, ref := range removeSources {
			v, err := findSourceView(views, ref)
			if err != nil {
				problems = append(problems, err.Error())
				continue
			}
			if seen[v.src.ID] {
				continue
			}
			seen[v.src.ID] = true
			if v.app != "" {
				problems = append(problems, fmt.Sprintf("%s (%s) belongs to application %s: uninstall that app instance instead", v.src.Name, v.src.ID, v.app))
				continue
			}
			if err := fillInput(v, idleDays); err != nil {
				return err
			}
			if !allowActive && (!v.measured || v.input > 0) {
				problems = append(problems, fmt.Sprintf("%s (%s) is not idle: %s events in %dd (--allow-active to override)", v.src.Name, v.src.ID, inputLabel(v), idleDays))
				continue
			}
			targets = append(targets, v)
		}
		for _, v := range targets {
			router := v.router
			if v.routerID == "" {
				router = "(none)"
			}
			fmt.Printf("plan: %-30s %-15s type=%-7s router=%s input_%dd=%s disabled=%v\n", v.src.Name, v.src.ID, v.src.Type, router, idleDays, inputLabel(v), v.src.Disable)
		}
		if len(problems) > 0 {
			return fmt.Errorf("nothing changed:\n  %s", strings.Join(problems, "\n  "))
		}
		if removeDryRun {
			fmt.Println("dry run: nothing changed")
			return nil
		}

		for _, v := range targets {
			if err := AppAPI.DisableDataSource(v.src.ID); err != nil {
				return fmt.Errorf("%s: %w", v.src.Name, err)
			}
			if v.routerID != "" {
				if err := AppAPI.SetSourceRouter(v.src.ID, ""); err != nil {
					return fmt.Errorf("%s: disabled, but disconnect from router %s failed: %w", v.src.Name, v.router, err)
				}
			}
			if err := AppAPI.DeleteDataSource(v.src.ID); err != nil {
				return fmt.Errorf("%s: disabled and disconnected, but delete failed: %w", v.src.Name, err)
			}
			fmt.Printf("removed: %s (%s)\n", v.src.Name, v.src.ID)
		}

		// Verify against a fresh read.
		after, err := AppAPI.ListConfigs()
		if err != nil {
			return fmt.Errorf("removed, but the verification read failed: %w", err)
		}
		left := loadSourceViews(after)
		var still []string
		for _, v := range targets {
			if _, ok := left[v.src.ID]; ok {
				still = append(still, v.src.ID)
			}
		}
		if len(still) > 0 {
			return fmt.Errorf("still listed after delete: %s", strings.Join(still, ", "))
		}
		fed := map[string]int{}
		for _, c := range after.Connections {
			if c != nil {
				fed[c.RouterID]++
			}
		}
		empty := map[string]string{}
		for _, v := range targets {
			if v.routerID != "" && fed[v.routerID] == 0 {
				empty[v.routerID] = v.router
			}
		}
		for id, name := range empty {
			fmt.Printf("note: router %s (%s) now has no sources\n", name, id)
		}
		fmt.Printf("verified: %d source(s) gone\n", len(targets))
		return nil
	},
}

func init() {
	streamCmd.AddCommand(sourceUsageCmd, removeSourceCmd)

	sourceUsageCmd.Flags().IntVar(&idleDays, "days", 30, "input window in days")
	sourceUsageCmd.Flags().StringVar(&usageTypeMatch, "type", "", "only sources of this type (hec, plugin, s3, ...)")

	removeSourceCmd.Flags().StringArrayVar(&removeSources, "source", nil, "data source name or src_ id (repeatable)")
	removeSourceCmd.Flags().IntVar(&idleDays, "idle-days", 30, "a source must have taken in 0 events over this many days")
	removeSourceCmd.Flags().BoolVar(&allowActive, "allow-active", false, "remove even if the source took in events (or is unmeasured)")
	removeSourceCmd.Flags().BoolVar(&removeDryRun, "dry-run", false, "check and print the plan; change nothing")
}
