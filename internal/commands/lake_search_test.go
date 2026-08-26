package commands

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/SecurityDo/ingext_api/model"
	"github.com/spf13/cobra"
)

func TestQualifyLakeIndex(t *testing.T) {
	cases := []struct {
		datalake string
		index    string
		want     string
	}{
		{"managed", "Office365", "managed-Office365"},
		{"managed", "managed-Office365", "managed-Office365"},
		{"managed", "custom-Office365", "custom-Office365"},
		{"", "Office365", "Office365"},
		{"managed", "", ""},
	}
	for _, tc := range cases {
		if got := qualifyLakeIndex(tc.datalake, tc.index); got != tc.want {
			t.Fatalf("qualifyLakeIndex(%q, %q) = %q, want %q", tc.datalake, tc.index, got, tc.want)
		}
	}
}

func TestParseFacetSpecs(t *testing.T) {
	facets, err := parseFacetSpecs([]string{"@fields.ClientIP", "Username=@fields.UserId"}, 20)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(facets) != 2 {
		t.Fatalf("unexpected facets: %+v", facets)
	}
	if facets[0].Field != "@fields.ClientIP" || facets[0].Title != "" || facets[0].Size != 20 {
		t.Fatalf("unexpected facet: %+v", facets[0])
	}
	if facets[1].Field != "@fields.UserId" || facets[1].Title != "Username" {
		t.Fatalf("unexpected facet: %+v", facets[1])
	}
	if _, err := parseFacetSpecs([]string{"Username="}, 20); err == nil {
		t.Fatalf("expected an error for a spec with no field")
	}
}

func TestParseFilterSpecs(t *testing.T) {
	filters, err := parseFilterSpecs([]string{"@source=Audit.AzureActiveDirectory, Audit.Exchange"}, "--must")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(filters) != 1 || filters[0].Field != "@source" {
		t.Fatalf("unexpected filters: %+v", filters)
	}
	if len(filters[0].Terms) != 2 || filters[0].Terms[1] != "Audit.Exchange" {
		t.Fatalf("terms not split or trimmed: %+v", filters[0].Terms)
	}
	for _, spec := range []string{"@source", "=Audit.Exchange", "@source="} {
		if _, err := parseFilterSpecs([]string{spec}, "--must"); err == nil {
			t.Fatalf("expected an error for %q", spec)
		}
	}
}

func TestPrintLakeSearchResponse(t *testing.T) {
	agg := func(buckets string) *model.SearchAggregate {
		var raw []json.RawMessage
		if err := json.Unmarshal([]byte(buckets), &raw); err != nil {
			t.Fatalf("bad test fixture: %v", err)
		}
		return &model.SearchAggregate{Buckets: raw}
	}

	resp := &model.LakeSearchResponse{
		Took:     1112,
		Cost:     0.0021,
		Total:    4,
		Filtered: 172,
		Hits: &model.LegacySearchHits{
			Total: 4,
			Hits: []*model.LegacySearchHit{
				{ID: "2erng0", Source: json.RawMessage("{\n  \"@source\": \"Audit.AzureActiveDirectory\"\n}")},
			},
		},
		Aggregations: map[string]*model.SearchAggregate{
			"@fields.UserId": agg(`[{"key":"kun@fluencysecurity.com","doc_count":2}]`),
			"@behaviors":     agg(`[]`),
			"dateHistogram":  agg(`[{"key":1787766789700,"doc_count":1}]`),
		},
	}

	cmd := &cobra.Command{}
	out := &bytes.Buffer{}
	cmd.SetOut(out)
	cmd.SetErr(out)

	lakeSearchHistogram = false
	printLakeSearchResponse(cmd, resp)
	got := out.String()

	if !strings.Contains(got, "4 hits (1 returned, 172 filtered, 1112 ms, cost 0.0021)") {
		t.Fatalf("summary line missing from output:\n%s", got)
	}
	if !strings.Contains(got, "kun@fluencysecurity.com") {
		t.Fatalf("facet buckets missing from output:\n%s", got)
	}
	// A facet with no buckets is noise in a terminal, and the histogram is only
	// printed when asked for.
	if strings.Contains(got, "@behaviors") {
		t.Fatalf("empty facet should not be printed:\n%s", got)
	}
	if strings.Contains(got, "dateHistogram") {
		t.Fatalf("histogram printed without --histogram:\n%s", got)
	}
	// One event per line, so grep and head stay useful.
	if !strings.Contains(got, `2erng0 {"@source":"Audit.AzureActiveDirectory"}`) {
		t.Fatalf("hit not printed on one line:\n%s", got)
	}

	out.Reset()
	lakeSearchHistogram = true
	t.Cleanup(func() { lakeSearchHistogram = false })
	printLakeSearchResponse(cmd, resp)
	if !strings.Contains(out.String(), "dateHistogram") {
		t.Fatalf("histogram missing with --histogram:\n%s", out.String())
	}
}

// TestPrintLakeSearchResponseBadBuckets: a facet whose buckets do not decode is
// reported and skipped, not fatal - the rest of the response is still worth
// printing.
func TestPrintLakeSearchResponseBadBuckets(t *testing.T) {
	resp := &model.LakeSearchResponse{
		Aggregations: map[string]*model.SearchAggregate{
			"@fields.UserId": {Buckets: []json.RawMessage{json.RawMessage(`{"key":42}`)}},
		},
	}
	cmd := &cobra.Command{}
	out := &bytes.Buffer{}
	cmd.SetOut(out)
	cmd.SetErr(out)

	printLakeSearchResponse(cmd, resp)
	if !strings.Contains(out.String(), "@fields.UserId") {
		t.Fatalf("expected the bad facet to be reported:\n%s", out.String())
	}
}
