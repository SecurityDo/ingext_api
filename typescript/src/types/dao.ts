/** Generic DAO request envelope used by many backend endpoints. */
export interface GenericDAORequestArgs<T> {
  id?: string;
  entry?: T;
  flag?: boolean;
}

export interface GenericDAORequest<T> {
  action: string;
  args?: GenericDAORequestArgs<T>;
}

export interface GenericDaoAddResponse {
  id: string;
}

export interface GenericDaoListResponse<T> {
  entries: T[];
}

/** Filter and facet shapes used by `SimpleSearchOption` (eventwatch et al). */
export interface SimpleFacetEntry {
  field: string;
  order: string;
  size: number;
}

export interface SimpleFilterEntry {
  field: string;
  /** Terms can be strings, numbers, or other JSON-compatible scalars. */
  terms: unknown[];
}

export interface SimpleFacetsOption {
  facets?: SimpleFacetEntry[];
  mustFilters?: SimpleFilterEntry[];
  mustNotFilters?: SimpleFilterEntry[];
}

export interface SimpleSearchOption {
  searchStr?: string;
  fetchOffset: number;
  fetchLimit: number;
  sortField?: string;
  sortOrder?: string;
  facets?: SimpleFacetsOption;
  range_from: number;
  range_to: number;
  range_field: string;
}
