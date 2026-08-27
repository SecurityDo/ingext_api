package model

import (
	"encoding/json"
	"time"
)

type BehaviorSummary struct {
	ID       string `json:"id" bson:"id"`
	From     int64  `json:"from" bson:"from"`
	To       int64  `json:"to" bson:"to"`
	Count    int    `json:"count" bson:"count"`
	Key      string `json:"key" bson:"key"`
	KeyType  string `json:"keyType,omitempty" bson:"keyType,omitempty"`
	DayIndex string `json:"dayIndex,omitempty" bson:"dayIndex,omitempty"`
	Interval string `json:"interval,omitempty" bson:"interval,omitempty"`

	BehaviorRules []string `json:"behaviorRules" bson:"behaviorRules"`
	Behaviors     []string `json:"behaviors" bson:"behaviors"`
	Risks         []string `json:"risks,omitempty" bson:"risks,omitempty"`
	RiskScore     int      `json:"riskScore" bson:"riskScore"`
	MitreTags     []string `json:"mitreTags,omitempty" bson:"mitreTags,omitempty"`

	SummaryList []*BehaviorRuleSummary `json:"summaryList,omitempty" bson:"summaryList,omitempty"`
	// Slots       []*TimeSlotSummary     `json:"slots,omitempty" bson:"slots,omitempty"`

	ScoreLevel string `json:"scoreLevel" bson:"scoreLevel"`

	// KeyContext *recommend.EntityContext `json:"keyContext" bson:"keyContext"`

	// incident management
	Comments              []*UserComment `json:"comments" bson:"comments"`
	Status                string         `json:"status,omitempty" bson:"status,omitempty"`
	Incident              bool           `json:"incident" bson:"incident"`
	IncidentDetectionTime time.Time      `json:"incidentDetectionTime,omitempty" bson:"incidentDetectionTime,omitempty"`
	ScoreAdjust           int            `json:"scoreAdjust"`
	UpdatedOn             int64          `json:"updatedOn" bson:"updatedOn"`

	NotifyFlag bool `json:"-"`

	// Investigations []*Investigation `json:"investigations,omitempty" bson:"investigations,omitempty"`

	// openAI embedding APIs
	VectorData      []float64 `json:"vectorData,omitempty" bson:"vectorData,omitempty"`
	FingerprintHash string    `json:"fingerprintHash,omitempty" bson:"fingerprintHash,omitempty"`
	Fingerprint     string    `json:"fingerprint,omitempty" bson:"fingerprint,omitempty"`
	RuleFP          string    `json:"ruleFP,omitempty" bson:"ruleFP,omitempty"`
}

type UserComment struct {
	Content  string   `json:"content"`
	Actions  []string `json:"actions"`
	Username string   `json:"username"`
	// NewState  string    `json:"newState"`
	CreatedOn int64 `json:"createdOn"`
}

type BehaviorRuleSummary struct {
	BehaviorRule       string                `json:"behaviorRule" bson:"behaviorRule"`
	Behavior           string                `json:"behavior" bson:"behavior"`
	From               int64                 `json:"from" bson:"from"`
	To                 int64                 `json:"to" bson:"to"`
	Count              int                   `json:"count" bson:"count"`
	AttributeSummaries []*ValueSummmaryEntry `json:"attributeSummaries" bson:"attributeSummaries"`
	Hits               []*BehaviorRuleHit    `json:"hits,omitempty" bson:"hits,omitempty"`
	RiskScore          int                   `json:"riskScore" bson:"riskScore"`
	//scoreSet           *risk.ScoreSet                              `json:"-" bson:"-"`
	//hitMap             map[string][]*eventWatchDao.BehaviorRuleHit `json:"-" bson:"-"`
	fieldMap map[string]*ValueSummmaryEntry `json:"-" bson:"-"`
	Risks    []string                       `json:"risks" bson:"risks"`
}

