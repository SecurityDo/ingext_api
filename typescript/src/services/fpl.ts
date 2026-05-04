import type { IngextClient } from "../client.js";
import type {
  GetMetricFPLResultResponse,
  GetTaskResponse,
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
}
