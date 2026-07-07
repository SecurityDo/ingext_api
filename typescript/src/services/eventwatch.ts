import type { IngextClient } from "../client.js";
import type { ElasticSearchResult } from "../types/elasticsearch.js";
import type { GenericDAORequest, SimpleSearchOption } from "../types/dao.js";
import type { EventWatchBucket } from "../types/eventwatch.js";

const DS = "api/ds";
const RULE_DAO = "eventwatch_bucket_dao";

interface ElasticSearchRequest {
  options: SimpleSearchOption;
}

/** Response from the eventwatch_bucket_dao "get" action. */
interface EventWatchRuleGetResponse {
  entry: EventWatchBucket | null;
}

/** Response from the eventwatch_bucket_dao "list" action. */
export interface EventWatchRuleListResponse {
  entries: EventWatchBucket[];
  tags: string[];
  groups: string[];
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

  /** List all eventwatch rules along with the distinct tags and groups. */
  async listRule(): Promise<EventWatchRuleListResponse> {
    const req: GenericDAORequest<EventWatchBucket> = { action: "list" };
    return await this.client.call<EventWatchRuleListResponse>(DS, RULE_DAO, req);
  }

  /** Get a single eventwatch rule by name; empty name returns the first rule. */
  async getRule(name = ""): Promise<EventWatchBucket | null> {
    const req: GenericDAORequest<EventWatchBucket> = {
      action: "get",
      args: { id: name },
    };
    const res = await this.client.call<EventWatchRuleGetResponse>(DS, RULE_DAO, req);
    return res.entry ?? null;
  }

  /** Create a new eventwatch rule. */
  async addRule(entry: EventWatchBucket): Promise<void> {
    const req: GenericDAORequest<EventWatchBucket> = {
      action: "add",
      args: { entry },
    };
    await this.client.call(DS, RULE_DAO, req);
  }

  /** Update an existing eventwatch rule. */
  async updateRule(entry: EventWatchBucket): Promise<void> {
    const req: GenericDAORequest<EventWatchBucket> = {
      action: "update",
      args: { entry },
    };
    await this.client.call(DS, RULE_DAO, req);
  }

  /** Flip the disabled state of an eventwatch rule identified by name. */
  async toggleRule(name: string): Promise<void> {
    const req: GenericDAORequest<EventWatchBucket> = {
      action: "toggle",
      args: { id: name },
    };
    await this.client.call(DS, RULE_DAO, req);
  }

  /** Remove an eventwatch rule identified by name. */
  async deleteRule(name: string): Promise<void> {
    const req: GenericDAORequest<EventWatchBucket> = {
      action: "delete",
      args: { id: name },
    };
    await this.client.call(DS, RULE_DAO, req);
  }
}
