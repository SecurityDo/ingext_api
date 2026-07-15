import type { IngextClient } from "../client.js";
import type { CollectorForWeb, CollectorT } from "../types/collector.js";

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

  /**
   * Recreate an existing collector on this site with its token preserved, which is
   * how a collector is migrated from one site to another without re-enrolling the
   * agent. The collector name must not already exist on the target site.
   */
  async collectorImport(collector: CollectorT): Promise<void> {
    await this.client.call<unknown>(DS, "collector_import", { collector });
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
