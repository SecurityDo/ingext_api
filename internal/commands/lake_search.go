package commands

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
	"text/tabwriter"
	"time"

	ingextAPI "github.com/SecurityDo/ingext_api/api"
	"github.com/SecurityDo/ingext_api/model"
	"github.com/spf13/cobra"
)

var (
	lakeSearchDatalake  string
	lakeSearchIndex     string
	lakeSearchQuery     string
	lakeSearchFrom      int64
	lakeSearchTo        int64
	lakeSearchFacets    []string
	lakeSearchFacetSize int
	lakeSearchMust      []string
	lakeSearchMustNot   []string
	lakeSearchLimit     int
	lakeSearchOffset    int
	lakeSearchSortField string
	lakeSearchSortOrder string
	lakeSearchHistogram bool
	lakeSearchJSON      bool
	lakeSearchOutput    string
)

var lakeSearchCmd = &cobra.Command{
	Use:   "search",
	Short: "Facet search over a datalake index",
	Long: `Run a facet search (lake_search) against one datalake index.

The search is a Lucene query over a time range, narrowed by must / must-not term
filters. Every --facet field is counted over whatever matches, and the counts
come back alongside the hits. An index also contributes its own configured
facets, so the response can carry counts that were not asked for.

--index takes the index name from 'ingext datalake list-index'. It is prefixed
with '<--datalake>-' unless it already contains a dash, so '--index Office365'
searches 'managed-Office365'.

A --facet is 'field' or 'Title=field'. A filter is 'field=term[,term...]', where
the terms of one filter are OR'ed and separate filters are AND'ed.

  ingext datalake search --index Office365 --query "71.178.173.2" \
      --facet "Username=@fields.UserId" --facet @fields.ClientIP \
      --must '@source=Audit.AzureActiveDirectory' --limit 20

With no --from/--to the range is the last hour. --limit is at least 1: the
platform has no hits-free search, and a limit of 0 crashes its search node.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if lakeSearchLimit < 1 {
			return fmt.Errorf("--limit must be at least 1: the platform has no facet-only search, and a limit of 0 crashes its search node")
		}

		from, to := lakeSearchFrom, lakeSearchTo
		if from == 0 && to == 0 {
			now := time.Now().UnixMilli()
			to = now
			from = now - int64(time.Hour/time.Millisecond)
		}

		facets, err := parseFacetSpecs(lakeSearchFacets, lakeSearchFacetSize)
		if err != nil {
			return err
		}
		mustFilters, err := parseFilterSpecs(lakeSearchMust, "--must")
		if err != nil {
			return err
		}
		mustNotFilters, err := parseFilterSpecs(lakeSearchMustNot, "--must-not")
		if err != nil {
			return err
		}

		resp, err := AppAPI.LakeSearch(&ingextAPI.LakeSearchOptions{
			Index:          qualifyLakeIndex(lakeSearchDatalake, lakeSearchIndex),
			Query:          lakeSearchQuery,
			RangeFrom:      from,
			RangeTo:        to,
			Facets:         facets,
			MustFilters:    mustFilters,
			MustNotFilters: mustNotFilters,
			FetchLimit:     lakeSearchLimit,
			FetchOffset:    lakeSearchOffset,
			SortField:      lakeSearchSortField,
			SortOrder:      lakeSearchSortOrder,
		})
		if err != nil {
			return err
		}

		if lakeSearchOutput != "" {
			jsonBytes, err := json.MarshalIndent(resp, "", "  ")
			if err != nil {
				return fmt.Errorf("marshal response: %w", err)
			}
			if err := os.WriteFile(lakeSearchOutput, jsonBytes, 0644); err != nil {
				return fmt.Errorf("write output file %s: %w", lakeSearchOutput, err)
			}
			cmd.PrintErrf("Response saved to %s\n", lakeSearchOutput)
		}

		if lakeSearchJSON {
			jsonBytes, err := json.MarshalIndent(resp, "", "  ")
			if err != nil {
				return fmt.Errorf("marshal response: %w", err)
			}
			cmd.Println(string(jsonBytes))
			return nil
		}

		printLakeSearchResponse(cmd, resp)
		return nil
	},
}

// qualifyLakeIndex turns an index name into the "<datalake>-<index>" name the
// endpoint indexes by. An index that already carries a datalake prefix is left
// alone, which is the same rule the platform applies elsewhere.
func qualifyLakeIndex(datalake, index string) string {
	if index == "" || strings.Contains(index, "-") {
		return index
	}
	if datalake == "" {
		return index
	}
	return datalake + "-" + index
}

// parseFacetSpecs parses "field" or "Title=field" into facet entries.
func parseFacetSpecs(specs []string, size int) ([]*model.FacetEntry, error) {
	var facets []*model.FacetEntry
	for _, spec := range specs {
		title, field := "", strings.TrimSpace(spec)
		if name, rest, found := strings.Cut(spec, "="); found {
			title = strings.TrimSpace(name)
			field = strings.TrimSpace(rest)
		}
		if field == "" {
			return nil, fmt.Errorf("--facet %q: field is missing", spec)
		}
		facets = append(facets, &model.FacetEntry{Title: title, Field: field, Size: size})
	}
	return facets, nil
}

// parseFilterSpecs parses "field=term[,term...]" into filter entries.
func parseFilterSpecs(specs []string, flag string) ([]*model.FilterEntry, error) {
	var filters []*model.FilterEntry
	for _, spec := range specs {
		field, termList, found := strings.Cut(spec, "=")
		field = strings.TrimSpace(field)
		if !found || field == "" {
			return nil, fmt.Errorf("%s %q: expected field=term[,term...]", flag, spec)
		}
		var terms []string
		for _, term := range strings.Split(termList, ",") {
			if term = strings.TrimSpace(term); term != "" {
				terms = append(terms, term)
			}
		}
		if len(terms) == 0 {
			return nil, fmt.Errorf("%s %q: at least one term is required", flag, spec)
		}
		filters = append(filters, &model.FilterEntry{Field: field, Terms: terms})
	}
	return filters, nil
}

func printLakeSearchResponse(cmd *cobra.Command, resp *model.LakeSearchResponse) {
	if resp == nil {
		cmd.PrintErrln("Empty response.")
		return
	}

	var returned int
	var total uint64
	if resp.Hits != nil {
		returned = len(resp.Hits.Hits)
		total = resp.Hits.Total
	}
	cmd.Printf("%d hits (%d returned, %d filtered, %d ms, cost %g)\n",
		total, returned, resp.Filtered, resp.Took, resp.Cost)

	// Facets first: they are usually what the search was for, and a page of
	// events would push them off the screen.
	names := make([]string, 0, len(resp.Aggregations))
	for name := range resp.Aggregations {
		if name == model.DateHistogramAggregation && !lakeSearchHistogram {
			continue
		}
		names = append(names, name)
	}
	sort.Strings(names)

	w := tabwriter.NewWriter(cmd.OutOrStdout(), 0, 0, 2, ' ', 0)
	for _, name := range names {
		agg := resp.Aggregations[name]
		if name == model.DateHistogramAggregation {
			buckets, err := agg.HistogramBuckets()
			if err != nil {
				cmd.PrintErrf("\n%s: %v\n", name, err)
				continue
			}
			fmt.Fprintf(w, "\n%s\t\n", name)
			for _, bucket := range buckets {
				if bucket.DocCount == 0 {
					continue
				}
				slot := time.UnixMilli(bucket.Key).Format("2006-01-02 15:04:05")
				fmt.Fprintf(w, "  %s\t%d\n", slot, bucket.DocCount)
			}
			continue
		}
		buckets, err := agg.FacetBuckets()
		if err != nil {
			cmd.PrintErrf("\n%s: %v\n", name, err)
			continue
		}
		if len(buckets) == 0 {
			continue
		}
		fmt.Fprintf(w, "\n%s\t\n", name)
		for _, bucket := range buckets {
			fmt.Fprintf(w, "  %s\t%d\n", bucket.Key, bucket.DocCount)
		}
	}
	w.Flush()

	if returned == 0 {
		return
	}
	cmd.Println()
	for _, hit := range resp.Hits.Hits {
		// One event per line: the documents of two indexes share no columns, so
		// there is no table to print. Use --json or --output for the full shape.
		cmd.Printf("%s %s\n", hit.ID, compactJSON(hit.Source))
	}
}

// compactJSON strips the indentation of a hit source so one event stays on one
// line. Unparsable input is returned as it arrived rather than dropped.
func compactJSON(raw json.RawMessage) string {
	var buf bytes.Buffer
	if err := json.Compact(&buf, raw); err != nil {
		return string(raw)
	}
	return buf.String()
}

func init() {
	lakeCmd.AddCommand(lakeSearchCmd)

	lakeSearchCmd.Flags().StringVar(&lakeSearchDatalake, "datalake", "managed", "datalake the index belongs to")
	lakeSearchCmd.Flags().StringVar(&lakeSearchIndex, "index", "", "datalake index to search")
	lakeSearchCmd.Flags().StringVarP(&lakeSearchQuery, "query", "q", "", "lucene query string (empty matches everything in range)")
	lakeSearchCmd.Flags().Int64Var(&lakeSearchFrom, "from", 0, "range start, epoch milliseconds (default: one hour ago)")
	lakeSearchCmd.Flags().Int64Var(&lakeSearchTo, "to", 0, "range end, epoch milliseconds (default: now)")
	lakeSearchCmd.Flags().StringArrayVar(&lakeSearchFacets, "facet", nil, "field to count terms on, as 'field' or 'Title=field' (repeatable)")
	lakeSearchCmd.Flags().IntVar(&lakeSearchFacetSize, "facet-size", 20, "number of terms to return per facet")
	lakeSearchCmd.Flags().StringArrayVar(&lakeSearchMust, "must", nil, "term filter to require, as 'field=term[,term...]' (repeatable)")
	lakeSearchCmd.Flags().StringArrayVar(&lakeSearchMustNot, "must-not", nil, "term filter to exclude, as 'field=term[,term...]' (repeatable)")
	lakeSearchCmd.Flags().IntVar(&lakeSearchLimit, "limit", 10, "hits to return (minimum 1)")
	lakeSearchCmd.Flags().IntVar(&lakeSearchOffset, "offset", 0, "hits to skip; offset+limit cannot exceed "+strconv.Itoa(ingextAPI.LakeSearchMaxPagination))
	lakeSearchCmd.Flags().StringVar(&lakeSearchSortField, "sort-field", "", "field to sort hits by (default: @timestamp)")
	lakeSearchCmd.Flags().StringVar(&lakeSearchSortOrder, "sort-order", "", "asc or desc (default: desc)")
	lakeSearchCmd.Flags().BoolVar(&lakeSearchHistogram, "histogram", false, "print the non-empty date histogram slots")
	lakeSearchCmd.Flags().BoolVar(&lakeSearchJSON, "json", false, "emit the full response as JSON on stdout")
	lakeSearchCmd.Flags().StringVar(&lakeSearchOutput, "output", "", "save the full JSON response to a file")

	_ = lakeSearchCmd.MarkFlagRequired("index")
}
