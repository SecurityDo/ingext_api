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
}
