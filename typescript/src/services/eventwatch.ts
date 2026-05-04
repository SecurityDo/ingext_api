import type { IngextClient } from "../client.js";
import type { ElasticSearchResult } from "../types/elasticsearch.js";
import type { SimpleSearchOption } from "../types/dao.js";

const DS = "api/ds";

interface ElasticSearchRequest {
  options: SimpleSearchOption;
}

function defaultOptions(
  searchString: string,
  rangeFrom: number,
  rangeTo: number,
  rangeField: string,
  sortField: string,
  sortOrder: string,
): SimpleSearchOption {
  return {
    searchStr: searchString,
    range_from: rangeFrom,
    range_to: rangeTo,
    range_field: rangeField,
    fetchLimit: 100,
    fetchOffset: 0,
    sortField,
    sortOrder,
    facets: {
      facets: [],
      mustFilters: [],
      mustNotFilters: [],
    },
  };
}

export class EventWatchService {
  constructor(private client: IngextClient) {}

  /** Summary search over BehaviorSummary records in the given time range. */
  async summarySearch(
    searchString: string,
    rangeFrom: number,
    rangeTo: number,
  ): Promise<ElasticSearchResult> {
    const req: ElasticSearchRequest = {
      options: defaultOptions(searchString, rangeFrom, rangeTo, "from", "to", "desc"),
    };
    return await this.client.call<ElasticSearchResult>(DS, "behavior_summary_search", req);
  }

  /** Timeline search (FSM behavior search) in the given time range. */
  async timelineSearch(
    searchString: string,
    rangeFrom: number,
    rangeTo: number,
  ): Promise<ElasticSearchResult> {
    const req: ElasticSearchRequest = {
      options: defaultOptions(
        searchString,
        rangeFrom,
        rangeTo,
        "timestamp",
        "timestamp",
        "desc",
      ),
    };
    return await this.client.call<ElasticSearchResult>(DS, "fsm_behavior_search", req);
  }

  /** Rule search (no time range). */
  async ruleSearch(searchString: string): Promise<ElasticSearchResult> {
    const req: ElasticSearchRequest = {
      options: defaultOptions(searchString, 0, 0, "", "name", "asc"),
    };
    return await this.client.call<ElasticSearchResult>(DS, "eventwatch_bucket_search", req);
  }
}
