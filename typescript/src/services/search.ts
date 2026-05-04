import type { IngextClient } from "../client.js";
import type { KQLSearchResponse } from "../types/kql.js";
import type {
  KQLSearchRequest,
  KQLValidateResponse,
} from "../types/search.js";

const DS = "api/ds";

export class SearchService {
  constructor(private client: IngextClient) {}

  /** Run a KQL query. */
  async kqlSearch(kql: string): Promise<KQLSearchResponse> {
    const req: KQLSearchRequest = { kql, rangeFrom: 0, rangeTo: 0 };
    return await this.client.call<KQLSearchResponse>(DS, "kql_search", req);
  }

  /** Parse-only: validate a KQL query without executing it. */
  async kqlValidate(kql: string): Promise<KQLValidateResponse> {
    return await this.client.call<KQLValidateResponse>(DS, "kql_validate", { kql });
  }
}
