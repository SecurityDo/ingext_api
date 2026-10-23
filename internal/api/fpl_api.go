package api

import (
	"fmt"

	fluencyAPI "github.com/SecurityDo/ingext_api/api"
	"github.com/SecurityDo/ingext_api/model"
)

// RunReport calls the run_fplv2_report API with the given request and returns the task id.
func (c *Client) RunReport(req *model.RunFPLV2Report) (uint, error) {
	svc := fluencyAPI.NewFPLService(c.ingextClient)
	id, err := svc.RunReport(req)
	if err != nil {
		c.Logger.Error("failed to run FPL report", "error", err)
		return 0, fmt.Errorf("failed to run FPL report: %w", err)
	}
	return id, nil
}

// GetTaskByID calls the get_fpl_task API with the given id and returns the task response.
func (c *Client) GetTaskByID(id uint) (*model.GetTaskResponse, error) {
	svc := fluencyAPI.NewFPLService(c.ingextClient)
	resp, err := svc.GetTaskByID(id)
	if err != nil {
		c.Logger.Error("failed to get FPL task", "id", id, "error", err)
		return nil, fmt.Errorf("failed to get FPL task: %w", err)
	}
	return resp, nil
}

// GetResultsByID calls the get_fpl_result API with the given id and returns the metric result.
func (c *Client) GetResultsByID(id uint) (*model.GetMetricFPLResultResponse, error) {
	svc := fluencyAPI.NewFPLService(c.ingextClient)
	resp, err := svc.GetResultsByID(id)
	if err != nil {
		c.Logger.Error("failed to get FPL result", "id", id, "error", err)
		return nil, fmt.Errorf("failed to get FPL result: %w", err)
	}
	return resp, nil
}

// ListReportArchive calls list_fpl_report_archive.
func (c *Client) ListReportArchive(req *model.ReportArchiveRequest) (*model.ListReportArchiveResponse, error) {
	svc := fluencyAPI.NewFPLService(c.ingextClient)
	resp, err := svc.ListReportArchive(req)
	if err != nil {
		c.Logger.Error("failed to list report archive", "name", req.Name, "error", err)
		return nil, fmt.Errorf("failed to list report archive: %w", err)
	}
	return resp, nil
}

// GetReportArchive calls get_fpl_report_archive.
func (c *Client) GetReportArchive(req *model.ReportArchiveRequest) (*model.ReportArchiveEntry, error) {
	svc := fluencyAPI.NewFPLService(c.ingextClient)
	entry, err := svc.GetReportArchive(req)
	if err != nil {
		c.Logger.Error("failed to get report archive day", "name", req.Name, "day", req.Day, "id", req.ID, "error", err)
		return nil, fmt.Errorf("failed to get report archive day: %w", err)
	}
	return entry, nil
}

// DownloadReportArchive calls download_fpl_report_archive.
func (c *Client) DownloadReportArchive(req *model.ReportArchiveRequest) (*model.DownloadReportArchiveResponse, error) {
	svc := fluencyAPI.NewFPLService(c.ingextClient)
	resp, err := svc.DownloadReportArchive(req)
	if err != nil {
		c.Logger.Error("failed to download report archive", "name", req.Name, "error", err)
		return nil, fmt.Errorf("failed to download report archive: %w", err)
	}
	return resp, nil
}

// BackfillReportArchive calls backfill_fpl_report_archive.
func (c *Client) BackfillReportArchive(req *model.ReportArchiveRequest) (*model.BackfillReportArchiveResponse, error) {
	svc := fluencyAPI.NewFPLService(c.ingextClient)
	resp, err := svc.BackfillReportArchive(req)
	if err != nil {
		c.Logger.Error("failed to backfill report archive", "name", req.Name, "error", err)
		return nil, fmt.Errorf("failed to backfill report archive: %w", err)
	}
	return resp, nil
}

// ReindexReportArchive calls reindex_fpl_report_archive.
func (c *Client) ReindexReportArchive(req *model.ReportArchiveRequest) (*model.ReindexReportArchiveResponse, error) {
	svc := fluencyAPI.NewFPLService(c.ingextClient)
	resp, err := svc.ReindexReportArchive(req)
	if err != nil {
		c.Logger.Error("failed to reindex report archive", "name", req.Name, "error", err)
		return nil, fmt.Errorf("failed to reindex report archive: %w", err)
	}
	return resp, nil
}
