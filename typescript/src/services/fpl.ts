import type { IngextClient } from "../client.js";
import type {
  BackfillReportArchiveResponse,
  DownloadReportArchiveResponse,
  GetMetricFPLResultResponse,
  GetTaskResponse,
  ListReportArchiveResponse,
  ReindexReportArchiveResponse,
  ReportArchiveEntry,
  ReportArchiveRequest,
  RunFPLV2Report,
} from "../types/fpl.js";

const DS = "api/ds";

export class FPLService {
  constructor(private client: IngextClient) {}

  /** Submit an FPL v2 report. Returns the task ID. */
  async runReport(req: RunFPLV2Report): Promise<number> {
    const res = await this.client.call<{ taskID: number }>(DS, "run_fplv2_report", req);
    return res.taskID;
  }

  /** Fetch task status by ID. */
  async getTaskByID(id: number): Promise<GetTaskResponse> {
    return await this.client.call<GetTaskResponse>(DS, "get_fplv2_task", {
      id,
      fpl: true,
    });
  }

  /** Fetch task results by ID. */
  async getResultsByID(id: number): Promise<GetMetricFPLResultResponse> {
    return await this.client.call<GetMetricFPLResultResponse>(DS, "get_fplv2_result", {
      id,
      fpl: true,
    });
  }

  /**
   * List the scheduled-report archive. Without `name`, the archived reports;
   * with it, that report's days in from..to (default the last 31) with
   * coverage and any running backfill.
   */
  async listReportArchive(req: ReportArchiveRequest = {}): Promise<ListReportArchiveResponse> {
    return await this.client.call<ListReportArchiveResponse>(DS, "list_fpl_report_archive", req);
  }

  /** One archived day, by `id` or `name` + `day`. */
  async getReportArchive(req: ReportArchiveRequest): Promise<ReportArchiveEntry> {
    const res = await this.client.call<{ entry: ReportArchiveEntry }>(DS, "get_fpl_report_archive", req);
    return res.entry;
  }

  /**
   * Archived results: one day (`id`, or `name` + `day`) or a range (`name` +
   * `from` [+ `to`], at most 92 days). `objects` and `limit` trim the result.
   */
  async downloadReportArchive(req: ReportArchiveRequest): Promise<DownloadReportArchiveResponse> {
    return await this.client.call<DownloadReportArchiveResponse>(DS, "download_fpl_report_archive", req);
  }

  /**
   * Run the report for past days (`day`, or `from`..`to`, at most 92, before
   * today) in the background. Completed days are skipped unless `force`.
   */
  async backfillReportArchive(req: ReportArchiveRequest): Promise<BackfillReportArchiveResponse> {
    return await this.client.call<BackfillReportArchiveResponse>(DS, "backfill_fpl_report_archive", req);
  }

  /** Rebuild the archive index from the stored files, for `name` or all reports. */
  async reindexReportArchive(req: ReportArchiveRequest = {}): Promise<ReindexReportArchiveResponse> {
    return await this.client.call<ReindexReportArchiveResponse>(DS, "reindex_fpl_report_archive", req);
  }
}