type BehaviorRuleHit struct {
	Name        string   `json:"name" bson:"name"`
	Description string   `json:"description,omitempty" bson:"description,omitempty"`
	Fields      []string `json:"fields,omitempty" bson:"fields,omitempty"`
	Values      []string `json:"values,omitempty" bson:"values,omitempty"`
	Scope       string   `json:"scope,omitempty" bson:"scope,omitempty"`
	Risks       []string `json:"risks,omitempty" bson:"risks,omitempty"`
}

type ValueSummmaryEntry struct {
	Key    string   `json:"key" bson:"key"`
	Aliase string   `json:"aliase,omitempty" bson:"aliase,omitempty"`
	Values []string `json:"values" bson:"values"`
	//	Value interface{} `json:"value"`
}

type ValueEntry struct {
	Key    string `json:"key"`
	Aliase string `json:"aliase,omitempty"`
	Value  string `json:"value"`
	//	Value interface{} `json:"value"`
}

type BehaviorEvent struct {
	// Sequence, OriginalKey, OriginalKeyType, Description, ScoreLevel and
	// AttributeMap are what eventwatch_rule_test returns on a hit; a behavior
	// event read back from a search carries the rest. AttributeMap is kept raw:
	// it has been null in every response seen so far. Note its capitalized JSON
	// key, which is the endpoint's, not a typo here.
	Sequence        int64              `json:"sequence,omitempty"`
	Timestamp       int64              `json:"timestamp"`
	Key             string             `json:"key"`
	OriginalKey     string             `json:"originalKey,omitempty"`
	OriginalKeyType string             `json:"originalKeyType,omitempty"`
	KeyType         string             `json:"keyType,omitempty"`
	Title           string             `json:"title,omitempty"`
	Description     string             `json:"description,omitempty"`
	BehaviorRule    string             `json:"behaviorRule"`
	Behavior        string             `json:"behavior"`
	RiskScore       int                `json:"riskScore"`
	ScoreLevel      string             `json:"scoreLevel,omitempty"`
	AttributeMap    json.RawMessage    `json:"AttributeMap,omitempty"`
	Attributes      []*ValueEntry      `json:"attributes"`
	Risks           []string           `json:"risks,omitempty"`
	RuleRisks       []string           `json:"ruleRisks,omitempty"`
	RuleHits        []*BehaviorRuleHit `json:"ruleHits,omitempty"`
}

type EventWatchBucket struct {
	ID         int64  `json:"id"`
	Repository string `json:"repository"`

	//Id            bson.ObjectId `bson:"_id,omitempty" json:"_id"`
	Name              string `json:"name,omitempty"`
	Group             string `json:"group"`
	Description       string `json:"description,omitempty"`
	Disabled          bool   `json:"disabled"`
	EnableTotal       bool   `json:"enableTotal"`
	EventType         string `json:"eventType"`
	ResourceType      string `json:"resourceType"`
	AggregationBucket bool   `json:"aggregationBucket"`
	IsProcessor       bool   `json:"isProcessor"`
	Processor         string `json:"processor"`
	RuleGroup         string `json:"ruleGroup"` // sigma

	// for compatibility
	BucketType string   `json:"bucketType,omitempty"` //  "user"  ||  "system"
	KeyType    string   `json:"keyType,omitempty"`
	TimeSlices []string `json:"timeSlices,omitempty"`

	Fields        []string    `json:"fields,omitempty"`
	EventSelector *EventMatch `json:"eventSelector"`

	// GroupBy string `json:"groupBy"`

	Aggregations []*Aggregation      `json:"aggregations"`
	TagActions   []*TagAction        `json:"tagActions"`
	Signals      []*EventWatchSignal `json:"signals"`

	BehaviorRule *BehaviorRuleT `json:"behaviorRule"`
	Translation  *TranslationT  `json:"translation"`

	//Rules []string `json:"rules,omitempty"`

	UpdatedOn  time.Time `json:"updatedOn"`
	CreatedOn  time.Time `json:"createdOn"`
	DeployedOn time.Time `json:"deployedOn"`

	Tags          []string `json:"tags"`
	SearchProfile string   `json:"searchProfile"`
	// GroupByEntity string    `json:"groupByEntity"`
	Discard bool `json:"discard"`

	Comment string          `json:"comment,omitempty"`
	History []*UpdateRecord `json:"history"`
}

