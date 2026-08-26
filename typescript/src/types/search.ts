export interface LakeFacetSearchOption {
  searchStr?: string;
  range_from: number;
  range_to: number;
  fetchLimit?: number;
  fetchOffset: number;
  sortField?: string;
  sortOrder?: string;
  facets?: FacetsOption;
}

export interface FacetsOption {
  dateFacets: DateFacetEntry[];
  facets: FacetEntry[];
  mustFilters: FilterEntry[];
  mustNotFilters: FilterEntry[];
}

export interface DateFacetEntry {
  name?: string;
  key?: string;
  value?: string;
}

export interface FacetEntry {
  title: string;
  field: string;
  order: string;
  size: number;
}

export interface FilterEntry {
  field: string;
  terms: string[];
}

export interface LakeSearchResponse {
  took: number;
  hits?: LegacySearchHits;
  aggregations?: Record<string, SearchAggregate>;
  query: unknown;
  terms?: string[];
  taskID?: string;
  range_from?: number;
  range_to?: number;
  cost?: number;
  total?: number;
  filtered?: number;
}

export interface LegacySearchHits {
  total: number;
  sortFieldType: string;
  hits: LegacySearchHit[];
  ascend: boolean;
  limit: number;
}

export interface LegacySearchHit {
  _source: unknown;
  _id?: string;
}

export interface SearchAggregate {
  /**
   * Raw bucket objects. The two shapes are not interchangeable: a term facet
   * keys on a string, the `dateHistogram` aggregation keys on an epoch
   * millisecond. Narrow with `FacetBucket` / `HistogramBucket`.
   */
  buckets: unknown[];
}

/**
 * Aggregation name the platform reserves for the date histogram over the search
 * range. Every other key of `LakeSearchResponse.aggregations` is a term facet,
 * named after the facet field.
 */
export const DATE_HISTOGRAM_AGGREGATION = "dateHistogram";

/** One term of a facet aggregation. */
export interface FacetBucket {
  key: string;
  doc_count: number;
}

/**
 * One slot of the date histogram. `key` is the epoch millisecond the slot
 * starts at; the platform returns a fixed number of slots spanning the search
 * range, empty ones included.
 */
export interface HistogramBucket {
  key: number;
  doc_count: number;
  subAggs?: { name: string; value: number }[];
}

/**
 * kargs payload for the lake_search endpoint.
 *
 * The endpoint also declares `dataType`, `partition` and `dayIndex`, but never
 * reads them, so they are left off the wire.
 */
export interface LakeSearchRequest {
  index?: string;
  options: LakeFacetSearchOption;
}

/** KQL search request payload for the kql_search endpoint. */
export interface KQLSearchRequest {
  index?: string;
  rangeFrom: number;
  rangeTo: number;
  kql?: string;
}

/** KQL validate response (mirrors lakeAPI.KQLValidateResponse on the Go side). */
export interface KQLValidateResponse {
  ok: boolean;
  error?: string;
  table?: string;
  timeRangeFound?: boolean;
}
