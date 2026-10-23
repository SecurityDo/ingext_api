package api

import (
	"github.com/SecurityDo/ingext_api/model"
)

// ListReportArchive calls list_fpl_report_archive. An empty Name lists the
// archived reports; with a Name it returns that report's days in From..To
// (default the last 31 days) with coverage and any running backfill.
func (s *FPLService) ListReportArchive(req *model.ReportArchiveRequest) (*model.ListReportArchiveResponse, error) {
	var resp model.ListReportArchiveResponse
	if err := s.call("list_fpl_report_archive", req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// GetReportArchive calls get_fpl_report_archive for one day, by ID or by
// Name + Day.
func (s *FPLService) GetReportArchive(req *model.ReportArchiveRequest) (*model.ReportArchiveEntry, error) {
	var resp model.GetReportArchiveResponse
	if err := s.call("get_fpl_report_archive", req, &resp); err != nil {
		return nil, err
	}
	return resp.Entry, nil
}

// DownloadReportArchive calls download_fpl_report_archive: one day (ID, or
// Name + Day) or a range (Name + From [+ To], at most 92 days). Objects keeps
// only the named result objects and Limit caps rows per table.
func (s *FPLService) DownloadReportArchive(req *model.ReportArchiveRequest) (*model.DownloadReportArchiveResponse, error) {
	var resp model.DownloadReportArchiveResponse
	if err := s.call("download_fpl_report_archive", req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// BackfillReportArchive calls backfill_fpl_report_archive for past days (Day,
// or From..To, at most 92, before today). It returns at once; the days run in
// the background and ListReportArchive shows progress. Completed days are
// skipped unless Force.
func (s *FPLService) BackfillReportArchive(req *model.ReportArchiveRequest) (*model.BackfillReportArchiveResponse, error) {
	var resp model.BackfillReportArchiveResponse
	if err := s.call("backfill_fpl_report_archive", req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// ReindexReportArchive calls reindex_fpl_report_archive: rebuild the index rows
// from the archived files, for one report (Name) or all of them.
func (s *FPLService) ReindexReportArchive(req *model.ReportArchiveRequest) (*model.ReindexReportArchiveResponse, error) {
	var resp model.ReindexReportArchiveResponse
	if err := s.call("reindex_fpl_report_archive", req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}
