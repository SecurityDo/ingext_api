export interface RunFPLV2Report {
  entry?: ReportTask;
  reportName?: string;
  arguments?: ProgramArgument[];
}

export interface GetTaskResponse {
  entry?: ReportTask;
  fpl?: string;
}

export interface GetMetricFPLResultResponse {
  task?: ReportTask;
  fpl?: string;
  result?: MetricFPLResponse;
}

export interface ReportTask {
  id?: number;
  name?: string;
  description?: string;
  state?: string;
  error?: string;
  kargs?: unknown;
  fplVersion: number;
  fpl?: string;
  took: number;
  cost?: number;
  arguments?: ProgramArgument[];
}

export interface ProgramArgument {
  name?: string;
  description?: string;
  value?: string;
  defaultValue?: string;
  optional?: boolean;
  type?: string;
  isList?: boolean;
  enumValues?: string[];
  enum?: boolean;
}

export interface MetricFPLResponse {
  taskID?: number;
  objects: FplObject[];
  error: string;
  env?: Record<string, unknown>;
  cost?: number;
  took?: number;
}

export interface FplObject {
  name?: string;
  description?: string;
  table?: FplTable;
  stream?: FplStream;
  alert?: FplAlertResult;
  list?: PrimitiveValue[];
}

export interface FplTable {
  rowCount?: number;
  columns?: FplColumnInfo[];
  rows?: Record<string, unknown>[];
  fplRows?: FplRowState[];
  resource?: string;
}

export interface FplColumnInfo {
  name: string;
  unit?: string;
  columnType?: string;
  aggType?: string;
  isVariable: boolean;
  isHidden?: boolean;
  dynamic?: boolean;
}

export interface FplValueState {
  v?: PrimitiveValue;
}

export interface FplRowState {
  key?: PrimitiveValue;
  values?: FplValueState[];
}

export interface PrimitiveValue {
  type: string;
  bool?: boolean;
  int?: number;
  float?: number;
  string?: string;
  list?: PrimitiveValue[];
  map?: Record<string, PrimitiveValue>;
  jsonObj?: Record<string, unknown>;
  jsonArray?: unknown[];
  /** Base64-encoded bytes (Go marshals []byte as base64). */
  blob?: string;
}

export interface FplStream {
  from?: number;
  to?: number;
  interval?: number;
  slots?: number[];
  dimensions?: string[];
  metrics?: FplStreamMetric[];
  isBool?: boolean;
  dimensionName?: string;
  name?: string;
  description?: string;
  unit?: string;
  metric?: string;
  namespace?: string;
}

export interface FplStreamMetric {
  flags?: number[];
  values?: number[];
  bools?: boolean[];
  key?: string;
  id?: number;
  tags?: TableColumnTag[];
}

export interface TableColumnTag {
  key?: string;
  values?: string[];
}

export interface FplAlertResult {
  alerts?: FplAlert[];
  slots?: number[];
  from?: number;
  to?: number;
  interval?: number;
  sliderStart?: number;
}

export interface FplAlert {
  name?: string;
  description?: string;
  isAlert?: boolean;
  eval?: AnomalyEvalT;
  key?: string;
  account?: string;
  accountID?: string;
  region?: string;
  dimensionTags?: StreamDimensionTag[];
  strategy?: string;
  timestamp?: number;
  alertMetric?: FplStreamMetric;
  values?: StreamHistogram[];
  associates?: StreamHistogram[];
}

export interface StreamDimensionTag {
  dimension?: string;
  key?: string;
  tags?: TableColumnTag[];
}

export interface StreamHistogram {
  name?: string;
  unit?: string;
  metric?: FplStreamMetric;
}

/**
 * The Go AnomalyEvalT has only unexported fields, so it serializes as `{}`.
 * Kept as an opaque shape for type completeness.
 */
export type AnomalyEvalT = Record<string, never>;
