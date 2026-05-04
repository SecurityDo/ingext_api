import type { LegacySearchHit, SearchAggregate } from "./search.js";

export interface ElasticSearchResult {
  took?: number;
  terminated_early?: boolean;
  num_reduce_phases?: number;
  _clusters?: SearchResultCluster;
  _scroll_id?: string;
  hits?: SearchHits;
  aggregations?: Record<string, SearchAggregate>;
  timed_out?: boolean;
  status?: number;
  pit_id?: string;
}

export interface SearchResultCluster {
  successful?: number;
  total?: number;
  skipped?: number;
}

export interface SearchHits {
  max_score?: number;
  hits?: LegacySearchHit[];
}

export interface TotalHits {
  value: number;
  relation: string;
}