type EventWatchSignalRef struct {
	Behavior      *BehaviorRuleT    `json:"behavior"`
	Signal        *EventWatchSignal `json:"signal"`
	Bucket        string            `json:"bucket"`
	SearchProfile string            `json:"searchProfile"`
	EventType     string            `json:"eventType"`
	// GroupByEntity string    `json:"groupByEntity"`
	GroupBy string `json:"groupBy"`

	// for correlation signals
	SignalSearchQuery string `json:"signalSearchQuery"`
	Query             string `json:"query"`
	SignalName        string `json:"signalName"`

	FieldEntity []*EventWatchEntityInfo `json:"fieldEntity,omitempty"`
}

type Aggregation struct {
	Type  string `json:"type"` // count,sum
	Name  string `json:"name"`
	Field string `json:"field,omitempty"`

	GroupBy       string `json:"groupBy"`
	GroupByEntity string `json:"groupByEntity"`

	// Keys []*FirstKeyEntry `json:"keys,omitempty"`
	AggregationBucket   string `json:"aggregationBucket"`
	AggregationKey      string `json:"aggregationKey"`
	AggregationName     string `json:"aggregationName"`
	ExternalAggregation bool   `json:"externalAggregation"`
}

type EventMatch struct {
	MatchAll    bool            `json:"matchAll"`
	EventFilter string          `json:"eventFilter"`
	Query       json.RawMessage `json:"query,omitempty"`

	MustFilters    []*BucketFilterEntry `json:"mustFilters"`
	MustNotFilters []*BucketFilterEntry `json:"mustNotFilters"`

	JsonQuery  json.RawMessage `json:"jsonQuery,omitempty"`
	JsonFilter *LambdaFilter   `json:"jsonFilter,omitempty"`
	SigmaQuery json.RawMessage `json:"sigmaQuery,omitempty"`

	LVDBQuery     json.RawMessage `json:"lvdbQuery,omitempty"`
	LVDBQueryFlag bool            `json:"lvdbQueryFlag"`
	LVDBQueryTree json.RawMessage `json:"lvdbQueryTree,omitempty"`
}

type BucketFilterEntry struct {
	Field string   `json:"field"`
	Terms []string `json:"terms"`
	// Terms []interface{} `json:"terms"`
	// "field", "eventinfo", "feed"
	FilterType string `json:"filterType"`
}

type TagAction struct {
	Field string `json:"field"`
	Value string `json:"value"`

	LookupFlag   bool   `json:"lookupFlag"`
	LookupTable  string `json:"lookupTable"`
	LookupKey    string `json:"lookupKey"`
	LookupColumn string `json:"lookupColumn"`
}

type LambdaFilter struct {
	Script string `json:"script"`
}

type UpdateRecord struct {
	Comment   string    `json:"comment"`
	Username  string    `json:"username"`
	CreatedOn time.Time `json:"createdOn"`
}

type EventWatchSignal struct {
	Type  string `json:"type"` // count,sum, unique,
	Name  string `json:"name"`
	Field string `json:"field,omitempty"`
	//Keys   []*FirstKeyEntry `json:"keys,omitempty"`
	Fields []string `json:"fields,omitempty"`

	// for aggregation type OR hit type (optional)
	GroupBy       string `json:"groupBy"`
	GroupByEntity string `json:"groupByEntity"`

	// for 'case' type signal
	KeyField       string `json:"keyField,omitempty"`
	CaseFlag       bool   `json:"caseFlag"`
	CaseGroupField string `json:"caseGroupField,omitempty"`
}

