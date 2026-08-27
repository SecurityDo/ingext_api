import type { IngextClient } from "../client.js";
import type { ElasticSearchResult } from "../types/elasticsearch.js";
import type {
  GenericDAORequest,
  GenericDAORequestArgs,
  SimpleSearchOption,
} from "../types/dao.js";
import type {
  BehaviorEvent,
  BehaviorEventFilterT,
  EventWatchBucket,
} from "../types/eventwatch.js";

const DS = "api/ds";
const RULE_DAO = "eventwatch_bucket_dao";
const FILTER_DAO = "behavior_filter_dao";

interface ElasticSearchRequest {
  options: SimpleSearchOption;
}

/** Response from the eventwatch_bucket_dao "get" action. */
interface EventWatchRuleGetResponse {
  entry: EventWatchBucket | null;
}

/**
 * Runs one rule against one event, without deploying the rule or storing what
 * it produces.
 *
 * `bucket` is the whole rule, the shape the DAO stores, and its `name` must be
 * set: the endpoint panics on a rule with no name ("runtime error: invalid
 * memory address or nil pointer dereference") rather than reporting it. `input`
 * is the event as a JSON object -- a JSON string holding an event is accepted
 * on the wire and then quietly ignored.
 */
export interface EventWatchRuleTestRequest {
  bucket: EventWatchBucket;
  input?: unknown;
}

/**
 * Whether the rule selected the event, and what it produced from it.
 *
 * The event is echoed back: under `input` normally, and under `output` on a hit
 * that produced a behavior event, where it has been null in every response seen
 * so far. `hit` only means the event selector matched -- a rule with no fields
 * or behavior rule configured hits without producing a signal or a behavior
 * event. `signals` are JSON documents carried as strings; decode them with
 * {@link decodeEventWatchSignals}.
 */
export interface EventWatchRuleTestResult {
  hit: boolean;
  input?: unknown;
  output?: unknown;
  signals: string[] | null;
  behaviorEvent: BehaviorEvent | null;
}

/**
 * One entry of `EventWatchRuleTestResult.signals`, unpacked. `ts` is in
 * seconds, unlike the behavior event's millisecond `timestamp`.
 */
export interface EventWatchTestSignal {
  signal: string;
  ts: number;
  count: number;
  key: string;
  valueMap: Record<string, string>;
}

/** Unpack the JSON documents the endpoint returns as strings. */
export function decodeEventWatchSignals(
  result: EventWatchRuleTestResult,
): EventWatchTestSignal[] {
  return (result.signals ?? []).map((raw, i) => {
    try {
      return JSON.parse(raw) as EventWatchTestSignal;
    } catch (err) {
      throw new Error(`signal ${i} is not valid JSON: ${String(err)}`);
    }
  });
}

/**
 * Build the DAO id of a behavior filter. A filter is keyed by its owning behavior
 * rule plus its name, and the two travel to the server joined by a slash in the
 * single `id` arg — the server splits on the first slash, so a name may contain
 * slashes but a behavior rule may not.
 */
export function behaviorFilterID(behaviorRule: string, name: string): string {
  return `${behaviorRule}/${name}`;
}

/** Build the args for the id-addressed actions, rejecting ids the server cannot resolve. */
function behaviorFilterArgs(
  behaviorRule: string,
  name: string,
): GenericDAORequestArgs<BehaviorEventFilterT> {
  if (!behaviorRule) {
    throw new Error("behavior filter behaviorRule is required");
  }
  if (!name) {
    throw new Error("behavior filter name is required");
  }
  if (behaviorRule.includes("/")) {
    throw new Error(`behavior filter behaviorRule "${behaviorRule}" must not contain a slash`);
  }
  return { id: behaviorFilterID(behaviorRule, name) };
}

/** Response from the behavior_filter_dao "get" action. */
interface BehaviorFilterGetResponse {
  entry: BehaviorEventFilterT | null;
}

/** Response from the behavior_filter_dao "list" action. */
interface BehaviorFilterListResponse {
  entries: BehaviorEventFilterT[];
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

