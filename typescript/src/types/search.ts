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
  buckets: unknown[];
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
