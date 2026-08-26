import type { IngextClient } from "../client.js";
import type { KQLSearchResponse } from "../types/kql.js";
import type {
  FacetEntry,
  FilterEntry,
  KQLSearchRequest,
  KQLValidateResponse,
  LakeSearchRequest,
  LakeSearchResponse,
} from "../types/search.js";

const DS = "api/ds";

/** The platform's hard cap on fetchOffset+fetchLimit for one lake search. */
export const LAKE_SEARCH_MAX_PAGINATION = 5000;

/**
 * A facet of size 0 is not "unlimited" on the platform - it returns no buckets
 * at all, which looks like a broken search rather than a missing default.
 */
const LAKE_SEARCH_DEFAULT_FACET_SIZE = 20;

/**
 * A fetchLimit of 0 is not a facet-only search: it panics the search node
 * ("index out of range [0] with length 0") because the top-hits heap is built
 * with a capacity of 0 and then indexed. There is no way to ask this endpoint
 * for counts without hits.
 */
const LAKE_SEARCH_DEFAULT_FETCH_LIMIT = 100;

/**
 * The platform keys the date histogram aggregation by the name of the first
 * date facet in the request. Without one it lands under the empty string.
 */
const LAKE_SEARCH_DATE_HISTOGRAM = "dateHistogram";

/**
 * lake_search has no server-side default for the sort field, and an empty one
 * fails its query parser.
 */
const LAKE_SEARCH_DEFAULT_SORT_FIELD = "@timestamp";
const LAKE_SEARCH_DEFAULT_SORT_ORDER = "desc";

/**
 * One lake facet search: a Lucene query and a time range, narrowed by term
 * filters, with term facets counted over whatever matches.
 */
export interface LakeSearchOptions {
  /**
   * The datalake index to search, in `<datalake>-<index>` form:
   * `managed-Office365` is index `Office365` of the `managed` datalake. The
   * bare index name is what list_data_tables reports as a stream table. An
   * omitted index searches the `default` index.
   */
  index?: string;

  /** Lucene query string. Empty matches every event in range. */
  query?: string;

  /**
   * Range bounds in epoch milliseconds. Both are required: without a range the
   * platform scans the whole index.
   */
  rangeFrom: number;
  rangeTo: number;

  /**
   * Fields to count terms on. Each needs a `field`; `size` defaults to 20 and
   * `title`, which the platform ignores, defaults to the field name. The
   * response is keyed by field. Do not assume it holds exactly the requested
   * set: a console response for this endpoint carried facets no one asked for,
   * so read the aggregations by name rather than by position.
   */
  facets?: Partial<FacetEntry>[];

  /**
   * Exact-term filters. Terms within one entry are OR'ed, and the entries are
   * AND'ed.
   */
  mustFilters?: FilterEntry[];
  mustNotFilters?: FilterEntry[];

  /**
   * How many hits to return; omitted or 0 means the default of 100. It cannot
   * be used to ask for counts alone - a fetchLimit of 0 reaches the platform
   * and panics the search node. fetchOffset+fetchLimit may not exceed
   * LAKE_SEARCH_MAX_PAGINATION.
   */
  fetchLimit?: number;
  fetchOffset?: number;

  /** Hit ordering; defaults to `@timestamp` / `desc`. */
  sortField?: string;
  sortOrder?: string;
}

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

  /**
   * Run a facet search against one datalake index and return the matching hits
   * together with the term counts.
   *
   * This is the search behind the console's event explorer: a Lucene query over
   * a time range, must / must-not term filters, and a facet per field to count.
   * The response also carries a `dateHistogram` aggregation - a fixed set of
   * slots spanning the range, empty ones included. That aggregation is keyed by
   * the name of the date facet the request asked for, so the request always
   * names one: without it the histogram comes back under the empty string.
   *
   * `total` counts the events that matched, `filtered` the ones that were read
   * and discarded.
   */
  async lakeSearch(options: LakeSearchOptions): Promise<LakeSearchResponse> {
    const req = buildLakeSearchRequest(options);
    return await this.client.call<LakeSearchResponse>(DS, "lake_search", req);
  }
}

/**
 * Validates the options and fills in the defaults the endpoint has none for.
 * The checks are client-side on purpose: each of these fails deep inside the
 * search service, where the error text says nothing about which field was wrong.
 */
export function buildLakeSearchRequest(options: LakeSearchOptions): LakeSearchRequest {
  if (!options) {
    throw new Error("lake search options are required");
  }
  const { rangeFrom, rangeTo } = options;
  if (!(rangeFrom > 0) || !(rangeTo > 0)) {
    throw new Error("lake search requires rangeFrom and rangeTo in epoch milliseconds");
  }
  if (rangeFrom >= rangeTo) {
    throw new Error(
      `lake search range is empty: rangeFrom ${rangeFrom} is not before rangeTo ${rangeTo}`,
    );
  }
  const requestedLimit = options.fetchLimit ?? 0;
  const fetchOffset = options.fetchOffset ?? 0;
  if (requestedLimit < 0 || fetchOffset < 0) {
    throw new Error("lake search fetchLimit and fetchOffset cannot be negative");
  }
  // A fetchLimit of 0 panics the search node rather than returning counts
  // alone, so an unset limit becomes the default instead of being sent.
  const fetchLimit = requestedLimit === 0 ? LAKE_SEARCH_DEFAULT_FETCH_LIMIT : requestedLimit;
  if (fetchOffset + fetchLimit > LAKE_SEARCH_MAX_PAGINATION) {
    throw new Error(
      `lake search pagination limit is ${LAKE_SEARCH_MAX_PAGINATION} entries: ` +
        `fetchOffset ${fetchOffset} + fetchLimit ${fetchLimit}`,
    );
  }

  const facets: FacetEntry[] = (options.facets ?? []).map((facet, i) => {
    if (!facet?.field) {
      throw new Error(`lake search facet ${i}: field is required`);
    }
    return {
      title: facet.title || facet.field,
      field: facet.field,
      order: facet.order ?? "",
      size: facet.size && facet.size > 0 ? facet.size : LAKE_SEARCH_DEFAULT_FACET_SIZE,
    };
  });

  return {
    index: options.index,
    options: {
      searchStr: options.query,
      range_from: rangeFrom,
      range_to: rangeTo,
      fetchLimit,
      fetchOffset,
      sortField: options.sortField || LAKE_SEARCH_DEFAULT_SORT_FIELD,
      sortOrder: options.sortOrder || LAKE_SEARCH_DEFAULT_SORT_ORDER,
      // facets, and the arrays inside it, are always sent: the endpoint walks
      // them without a nil check. The date facet is named so the histogram
      // lands under DATE_HISTOGRAM_AGGREGATION rather than under "".
      facets: {
        dateFacets: [{ name: LAKE_SEARCH_DATE_HISTOGRAM }],
        facets,
        mustFilters: checkFilters("mustFilter", options.mustFilters),
        mustNotFilters: checkFilters("mustNotFilter", options.mustNotFilters),
      },
    },
  };
}

/**
 * Rejects the two filters that silently do nothing useful: one with no field,
 * and one with no terms to match.
 */
function checkFilters(kind: string, filters?: FilterEntry[]): FilterEntry[] {
  return (filters ?? []).map((filter, i) => {
    if (!filter?.field) {
      throw new Error(`lake search ${kind} ${i}: field is required`);
    }
    if (!filter.terms?.length) {
      throw new Error(`lake search ${kind} ${i} (${filter.field}): at least one term is required`);
    }
    return filter;
  });
}