type BehaviorRuleT struct {
	Name        string `json:"name,omitempty"`
	Behavior    string `json:"behavior,omitempty"`
	Description string `json:"description,omitempty"`

	KeyType string `json:"keyType,omitempty"`
	Key     string `json:"key,omitempty"`

	// AlternativeKeys []*BehaviorRuleKey `json:"alternativeKeys,omitempty"`

	TimelineFlag bool `json:"timelineFlag"`
	CaseFlag     bool `json:"caseFlag"`

	//only if there is occurrence rule
	Fields []string `json:"fields,omitempty"`
	//for print purpose
	//KeyFields []string `json:"keyFields,omitempty"`
	Risks      []string               `json:"risks,omitempty"`
	Rules      []*BehaviorCorrelation `json:"rules"`
	Attributes []*FieldAttribute      `json:"attributes"`

	References     []string `json:"references,omitempty"`
	MitreTags      []string `json:"mitreTags,omitempty"`
	FalsePositives []string `json:"falsepositives,omitempty"`
	Level          string   `json:"level,omitempty"`
	ImportSource   string   `json:"importSource,omitempty"`
}

type TranslationT struct {
	KeyType string `json:"keyType,omitempty"`
	Key     string `json:"key,omitempty"`
}

type BehaviorCorrelation struct {
	Name        string `json:"name"`
	Type        string `json:"type"`
	Description string `json:"description,omitempty"`

	// Fields []string `json:"fields,omitempty"`
	Aggregation    *FieldAggregation `json:"aggregation,omitempty"`
	FirstOccurance *FieldOccurrence  `json:"first,omitempty"`

	CaseFlag bool     `json:"caseFlag"`
	Risks    []string `json:"risks,omitempty"`

	LastUpdate time.Time `json:"lastUpdate"`
}

type FieldAggregation struct {
	Window *CorrelationWindow `json:"window"`
	// Signal  string             `json:"signal"`
	Field        string      `json:"field,omitempty"`
	AggType      string      `json:"aggType"`
	Match        *ValueMatch `json:"match"`
	AnomalyFlag  bool        `json:"anomalyFlag"`
	AnomalyScore float64     `json:"anomalyScore"`

	Rule *CorrelationRule `json:"-"`
}

type FieldOccurrence struct {
	Window      *CorrelationWindow `json:"window"`
	Fields      []string           `json:"fields,omitempty"`
	GlobalScope bool               `json:"globalScope,omitempty"`
	GlobalGroup string             `json:"globalGroup,omitempty"`
	LocalScope  bool               `json:"localScope,omitempty"`
	EntityScope bool               `json:"entityScope,omitempty"`

	GlobalRule *CorrelationRule `json:"-"`
	LocalRule  *CorrelationRule `json:"-"`
	EntityRule *CorrelationRule `json:"-"`
}

type CorrelationRule struct {
	Name        string `json:"name"`
	UniqueName  string `json:"uniqueName"`
	Description string `json:"description,omitempty"`
	Disabled    bool   `json:"disabled"`

	Type           string                 `json:"type"`
	SignalSet      *SignalSetRule         `json:"set,omitempty"`
	Aggregation    *SignalAggregationRule `json:"aggregation,omitempty"`
	FirstOccurance *OccuranceRule         `json:"first,omitempty"`
	LastOccurance  *OccuranceRule         `json:"last,omitempty"`
	EmitSignal     string                 `json:"emit"`
	Text           string                 `json:"text"`
	KeyField       string                 `json:"keyField"`
	CaseFlag       bool                   `json:"caseFlag"`

	UpdatedOn time.Time `json:"updatedOn"`
	CreatedOn time.Time `json:"createdOn"`
	CreatedTs int64     `json:"createdTs"`

	Tags []string `json:"tags"`
}

type OccuranceRule struct {
	Window *CorrelationWindow `json:"window"`
	Signal string             `json:"signal"`
	Fields []string           `json:"fields"`
	// Key    string             `json:"key"` // groupBy field, optional
}
type CorrelationWindow struct {
	Unit   string `json:"unit"`
	Length int    `json:"length"`
	Text   string `json:"text"`
}

type SignalSetRule struct {
	Includes []string           `json:"includes"`
	Excludes []string           `json:"excludes"`
	Window   *CorrelationWindow `json:"window"`
	// KeyField string             `json:"keyField"`
}

type ValueMatch struct {
	Operator string    `json:"operator"`
	Operands []float64 `json:"operands"`
	Text     string    `json:"text"`
}

