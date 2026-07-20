package commands

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/SecurityDo/ingext_api/model"
	"github.com/spf13/cobra"
)

var (
	eventwatchQuery   string
	eventwatchFrom    int64
	eventwatchTo      int64
	eventwatchName    string
	eventwatchGroup   string
	eventwatchContent string
)

var eventwatchCmd = &cobra.Command{
	Use:   "eventwatch",
	Short: "EventWatch service",
}

var eventwatchSummarySearchCmd = &cobra.Command{
	Use:   "search_summary",
	Short: "Run summary search",
	RunE: func(cmd *cobra.Command, args []string) error {
		from, to := eventwatchFrom, eventwatchTo
		if from == 0 && to == 0 {
			now := time.Now().UnixMilli()
			to = now
			from = now - int64(time.Hour/time.Millisecond)
		}
		resp, err := AppAPI.SummarySearch(eventwatchQuery, from, to)
		if err != nil {
			return err
		}
		return printEventwatchHits(cmd, resp, "BehaviorSummary")
	},
}

var eventwatchTimelineSearchCmd = &cobra.Command{
	Use:   "search_timeline",
	Short: "Run timeline search (fsm_behavior_search)",
	RunE: func(cmd *cobra.Command, args []string) error {
		from, to := eventwatchFrom, eventwatchTo
		if from == 0 && to == 0 {
			now := time.Now().UnixMilli()
			to = now
			from = now - int64(time.Hour/time.Millisecond)
		}
		resp, err := AppAPI.TimelineSearch(eventwatchQuery, from, to)
		if err != nil {
			return err
		}
		return printEventwatchHits(cmd, resp, "BehaviorEvent")
	},
}

var eventwatchRuleSearchCmd = &cobra.Command{
	Use:   "search_rule",
	Short: "Run rule search (eventwatch_bucket_search)",
	RunE: func(cmd *cobra.Command, args []string) error {
		resp, err := AppAPI.RuleSearch(eventwatchQuery)
		if err != nil {
			return err
		}
		return printEventwatchHits(cmd, resp, "BehaviorRule")
	},
}

// --- EventWatch rule DAO commands (eventwatch_bucket_dao) ---

var eventwatchRuleListCmd = &cobra.Command{
	Use:   "rule_list",
	Short: "List eventwatch rules",
	RunE: func(cmd *cobra.Command, args []string) error {
		resp, err := AppAPI.ListRule()
		if err != nil {
			return err
		}
		if len(resp.Entries) == 0 {
			cmd.PrintErrln("No rules found.")
			return nil
		}
		for _, entry := range resp.Entries {
			state := "enabled"
			if entry.Disabled {
				state = "disabled"
			}
			cmd.Printf("Name: %s, Group: %s, State: %s\n", entry.Name, entry.Group, state)
		}
		return nil
	},
}

var eventwatchRuleGetCmd = &cobra.Command{
	Use:   "rule_get",
	Short: "Get a single eventwatch rule by name",
	RunE: func(cmd *cobra.Command, args []string) error {
		entry, err := AppAPI.GetRule(eventwatchName)
		if err != nil {
			return err
		}
		if entry == nil {
			cmd.PrintErrln("Rule not found.")
			return nil
		}
		b, err := json.MarshalIndent(entry, "", "  ")
		if err != nil {
			return err
		}
		cmd.Println(string(b))
		return nil
	},
}

var eventwatchRuleAddCmd = &cobra.Command{
	Use:   "rule_add",
	Short: "Add an eventwatch rule from a JSON definition",
	// Example usage:
	// 1. ingext eventwatch rule_add --content "@./rule.json"
	// 2. cat rule.json | ingext eventwatch rule_add --content -
	RunE: func(cmd *cobra.Command, args []string) error {
		entry, err := loadEventwatchRule(cmd)
		if err != nil {
			return err
		}
		if err := AppAPI.AddRule(entry); err != nil {
			return err
		}
		cmd.PrintErrf("Rule '%s' added successfully\n", entry.Name)
		return nil
	},
}

var eventwatchRuleUpdateCmd = &cobra.Command{
	Use:   "rule_update",
	Short: "Update an eventwatch rule from a JSON definition",
	// Example usage:
	// 1. ingext eventwatch rule_update --content "@./rule.json"
	// 2. cat rule.json | ingext eventwatch rule_update --content -
	RunE: func(cmd *cobra.Command, args []string) error {
		entry, err := loadEventwatchRule(cmd)
		if err != nil {
			return err
		}
		if err := AppAPI.UpdateRule(entry); err != nil {
			return err
		}
		cmd.PrintErrf("Rule '%s' updated successfully\n", entry.Name)
		return nil
	},
}

var eventwatchRuleToggleCmd = &cobra.Command{
	Use:   "rule_toggle",
	Short: "Toggle the disabled state of an eventwatch rule",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := AppAPI.ToggleRule(eventwatchName); err != nil {
			return err
		}
		cmd.PrintErrf("Rule '%s' toggled successfully\n", eventwatchName)
		return nil
	},
}

var eventwatchRuleDeleteCmd = &cobra.Command{
	Use:   "rule_delete",
	Short: "Delete an eventwatch rule by name",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := AppAPI.DeleteRule(eventwatchName); err != nil {
			return err
		}
		cmd.PrintErrf("Rule '%s' deleted successfully\n", eventwatchName)
		return nil
	},
}

var eventwatchGroupDeleteCmd = &cobra.Command{
	Use:   "group_delete",
	Short: "Delete an eventwatch rule group by name",
	RunE: func(cmd *cobra.Command, args []string) error {
		count, err := AppAPI.DeleteRuleGroup(eventwatchGroup)
		if err != nil {
			return err
		}
		cmd.PrintErrf("Group '%s' deleted successfully (%d rules)\n", eventwatchGroup, count)
		return nil
	},
}

