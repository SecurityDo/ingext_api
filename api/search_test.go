package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/SecurityDo/ingext_api/client"
	fsb "github.com/SecurityDo/ingext_api/fsb"
	"github.com/SecurityDo/ingext_api/model"
)

// capturedLakeSearch is the lake_search kargs as it arrives on the wire.
type capturedLakeSearch struct {
	Index   string                    `json:"index"`
	Options *capturedLakeSearchOption `json:"options"`
}

type capturedLakeSearchOption struct {
	SearchStr   string              `json:"searchStr"`
	RangeFrom   int64               `json:"range_from"`
	RangeTo     int64               `json:"range_to"`
	FetchLimit  int                 `json:"fetchLimit"`
	FetchOffset int                 `json:"fetchOffset"`
	SortField   string              `json:"sortField"`
	SortOrder   string              `json:"sortOrder"`
	Facets      *model.FacetsOption `json:"facets"`
}

func newSearchServiceForTest(t *testing.T, handler http.HandlerFunc) *SearchService {
	t.Helper()
	ts := httptest.NewServer(handler)
	t.Cleanup(ts.Close)

	flClient := client.NewIngextClient(ts.URL, "", false, nil)
	return NewSearchService(flClient)
}

// newLakeSearchRecorder serves a fixed response and records the kargs of the
// single call it expects.
func newLakeSearchRecorder(t *testing.T, response interface{}) (*SearchService, *capturedLakeSearch, *json.RawMessage) {
	t.Helper()
	captured := &capturedLakeSearch{}
	rawKargs := new(json.RawMessage)

	svc := newSearchServiceForTest(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/ds/lake_search" {
			t.Fatalf("unexpected path %s", r.URL.Path)
		}
		var req fsb.CallRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatalf("failed to decode request: %v", err)
		}
		if req.Kargs == nil {
			t.Fatalf("missing kargs in call request")
		}
		*rawKargs = json.RawMessage(req.Kargs.GetBytes())
		if err := json.Unmarshal(req.Kargs.GetBytes(), captured); err != nil {
			t.Fatalf("failed to decode kargs: %v", err)
		}

		body, err := json.Marshal(response)
		if err != nil {
			t.Fatalf("failed to marshal response: %v", err)
		}
		respPayload := fsb.CallResponse{Verdict: "OK", Response: fsb.NewJNodeByte(body)}
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(respPayload); err != nil {
			t.Fatalf("failed to write response: %v", err)
		}
	})
	return svc, captured, rawKargs
}

