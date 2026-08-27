export interface BehaviorSummary {
  id: string;
  from: number;
  to: number;
  count: number;
  key: string;
  keyType?: string;
  dayIndex?: string;
  interval?: string;
  behaviorRules: string[];
  behaviors: string[];
  risks?: string[];
  riskScore: number;
  mitreTags?: string[];
  summaryList?: BehaviorRuleSummary[];
  scoreLevel: string;
  comments: UserComment[];
  status?: string;
  incident: boolean;
  scoreAdjust: number;
  updatedOn: number;
  vectorData?: number[];
  fingerprintHash?: string;
  fingerprint?: string;
  ruleFP?: string;
}

export interface UserComment {
  content: string;
  actions: string[];
  username: string;
  createdOn: number;
}

export interface BehaviorRuleSummary {
  behaviorRule: string;
  behavior: string;
  from: number;
  to: number;
  count: number;
  attributeSummaries: ValueSummmaryEntry[];
  hits?: BehaviorRuleHit[];
  riskScore: number;
  risks: string[];
}

export interface BehaviorRuleHit {
  name: string;
  description?: string;
  fields?: string[];
  values?: string[];
  scope?: string;
  risks?: string[];
}

export interface ValueSummmaryEntry {
  key: string;
  aliase?: string;
  values: string[];
}

export interface ValueEntry {
  key: string;
  aliase?: string;
  value: string;
}

/**
 * `sequence`, `originalKey`, `originalKeyType`, `description`, `scoreLevel` and
 * `AttributeMap` are what eventwatch_rule_test returns on a hit; a behavior
 * event read back from a search carries the rest. Note `AttributeMap`'s
 * capitalized key, which is the endpoint's.
 */
export interface BehaviorEvent {
  sequence?: number;
  timestamp: number;
  key: string;
  originalKey?: string;
  originalKeyType?: string;
  keyType?: string;
  title?: string;
  description?: string;
  behaviorRule: string;
  behavior: string;
  riskScore: number;
  scoreLevel?: string;
  AttributeMap?: unknown;
  attributes: ValueEntry[];
  risks?: string[];
  ruleRisks?: string[];
  ruleHits?: BehaviorRuleHit[];
}

export interface EventWatchBucket {
  id: number;
  repository: string;
  name?: string;
  group: string;
  description?: string;
  disabled: boolean;
  enableTotal: boolean;
  eventType: string;
  resourceType: string;
  aggregationBucket: boolean;
  isProcessor: boolean;
  processor: string;
  ruleGroup: string;
  bucketType?: string;
  keyType?: string;
  timeSlices?: string[];
  fields?: string[];
  eventSelector: EventMatch;
  aggregations: Aggregation[];
  tagActions: TagAction[];
  signals: EventWatchSignal[];
  behaviorRule: BehaviorRuleT;
  translation: TranslationT;
  updatedOn: string;
  createdOn: string;
  deployedOn: string;
  tags: string[];
  searchProfile: string;
  discard: boolean;
  comment?: string;
  history: UpdateRecord[];
}

export interface EventWatchSignalRef {
  behavior: BehaviorRuleT;
  signal: EventWatchSignal;
  bucket: string;
  searchProfile: string;
  eventType: string;
  groupBy: string;
  signalSearchQuery: string;
  query: string;
  signalName: string;
  fieldEntity?: EventWatchEntityInfo[];
}

export interface Aggregation {
  type: string;
  name: string;
  field?: string;
  groupBy: string;
  groupByEntity: string;
  aggregationBucket: string;
  aggregationKey: string;
  aggregationName: string;
  externalAggregation: boolean;
}

export interface EventMatch {
  matchAll: boolean;
  eventFilter: string;
  query?: unknown;
  mustFilters: BucketFilterEntry[];
  mustNotFilters: BucketFilterEntry[];
  jsonQuery?: unknown;
  jsonFilter?: LambdaFilter;
  sigmaQuery?: unknown;
  lvdbQuery?: unknown;
  lvdbQueryFlag: boolean;
  lvdbQueryTree?: unknown;
}

export interface BucketFilterEntry {
  field: string;
  terms: string[];
  filterType: string;
}

export interface TagAction {
  field: string;
  value: string;
  lookupFlag: boolean;
  lookupTable: string;
  lookupKey: string;
  lookupColumn: string;
}

export interface LambdaFilter {
  script: string;
}

export interface UpdateRecord {
  comment: string;
  username: string;
  createdOn: string;
}

export interface EventWatchSignal {
  type: string;
  name: string;
  field?: string;
  fields?: string[];
  groupBy: string;
  groupByEntity: string;
  keyField?: string;
  caseFlag: boolean;
  caseGroupField?: string;
}

export interface BehaviorRuleT {
  name?: string;
  behavior?: string;
  description?: string;
  keyType?: string;
  key?: string;
  timelineFlag: boolean;
  caseFlag: boolean;
  fields?: string[];
  risks?: string[];
  rules: BehaviorCorrelation[];
  attributes: FieldAttribute[];
  references?: string[];
  mitreTags?: string[];
  falsepositives?: string[];
  level?: string;
  importSource?: string;
}

