/**
 * KQL scalar type names that appear in `ColumnDef.DataType` over the wire.
 * Mirrors the strings emitted by `KType.String()` in
 * `kql/model/scalar_type.go`.
 */
export type KType =
  | "null"
  | "bool"
  | "int"
  | "long"
  | "real"
  | "string"
  | "datetime"
  | "timespan"
  | "guid"
  | "decimal"
  | "dynamic"
  | "unknown";

/** A single cell value in a KQL row. */
export type KCell =
  | string
  | number
  | boolean
  | null
  | unknown[]
  | Record<string, unknown>;

export interface ColumnDef {
  ColumnName: string;
  DataType: KType | string;
}

export interface DataTable {
  TableName: string;
  Columns: ColumnDef[];
  /** Each row is an array of cells positionally aligned with `Columns`. */
  Rows: KCell[][];
}

export interface DataSet {
  Tables: DataTable[];
}

export interface KQLSubSearchRequest {
  kql: string;
  strategy: number;
  timeRangeFrom?: number;
  timeRangeTo?: number;
  timeRangeFound?: boolean;
  whereOpIndices?: number[];
}

export interface KQLSubSearchJobResult {
  totalBytes?: number;
  totalRows?: number;
  data?: DataSet;
  groupStates?: unknown;
  facetResults?: unknown;
  strategy: number;
  schemaNames?: string[];
}

export interface KQLSearchResponse {
  total: number;
  totalBytes?: number;
  data?: DataSet;
  range_from?: number;
  range_to?: number;
  cost?: number;
}