// TestSearchService_LakeSearchRequest locks in the wire shape, including the
// defaults the endpoint has none of: an empty sort field fails its query parser,
// and a facet of size 0 comes back with no buckets at all.
func TestSearchService_LakeSearchRequest(t *testing.T) {
	svc, captured, rawKargs := newLakeSearchRecorder(t, model.LakeSearchResponse{Took: 12})

	_, err := svc.LakeSearch(&LakeSearchOptions{
		Index:     "managed-Office365",
		Query:     "71.178.173.2",
		RangeFrom: 1787766069700,
		RangeTo:   1787769669701,
		Facets: []*model.FacetEntry{
			{Title: "Username", Field: "@fields.UserId"},
			{Field: "@fields.ClientIP", Size: 5},
		},
		MustFilters:    []*model.FilterEntry{{Field: "@source", Terms: []string{"Audit.AzureActiveDirectory"}}},
		MustNotFilters: []*model.FilterEntry{{Field: "@fields.UserType", Terms: []string{"0"}}},
		FetchLimit:     1000,
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if captured.Index != "managed-Office365" {
		t.Fatalf("unexpected index %q", captured.Index)
	}
	// dataType, partition and dayIndex are declared by the endpoint but never
	// read, so they stay off the wire.
	var keys map[string]json.RawMessage
	if err := json.Unmarshal(*rawKargs, &keys); err != nil {
		t.Fatalf("failed to decode kargs: %v", err)
	}
	if len(keys) != 2 || keys["index"] == nil || keys["options"] == nil {
		t.Fatalf("kargs should carry only index and options, got %v", keys)
	}
	opts := captured.Options
	if opts == nil {
		t.Fatalf("missing options in kargs")
	}
	if opts.SearchStr != "71.178.173.2" {
		t.Fatalf("unexpected searchStr %q", opts.SearchStr)
	}
	if opts.RangeFrom != 1787766069700 || opts.RangeTo != 1787769669701 {
		t.Fatalf("time range lost: %+v", opts)
	}
	if opts.FetchLimit != 1000 || opts.FetchOffset != 0 {
		t.Fatalf("unexpected pagination: %+v", opts)
	}
	if opts.SortField != "@timestamp" || opts.SortOrder != "desc" {
		t.Fatalf("sort defaults not applied: %+v", opts)
	}
	if opts.Facets == nil {
		t.Fatalf("facets must always be sent: the endpoint walks them without a nil check")
	}
	if len(opts.Facets.Facets) != 2 {
		t.Fatalf("unexpected facets: %+v", opts.Facets.Facets)
	}
	if got := opts.Facets.Facets[0]; got.Title != "Username" || got.Field != "@fields.UserId" || got.Size != 20 {
		t.Fatalf("facet default size or title lost: %+v", got)
	}
	if got := opts.Facets.Facets[1]; got.Title != "@fields.ClientIP" || got.Size != 5 {
		t.Fatalf("explicit facet size overridden, or title not defaulted to the field: %+v", got)
	}
	if len(opts.Facets.MustFilters) != 1 || opts.Facets.MustFilters[0].Field != "@source" {
		t.Fatalf("unexpected mustFilters: %+v", opts.Facets.MustFilters)
	}
	if len(opts.Facets.MustNotFilters) != 1 || opts.Facets.MustNotFilters[0].Field != "@fields.UserType" {
		t.Fatalf("unexpected mustNotFilters: %+v", opts.Facets.MustNotFilters)
	}
}

// TestSearchService_LakeSearchSendsFacetsObject covers the crash case: the
// endpoint ranges over facets.mustFilters and friends without a nil check, so a
// search with no facets at all must still carry the object and its arrays.
func TestSearchService_LakeSearchSendsFacetsObject(t *testing.T) {
	svc, _, rawKargs := newLakeSearchRecorder(t, model.LakeSearchResponse{})

	_, err := svc.LakeSearch(&LakeSearchOptions{
		Index:     "managed-Office365",
		RangeFrom: 1,
		RangeTo:   2,
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	var payload struct {
		Options struct {
			FetchLimit int `json:"fetchLimit"`
			Facets     struct {
				Facets         []json.RawMessage       `json:"facets"`
				MustFilters    []json.RawMessage       `json:"mustFilters"`
				MustNotFilters []json.RawMessage       `json:"mustNotFilters"`
				DateFacets     []*model.DateFacetEntry `json:"dateFacets"`
			} `json:"facets"`
		} `json:"options"`
	}
	if err := json.Unmarshal(*rawKargs, &payload); err != nil {
		t.Fatalf("failed to decode kargs: %v", err)
	}
	facets := payload.Options.Facets
	if facets.Facets == nil || facets.MustFilters == nil || facets.MustNotFilters == nil {
		t.Fatalf("facet arrays must be sent as empty arrays, not null: %+v", facets)
	}
	// The platform keys the histogram aggregation by the name of the first date
	// facet. Without one it lands under "" instead of "dateHistogram".
	if len(facets.DateFacets) != 1 || facets.DateFacets[0].Name != model.DateHistogramAggregation {
		t.Fatalf("date facet must be named %q, got %+v", model.DateHistogramAggregation, facets.DateFacets)
	}
	// A fetchLimit of 0 reaches the platform and panics its search node, so an
	// unset limit has to become the default instead.
	if payload.Options.FetchLimit != 100 {
		t.Fatalf("expected an unset fetchLimit to default to 100, got %d", payload.Options.FetchLimit)
	}
}

// TestSearchService_LakeSearchResponse decodes a trimmed copy of a real
// response. The histogram key matters: it is an epoch millisecond too large to
// survive a float64 round trip.
func TestSearchService_LakeSearchResponse(t *testing.T) {
	const body = `{
	  "took": 1112,
	  "hits": {
	    "total": 4,
	    "hits": [{"_source": {"@timestamp": 1787767995000}, "_id": "2erng0"}],
	    "ascend": false,
	    "limit": 0
	  },
	  "aggregations": {
	    "@fields.UserId": {"TokenEntity": "", "buckets": [
	      {"key": "kun@fluencysecurity.com", "doc_count": 2},
	      {"key": "howard@fluencysecurity.com", "doc_count": 1}
	    ]},
	    "@behaviors": {"TokenEntity": "", "buckets": []},
	    "dateHistogram": {"TokenEntity": "", "buckets": [
	      {"doc_count": 0, "key": 1787766069700},
	      {"doc_count": 1, "key": 1787766789700}
	    ]}
	  },
	  "query": null,
	  "cost": 0.0021,
	  "total": 4,
	  "filtered": 172
	}`

	svc := newSearchServiceForTest(t, func(w http.ResponseWriter, r *http.Request) {
		respPayload := fsb.CallResponse{Verdict: "OK", Response: fsb.NewJNodeByte([]byte(body))}
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(respPayload); err != nil {
			t.Fatalf("failed to write response: %v", err)
		}
	})

	resp, err := svc.LakeSearch(&LakeSearchOptions{Index: "managed-Office365", RangeFrom: 1, RangeTo: 2, FetchLimit: 10})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if resp.Took != 1112 || resp.Total != 4 || resp.Filtered != 172 || resp.Cost != 0.0021 {
		t.Fatalf("unexpected search totals: %+v", resp)
	}
	if resp.Hits == nil || resp.Hits.Total != 4 || len(resp.Hits.Hits) != 1 {
		t.Fatalf("unexpected hits: %+v", resp.Hits)
	}
	if resp.Hits.Hits[0].ID != "2erng0" {
		t.Fatalf("hit id lost: %+v", resp.Hits.Hits[0])
	}

	buckets, err := resp.Aggregations["@fields.UserId"].FacetBuckets()
	if err != nil {
		t.Fatalf("failed to decode facet buckets: %v", err)
	}
	if len(buckets) != 2 || buckets[0].Key != "kun@fluencysecurity.com" || buckets[0].DocCount != 2 {
		t.Fatalf("unexpected facet buckets: %+v", buckets)
	}

	// A facet the index computed but found nothing for is present and empty,
	// which is not the same as absent.
	empty, err := resp.Aggregations["@behaviors"].FacetBuckets()
	if err != nil {
		t.Fatalf("failed to decode empty facet: %v", err)
	}
	if len(empty) != 0 {
		t.Fatalf("expected no buckets, got %+v", empty)
	}

	slots, err := resp.Aggregations[model.DateHistogramAggregation].HistogramBuckets()
	if err != nil {
		t.Fatalf("failed to decode histogram buckets: %v", err)
	}
	if len(slots) != 2 || slots[0].Key != 1787766069700 || slots[0].DocCount != 0 {
		t.Fatalf("unexpected histogram slots: %+v", slots[0])
	}
	if slots[1].Key != 1787766789700 || slots[1].DocCount != 1 {
		t.Fatalf("histogram slot key or count lost: %+v", slots[1])
	}

	// The two bucket shapes do not interchange: reading the histogram as a term
	// facet has to fail rather than return zero-valued keys.
	if _, err := resp.Aggregations[model.DateHistogramAggregation].FacetBuckets(); err == nil {
		t.Fatalf("expected an error decoding histogram buckets as a term facet")
	}
}

// TestSearchService_LakeSearchValidation covers the requests that are rejected
// before any call is made: each one fails inside the search service with an
// error that never names the field at fault.
func TestSearchService_LakeSearchValidation(t *testing.T) {
	svc := newSearchServiceForTest(t, func(w http.ResponseWriter, r *http.Request) {
		t.Fatalf("no call should reach the platform")
	})

	cases := []struct {
		name string
		opts *LakeSearchOptions
	}{
		{"nil options", nil},
		{"no time range", &LakeSearchOptions{Index: "managed-Office365"}},
		{"open-ended range", &LakeSearchOptions{Index: "managed-Office365", RangeFrom: 100}},
		{"empty range", &LakeSearchOptions{Index: "managed-Office365", RangeFrom: 200, RangeTo: 100}},
		{"negative limit", &LakeSearchOptions{RangeFrom: 1, RangeTo: 2, FetchLimit: -1}},
		{"pagination over cap", &LakeSearchOptions{RangeFrom: 1, RangeTo: 2, FetchOffset: 4000, FetchLimit: 1001}},
		{"defaulted limit over cap", &LakeSearchOptions{RangeFrom: 1, RangeTo: 2, FetchOffset: 4950}},
		{"facet without field", &LakeSearchOptions{RangeFrom: 1, RangeTo: 2, Facets: []*model.FacetEntry{{Title: "Username"}}}},
		{"filter without field", &LakeSearchOptions{RangeFrom: 1, RangeTo: 2, MustFilters: []*model.FilterEntry{{Terms: []string{"x"}}}}},
		{"filter without terms", &LakeSearchOptions{RangeFrom: 1, RangeTo: 2, MustNotFilters: []*model.FilterEntry{{Field: "@source"}}}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if _, err := svc.LakeSearch(tc.opts); err == nil {
				t.Fatalf("expected an error")
			}
		})
	}

	// The cap itself is allowed - only crossing it is not.
	if _, err := (&LakeSearchOptions{RangeFrom: 1, RangeTo: 2, FetchOffset: 4000, FetchLimit: 1000}).toRequest(); err != nil {
		t.Fatalf("expected the pagination cap to be inclusive, got %v", err)
	}
}