export interface TranslationT {
  keyType?: string;
  key?: string;
}

export interface BehaviorCorrelation {
  name: string;
  type: string;
  description?: string;
  aggregation?: FieldAggregation;
  first?: FieldOccurrence;
  caseFlag: boolean;
  risks?: string[];
  lastUpdate: string;
}

export interface FieldAggregation {
  window: CorrelationWindow;
  field?: string;
  aggType: string;
  match: ValueMatch;
  anomalyFlag: boolean;
  anomalyScore: number;
}

export interface FieldOccurrence {
  window: CorrelationWindow;
  fields?: string[];
  globalScope?: boolean;
  globalGroup?: string;
  localScope?: boolean;
  entityScope?: boolean;
}

export interface CorrelationRule {
  name: string;
  uniqueName: string;
  description?: string;
  disabled: boolean;
  type: string;
  set?: SignalSetRule;
  aggregation?: SignalAggregationRule;
  first?: OccuranceRule;
  last?: OccuranceRule;
  emit: string;
  text: string;
  keyField: string;
  caseFlag: boolean;
  updatedOn: string;
  createdOn: string;
  createdTs: number;
  tags: string[];
}

export interface OccuranceRule {
  window: CorrelationWindow;
  signal: string;
  fields: string[];
}

export interface CorrelationWindow {
  unit: string;
  length: number;
  text: string;
}

export interface SignalSetRule {
  includes: string[];
  excludes: string[];
  window: CorrelationWindow;
}

export interface ValueMatch {
  operator: string;
  operands: number[];
  text: string;
}

export interface FieldAttribute {
  field: string;
  aliase: string;
  objectFlag: boolean;
  objectField: ObjectFieldSelector;
  entityInfo: string;
}

export interface ObjectFieldSelector {
  keyField: string;
  valueField: string;
  key: string;
}

export interface EventWatchEntityInfo {
  field: string;
  Entities: string[];
}

export interface SignalAggregationRule {
  window: CorrelationWindow;
  signal: string;
  aggType: string;
  match: ValueMatch;
  key: string;
}

/**
 * Match types for `BehaviorEventFilterEntryT.matchType`. Every type except
 * MATCH_TYPE_RANGE compares against the entry's `values` list; MATCH_TYPE_RANGE
 * uses `numberMatch` instead and ignores `values`.
 */
/** Case-insensitive equality. */
export const MATCH_TYPE_STRING = "string_match";
/** Case-insensitive prefix. */
export const MATCH_TYPE_PREFIX = "prefix";
/** Case-insensitive suffix. */
export const MATCH_TYPE_SUFFIX = "suffix";
/** Go RE2 regex; rejected at add/update time if it does not compile. */
export const MATCH_TYPE_REGEX = "regex";
/** Numeric comparison driven by `numberMatch`. */
export const MATCH_TYPE_RANGE = "range";
/** Case-insensitive substring. */
export const MATCH_TYPE_CONTAIN = "contain";
/** Membership in the named entityinfo table. */
export const MATCH_TYPE_ENTITY = "entityinfo";

/** Actions for `BehaviorEventFilterT.action`, applied to a matched behavior event. */
/** Drop the event entirely. */
export const FILTER_ACTION_DISCARD = "discard";
/** Keep the event, exclude it from the behavior summary. */
export const FILTER_ACTION_SKIP_SUMMARY = "skip_summary";
/** Keep the event; use with riskMask/riskAppend to adjust risks only. */
export const FILTER_ACTION_PASS = "pass";

/**
 * One match condition inside a behavior filter. All entries of a filter must
 * match for the filter to fire, unless the filter sets `matchAll`.
 */
export interface BehaviorEventFilterEntryT {
  /**
   * The behavior event field to test: "key", "keyType", "behavior.name",
   * "behavior.group", or a behavior attribute name.
   */
  field: string;
  values: string[];
  /** Required when `matchType` is MATCH_TYPE_RANGE, ignored otherwise. */
  numberMatch?: ValueMatch;
  matchType?: string;
  /** Derived server-side from `field` on add/update; callers do not set it. */
  isAttribute: boolean;
  /** Negates this entry's result. */
  exclude: boolean;
}

/**
 * A behavior event filter, stored per behavior rule. `behaviorRule` plus `name`
 * is the identity: together they form the etcd key
 * `acc_$account/fsm/filters/$behaviorRule/$name`, so `name` only has to be unique
 * inside a rule. A `behaviorRule` of "*" applies the filter to every behavior rule.
 */
export interface BehaviorEventFilterT {
  id: number;
  repository?: string;
  group: string;
  behaviorRule: string;
  name: string;
  disabled: boolean;
  description: string;
  filters: BehaviorEventFilterEntryT[];
  action: string;
  riskMask: string[];
  riskAppend: string[];
  attributes: string[];
  /**
   * Makes the filter fire on every event of the rule, ignoring `filters`.
   * `filters` must still be non-empty: the server rejects a filter with no entries.
   */
  matchAll: boolean;
  whitelist: boolean;
  createdOn: string;
  updatedOn: string;
  lastHit: string;
}
