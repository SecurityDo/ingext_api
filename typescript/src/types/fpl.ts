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
  type: number;
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

// ---- scheduled-report archive ----
// A daily scheduled report with archive enabled keeps each day's result for 12
// months. Days are "YYYY-MM-DD" in the archive timezone; ranges include both ends.

export interface ReportArchive {
  enabled: boolean;
  /** Day boundary (IANA name); empty means UTC. */
  timezone?: string;
  argFrom?: string;
  argTo?: string;
}

/** kargs of every *_fpl_report_archive function; each reads what it needs. */
export interface ReportArchiveRequest {
  id?: number;
  name?: string;
  day?: string;
  from?: string;
  to?: string;
  state?: "completed" | "aborted" | "running";
  objects?: string[];
  limit?: number;
  force?: boolean;
}

export interface ReportArchiveEntry {
  id?: number;
  account: string;
  reportName: string;
  taskID: number;
  state: string;
  errorMsg?: string;
  rangeFrom: string;
  rangeTo: string;
  createdOn: string;
  completedOn?: string;
  took: number;
  cost: number;
  /** Output rows. */
  total: number;
  resultBytes: number;
  resultPath?: string;
  env?: {
    timezone: string;
    scriptHash?: string;
    processorName?: string;
    arguments?: ProgramArgument[];
    reportConfig?: unknown;
  };
}

export interface ReportArchiveCoverage {
  firstDay: string;
  lastDay: string;
  expected: number;
  completed: number;
  aborted: number;
  running: number;
  missing: number;
  missingDays?: string[];
  abortedDays?: string[];
  runningDays?: string[];
  scriptHashes?: string[];
}

export interface ArchivedReport {
  name: string;
  enabled: boolean;
  timezone: string;
  firstDay?: string;
  lastDay?: string;
  days: number;
  completed: number;
}

export interface ReportArchiveBackfillStatus {
  name: string;
  days: string[];
  done: number;
  failed: number;
  current?: string;
  startedBy?: string;
  startedOn: string;
}

/** Without a name only `reports` is set; with one, the rest. */
export interface ListReportArchiveResponse {
  reports?: ArchivedReport[];
  name?: string;
  enabled?: boolean;
  timezone?: string;
  coverage?: ReportArchiveCoverage;
  entries?: ReportArchiveEntry[];
  backfill?: ReportArchiveBackfillStatus | null;
}

export interface ReportArchiveDay {
  entry: ReportArchiveEntry;
  result: MetricFPLResponse;
}

export interface DownloadReportArchiveResponse {
  name: string;
  timezone: string;
  coverage?: ReportArchiveCoverage;
  days: ReportArchiveDay[];
}

export interface BackfillReportArchiveResponse {
  name: string;
  timezone: string;
  queued: string[];
  skipped: string[];
}

export interface ReindexReportArchiveResponse {
  name: string;
  result: { restored: number; existing: number; failed: number; errors?: string[] };
}
