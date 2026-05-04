import type { IngextClient } from "../client.js";
import type { CollectorForWeb } from "../types/collector.js";

const DS = "api/ds";

export class CollectorService {
  constructor(private client: IngextClient) {}

  async collectorList(): Promise<CollectorForWeb[]> {
    const res = await this.client.call<{ entries: CollectorForWeb[] }>(
      DS,
      "collector_list",
      {},
    );
    return res.entries ?? [];
  }

  async collectorStatus(
    collector: string,
    cargs?: Record<string, unknown>,
  ): Promise<Record<string, unknown>> {
    return await this.client.call<Record<string, unknown>>(DS, "get_system_status", {
      collector,
      cargs,
    });
  }
}