// loadEventwatchRule reads an EventWatchBucket JSON definition from the --content
// flag, which accepts inline JSON, an "@path" file reference, or "-" for stdin.
func loadEventwatchRule(cmd *cobra.Command) (*model.EventWatchBucket, error) {
	var raw string
	switch {
	case eventwatchContent == "-":
		b, err := io.ReadAll(cmd.InOrStdin())
		if err != nil {
			return nil, fmt.Errorf("failed to read from stdin: %w", err)
		}
		raw = string(b)
	case len(eventwatchContent) > 1 && eventwatchContent[0] == '@':
		filePath := eventwatchContent[1:]
		b, err := os.ReadFile(filePath)
		if err != nil {
			return nil, fmt.Errorf("failed to read file '%s': %w", filePath, err)
		}
		raw = string(b)
	default:
		raw = eventwatchContent
	}

	if len(raw) == 0 {
		return nil, fmt.Errorf("rule content is empty")
	}

	var entry model.EventWatchBucket
	if err := json.Unmarshal([]byte(raw), &entry); err != nil {
		return nil, fmt.Errorf("failed to parse rule JSON: %w", err)
	}
	if entry.Name == "" {
		return nil, fmt.Errorf("rule name field missing")
	}
	return &entry, nil
}

func printEventwatchHits(cmd *cobra.Command, resp *model.ElasticSearchResult, sourceType string) error {
	if resp.Hits == nil || len(resp.Hits.Hits) == 0 {
		cmd.PrintErrln("No hits found.")
		return nil
	}
	for _, hit := range resp.Hits.Hits {
		switch sourceType {
		case "BehaviorSummary":
			var src model.BehaviorSummary
			if err := json.Unmarshal(hit.Source, &src); err != nil {
				cmd.PrintErrf("skip hit %s: invalid _source: %v\n", hit.ID, err)
				continue
			}
			cmd.Printf("Key: %s, RiskScore: %d\n", src.Key, src.RiskScore)
		case "BehaviorEvent":
			var src model.BehaviorEvent
			if err := json.Unmarshal(hit.Source, &src); err != nil {
				cmd.PrintErrf("skip hit %s: invalid _source: %v\n", hit.ID, err)
				continue
			}
			cmd.Printf("Key: %s, RiskScore: %d\n", src.Key, src.RiskScore)
		case "BehaviorRule":
			var src model.EventWatchBucket
			if err := json.Unmarshal(hit.Source, &src); err != nil {
				cmd.PrintErrf("skip hit %s: invalid _source: %v\n", hit.ID, err)
				continue
			}
			cmd.Printf("Name: %s, Group: %s, Repository: %s\n", src.Name, src.Group, src.Repository)
		default:
			cmd.PrintErrf("skip hit %s: unknown source type %q\n", hit.ID, sourceType)
		}
	}
	return nil
}

func init() {
	RootCmd.AddCommand(eventwatchCmd)
	eventwatchCmd.AddCommand(eventwatchSummarySearchCmd, eventwatchTimelineSearchCmd, eventwatchRuleSearchCmd)
	eventwatchCmd.AddCommand(
		eventwatchRuleListCmd,
		eventwatchRuleGetCmd,
		eventwatchRuleAddCmd,
		eventwatchRuleUpdateCmd,
		eventwatchRuleToggleCmd,
		eventwatchRuleDeleteCmd,
		eventwatchGroupDeleteCmd,
	)

	eventwatchSummarySearchCmd.Flags().StringVar(&eventwatchQuery, "query", "", "Search query")
	eventwatchSummarySearchCmd.Flags().Int64Var(&eventwatchFrom, "from", 0, "Range start (Unix ms); default: 1 hour ago")
	eventwatchSummarySearchCmd.Flags().Int64Var(&eventwatchTo, "to", 0, "Range end (Unix ms); default: now")

	eventwatchTimelineSearchCmd.Flags().StringVar(&eventwatchQuery, "query", "", "Search query")
	eventwatchTimelineSearchCmd.Flags().Int64Var(&eventwatchFrom, "from", 0, "Range start (Unix ms); default: 1 hour ago")
	eventwatchTimelineSearchCmd.Flags().Int64Var(&eventwatchTo, "to", 0, "Range end (Unix ms); default: now")

	eventwatchRuleSearchCmd.Flags().StringVar(&eventwatchQuery, "query", "", "Search query")

	eventwatchRuleGetCmd.Flags().StringVar(&eventwatchName, "name", "", "Rule name (empty returns the first rule)")

	eventwatchRuleAddCmd.Flags().StringVar(&eventwatchContent, "content", "", "Rule JSON, '@path' file, or '-' for stdin")
	_ = eventwatchRuleAddCmd.MarkFlagRequired("content")

	eventwatchRuleUpdateCmd.Flags().StringVar(&eventwatchContent, "content", "", "Rule JSON, '@path' file, or '-' for stdin")
	_ = eventwatchRuleUpdateCmd.MarkFlagRequired("content")

	eventwatchRuleToggleCmd.Flags().StringVar(&eventwatchName, "name", "", "Rule name")
	_ = eventwatchRuleToggleCmd.MarkFlagRequired("name")

	eventwatchRuleDeleteCmd.Flags().StringVar(&eventwatchName, "name", "", "Rule name")
	_ = eventwatchRuleDeleteCmd.MarkFlagRequired("name")

	eventwatchGroupDeleteCmd.Flags().StringVar(&eventwatchGroup, "group", "", "Rule group name")
	_ = eventwatchGroupDeleteCmd.MarkFlagRequired("group")
}
