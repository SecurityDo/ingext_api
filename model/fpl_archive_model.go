package model

import (
	"encoding/json"
	"time"
)

// The scheduled-report archive: a daily scheduled report with archive enabled
// runs once a day for exactly the previous day (in the archive timezone) and
// keeps each result for 12 months. Days are "YYYY-MM-DD" in that timezone and
// ranges include both ends.

// ReportArchive is the archive block of an FPL v2 report.
type ReportArchive struct {
	Enabled bool `json:"enabled"`
	// Timezone is the day boundary (IANA name); empty means UTC.
	Timezone string `json:"timezone,omitempty"`
	// ArgFrom and ArgTo name the script's window arguments; empty means "from"/"to".
	ArgFrom string `json:"argFrom,omitempty"`
	ArgTo   string `json:"argTo,omitempty"`
}

// ReportArchiveRequest is the kargs of every *_fpl_report_archive function;
// each reads the fields it needs.
type ReportArchiveRequest struct {
	ID      int64    `json:"id,omitempty"`
	Name    string   `json:"name,omitempty"`
	Day     string   `json:"day,omitempty"`
	From    string   `json:"from,omitempty"`
	To      string   `json:"to,omitempty"`
	State   string   `json:"state,omitempty"`
	Objects []string `json:"objects,omitempty"`
	Limit   int      `json:"limit,omitempty"`
	Force   bool     `json:"force,omitempty"`
}

// ReportArchiveEntry is one archived day.
type ReportArchiveEntry struct {
	ID          int64      `json:"id,omitempty"`
	Account     string     `json:"account"`
	ReportName  string     `json:"reportName"`
	TaskID      uint       `json:"taskID"`
	State       string     `json:"state"` // running, completed, aborted
	ErrorMsg    string     `json:"errorMsg,omitempty"`
	RangeFrom   time.Time  `json:"rangeFrom"`
	RangeTo     time.Time  `json:"rangeTo"`
	CreatedOn   time.Time  `json:"createdOn"`
	CompletedOn *time.Time `json:"completedOn,omitempty"`
	Took        int64      `json:"took"`
	Cost        float64    `json:"cost"`
	// Total is the number of output rows.
	Total       int64                `json:"total"`
	ResultBytes int64                `json:"resultBytes"`
	ResultPath  string               `json:"resultPath,omitempty"`
	Env         *ReportArchiveRunEnv `json:"env,omitempty"`
}

// ReportArchiveRunEnv is what the day was run with.
type ReportArchiveRunEnv struct {
	Timezone      string          `json:"timezone"`
	ScriptHash    string          `json:"scriptHash,omitempty"`
	ProcessorName string          `json:"processorName,omitempty"`
	Arguments     json.RawMessage `json:"arguments,omitempty"`
	ReportConfig  json.RawMessage `json:"reportConfig,omitempty"`
}

// ReportArchiveCoverage is what a requested day range holds.
type ReportArchiveCoverage struct {
	FirstDay     string   `json:"firstDay"`
	LastDay      string   `json:"lastDay"`
	Expected     int      `json:"expected"`
	Completed    int      `json:"completed"`
	Aborted      int      `json:"aborted"`
	Running      int      `json:"running"`
	Missing      int      `json:"missing"`
	MissingDays  []string `json:"missingDays,omitempty"`
	AbortedDays  []string `json:"abortedDays,omitempty"`
	RunningDays  []string `json:"runningDays,omitempty"`
	ScriptHashes []string `json:"scriptHashes,omitempty"`
}

// ArchivedReport is one report in the list without a name.
type ArchivedReport struct {
	Name      string `json:"name"`
	Enabled   bool   `json:"enabled"`
	Timezone  string `json:"timezone"`
	FirstDay  string `json:"firstDay,omitempty"`
	LastDay   string `json:"lastDay,omitempty"`
	Days      int64  `json:"days"`
	Completed int64  `json:"completed"`
}

// ReportArchiveBackfillStatus is a backfill in progress.
type ReportArchiveBackfillStatus struct {
	Name      string    `json:"name"`
	Days      []string  `json:"days"`
	Done      int       `json:"done"`
	Failed    int       `json:"failed"`
	Current   string    `json:"current,omitempty"`
	StartedBy string    `json:"startedBy,omitempty"`
	StartedOn time.Time `json:"startedOn"`
}

// ListReportArchiveResponse covers both shapes of list_fpl_report_archive:
// without a name only Reports is set; with one, the rest.
type ListReportArchiveResponse struct {
	Reports  []*ArchivedReport            `json:"reports,omitempty"`
	Name     string                       `json:"name,omitempty"`
	Enabled  bool                         `json:"enabled,omitempty"`
	Timezone string                       `json:"timezone,omitempty"`
	Coverage *ReportArchiveCoverage       `json:"coverage,omitempty"`
	Entries  []*ReportArchiveEntry        `json:"entries,omitempty"`
	Backfill *ReportArchiveBackfillStatus `json:"backfill,omitempty"`
}

type GetReportArchiveResponse struct {
	Entry *ReportArchiveEntry `json:"entry"`
}

// ReportArchiveDay is one downloaded day; only completed days are returned.
type ReportArchiveDay struct {
	Entry  *ReportArchiveEntry `json:"entry"`
	Result *MetricFPLResponse  `json:"result"`
}

type DownloadReportArchiveResponse struct {
	Name     string                 `json:"name"`
	Timezone string                 `json:"timezone"`
	Coverage *ReportArchiveCoverage `json:"coverage,omitempty"`
	Days     []*ReportArchiveDay    `json:"days"`
}

type BackfillReportArchiveResponse struct {
	Name     string   `json:"name"`
	Timezone string   `json:"timezone"`
	Queued   []string `json:"queued"`
	Skipped  []string `json:"skipped"`
}

type ReindexReportArchiveResult struct {
	Restored int      `json:"restored"`
	Existing int      `json:"existing"`
	Failed   int      `json:"failed"`
	Errors   []string `json:"errors,omitempty"`
}

type ReindexReportArchiveResponse struct {
	Name   string                      `json:"name"`
	Result *ReindexReportArchiveResult `json:"result"`
}
