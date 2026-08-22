import type { IngextClient } from "../client.js";
import type {
  LakeFacetSearchOption,
  LakeSearchResponse,
} from "../types/search.js";

const DS = "api/ds";

interface ResourceSearchRequest {
  options: LakeFacetSearchOption;
  customer: string;
  resource: string;
}

interface ResourceDumpDeleteRequest {
  customer: string;
}

interface ResourceDumpDeleteResponse {
  deleted: number;
}

export class ResourceService {
  constructor(private client: IngextClient) {}

  async search(resourceType: string, customer: string): Promise<LakeSearchResponse> {
    const req: ResourceSearchRequest = {
      options: {
        range_from: 0,
        range_to: 0,
        fetchOffset: 0,
        fetchLimit: 1000,
        facets: {
          dateFacets: [],
          facets: [
            { title: "Groups", field: "@office365User.groups", order: "", size: 20 },
          ],
          mustFilters: [],
          mustNotFilters: [],
        },
      },
      customer,
      resource: resourceType,
    };
    return await this.client.call<LakeSearchResponse>(DS, "resource_search", req);
  }

  /**
   * Purge every resource dump belonging to one customer, across all resource
   * types, and return the number of dump folders removed.
   *
   * This is the same cleanup that deleting the owning plugin data source
   * performs; the endpoint exists for the dumps of data sources that were
   * already removed while that purge silently matched nothing. It destroys
   * collected data with no undo.
   *
   * The customer is matched literally - `"_all_"` is a customer name here, not
   * a wildcard. An unknown customer is not an error: it deletes nothing and
   * returns 0, which is the only way to tell a purge from a no-op.
   */
  async dumpDelete(customer: string): Promise<number> {
    if (!customer) {
      throw new Error("resource dump customer is required");
    }
    const req: ResourceDumpDeleteRequest = { customer };
    const resp = await this.client.call<ResourceDumpDeleteResponse>(
      DS,
      "ingext_resource_dump_delete",
      req,
    );
    return resp.deleted ?? 0;
  }
}