type FieldAttribute struct {
	Field  string `json:"field"`
	Aliase string `json:"aliase"`

	ObjectFlag  bool                `json:"objectFlag"`
	ObjectField ObjectFieldSelector `json:"objectField"`

	EntityInfo string `json:"entityInfo"`
	// PrimaryFlag bool   `json:"primaryFlag"`
}

type ObjectFieldSelector struct {
	KeyField   string `json:"keyField"`
	ValueField string `json:"valueField"`
	Key        string `json:"key"`
}

type EventWatchEntityInfo struct {
	Field    string   `json:"field"`
	Entities []string `json:"Entities"`
}

type SignalAggregationRule struct {
	Window  *CorrelationWindow `json:"window"`
	Signal  string             `json:"signal"`
	AggType string             `json:"aggType"`
	Match   *ValueMatch        `json:"match"`
	Key     string             `json:"key"` // groupBy field, optional
}

// Match types for BehaviorEventFilterEntryT.MatchType. Every type except
// MatchTypeRange compares against the entry's Values list; MatchTypeRange uses
// NumberMatch instead and ignores Values.
const (
	MatchTypeString  = "string_match" // case-insensitive equality
	MatchTypePrefix  = "prefix"       // case-insensitive prefix
	MatchTypeSuffix  = "suffix"       // case-insensitive suffix
	MatchTypeRegex   = "regex"        // Go RE2 regex; rejected at add/update time if it does not compile
	MatchTypeRange   = "range"        // numeric comparison driven by NumberMatch
	MatchTypeContain = "contain"      // case-insensitive substring
	MatchTypeEntity  = "entityinfo"   // membership in the named entityinfo table
)

// Actions for BehaviorEventFilterT.Action, applied to a behavior event that the
// filter matches.
const (
	FilterActionDiscard     = "discard"      // drop the event entirely
	FilterActionSkipSummary = "skip_summary" // keep the event, exclude it from the behavior summary
	FilterActionPass        = "pass"         // keep the event; use with RiskMask/RiskAppend to adjust risks only
)

// BehaviorEventFilterEntryT is one match condition inside a behavior filter. All
// entries of a filter must match for the filter to fire, unless the filter sets
// MatchAll.
type BehaviorEventFilterEntryT struct {
	// Field is the behavior event field to test: "key", "keyType",
	// "behavior.name", "behavior.group", or a behavior attribute name.
	Field  string   `json:"field"`
	Values []string `json:"values"`
	// NumberMatch is required when MatchType is MatchTypeRange and ignored otherwise.
	NumberMatch *ValueMatch `json:"numberMatch,omitempty"`
	MatchType   string      `json:"matchType,omitempty"`
	// IsAttribute is derived server-side from Field on add/update; callers do not set it.
	IsAttribute bool `json:"isAttribute"`
	// Exclude negates this entry's result.
	Exclude bool `json:"exclude"`
}

// BehaviorEventFilterT is a behavior event filter, stored per behavior rule.
// BehaviorRule plus Name is the identity: together they form the etcd key
// acc_$account/fsm/filters/$behaviorRule/$name, so Name only has to be unique
// inside a rule. A BehaviorRule of "*" applies the filter to every behavior rule.
type BehaviorEventFilterT struct {
	ID         int64  `json:"id"`
	Repository string `json:"repository,omitempty"`
	Group      string `json:"group"`

	BehaviorRule string                       `json:"behaviorRule"`
	Name         string                       `json:"name"`
	Disabled     bool                         `json:"disabled"`
	Description  string                       `json:"description"`
	Filters      []*BehaviorEventFilterEntryT `json:"filters"`
	Action       string                       `json:"action"`
	RiskMask     []string                     `json:"riskMask"`
	RiskAppend   []string                     `json:"riskAppend"`
	Attributes   []string                     `json:"attributes"`
	// MatchAll makes the filter fire on every event of the rule, ignoring Filters.
	// Filters must still be non-empty: the server rejects a filter with no entries.
	MatchAll  bool `json:"matchAll"`
	Whitelist bool `json:"whitelist"`

	CreatedOn time.Time `json:"createdOn"`
	UpdatedOn time.Time `json:"updatedOn"`
	LastHit   time.Time `json:"lastHit"`
}