  /**
   * Run a rule definition that need not be deployed against one event. The
   * rule's `name` must be set and the event must be an object -- see
   * {@link EventWatchRuleTestRequest}.
   */
  async testRule(
    rule: EventWatchBucket,
    event?: unknown,
  ): Promise<EventWatchRuleTestResult> {
    if (!rule?.name?.trim()) {
      throw new Error(
        "rule name is required: the endpoint panics on a rule with no name",
      );
    }
    if (event !== undefined) {
      if (event === null || typeof event !== "object" || Array.isArray(event)) {
        throw new Error(
          "event must be a JSON object: anything else is accepted and then ignored",
        );
      }
    }
    const req: EventWatchRuleTestRequest = { bucket: rule, input: event };
    return await this.client.call<EventWatchRuleTestResult>(
      DS,
      "eventwatch_rule_test",
      req,
    );
  }

  /**
   * Read the stored rule of that name and run it against one event. A name that
   * is not deployed throws from the read ("bucket not found"), before anything
   * is tested.
   */
  async testDeployedRule(
    name: string,
    event?: unknown,
  ): Promise<EventWatchRuleTestResult> {
    if (!name.trim()) {
      throw new Error(
        "rule name is required: an empty one returns whichever rule the DAO lists first",
      );
    }
    const rule = await this.getRule(name);
    if (!rule) {
      throw new Error(`rule ${name} not found`);
    }
    return await this.testRule(rule, event);
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

  /** List every behavior filter of the account, across all behavior rules. */
  async listBehaviorFilter(): Promise<BehaviorEventFilterT[]> {
    const req: GenericDAORequest<BehaviorEventFilterT> = { action: "list" };
    const res = await this.client.call<BehaviorFilterListResponse>(DS, FILTER_DAO, req);
    return res.entries ?? [];
  }

  /**
   * Get a single behavior filter. Both the owning behavior rule and the filter
   * name are required; the server reports a missing filter as an error rather
   * than an empty result.
   */
  async getBehaviorFilter(
    behaviorRule: string,
    name: string,
  ): Promise<BehaviorEventFilterT | null> {
    const req: GenericDAORequest<BehaviorEventFilterT> = {
      action: "get",
      args: behaviorFilterArgs(behaviorRule, name),
    };
    const res = await this.client.call<BehaviorFilterGetResponse>(DS, FILTER_DAO, req);
    return res.entry ?? null;
  }

  /**
   * Create a new behavior filter. The entry needs `name`, `behaviorRule`, and at
   * least one entry in `filters`; the server stamps `createdOn`/`updatedOn` and
   * rejects a name already used under the same behavior rule.
   */
  async addBehaviorFilter(entry: BehaviorEventFilterT): Promise<void> {
    const req: GenericDAORequest<BehaviorEventFilterT> = { action: "add", args: { entry } };
    await this.client.call(DS, FILTER_DAO, req);
  }

  /**
   * Replace an existing behavior filter, identified by the entry's
   * `behaviorRule` and `name` rather than by an id. Changing either field
   * addresses a different filter, so the update fails as "behavior filter does
   * not exist" instead of renaming: add the new filter and delete the old one.
   */
  async updateBehaviorFilter(entry: BehaviorEventFilterT): Promise<void> {
    const req: GenericDAORequest<BehaviorEventFilterT> = { action: "update", args: { entry } };
    await this.client.call(DS, FILTER_DAO, req);
  }

  /** Flip the disabled state of a behavior filter. */
  async toggleBehaviorFilter(behaviorRule: string, name: string): Promise<void> {
    const req: GenericDAORequest<BehaviorEventFilterT> = {
      action: "toggle",
      args: behaviorFilterArgs(behaviorRule, name),
    };
    await this.client.call(DS, FILTER_DAO, req);
  }

  /** Remove a behavior filter. */
  async deleteBehaviorFilter(behaviorRule: string, name: string): Promise<void> {
    const req: GenericDAORequest<BehaviorEventFilterT> = {
      action: "delete",
      args: behaviorFilterArgs(behaviorRule, name),
    };
    await this.client.call(DS, FILTER_DAO, req);
  }
}
