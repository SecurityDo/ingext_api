package model

// SentinelOne alert-investigation response model — the shape returned by
// /api/ds/investigate_sentinelone_alert. Mirrors the platform's
// sentinelone.AlertInvestigation field for field; the json tags must stay
// identical.
//
// Two things to know before reading a result:
//
//   - Only the alert lookup is fatal. Every other section degrades to a status
//     in CollectionStatus plus a warning, so a successful call can carry a
//     failed section. Read the other sections normally and report the failure.
//
//   - AlertSummary.Raw is never populated over /api/ds: the backend's
//     includeRawAlert knob is not settable on the wire.

import "time"

// Per-section collection outcomes reported in CollectionStatus.
const (
	S1StatusComplete    = "complete"    // everything asked for was collected
	S1StatusTruncated   = "truncated"   // collected, but a cap was hit
	S1StatusPartial     = "partial"     // some sub-queries succeeded, some did not
	S1StatusFailed      = "failed"      // section could not be collected; see Warnings
	S1StatusSkipped     = "skipped"     // caller turned the section off
	S1StatusUnavailable = "unavailable" // the tenant's SentinelOne API does not support it
)

// Threat correlation methods reported in Associations.Method.
const (
	S1MethodExternalID         = "externalId"         // SentinelOne's own identifier links them
	S1MethodStoryline          = "storyline"          // same process lineage
	S1MethodEndpointTimeWindow = "endpointTimeWindow" // same endpoint inside the window only
	S1MethodNone               = "none"               // nothing correlated — a real answer, not an error
)

// Confidence attached to a correlation method, reported in
// Associations.Confidence. This is what says how much to trust Threats.
const (
	S1ConfidenceHigh   = "high"   // paired with S1MethodExternalID
	S1ConfidenceMedium = "medium" // paired with S1MethodStoryline
	S1ConfidenceLow    = "low"    // paired with S1MethodEndpointTimeWindow; may be coincidence
	S1ConfidenceNone   = "none"   // paired with S1MethodNone
)

// AlertInvestigation is the assembled result of investigate_sentinelone_alert.
type AlertInvestigation struct {
	Alert            *AlertSummary          `json:"alert"`
	Associations     *Associations          `json:"associations"`
	Threats          []*ThreatSummary       `json:"threats"`
	Activities       []*ActivitySummary     `json:"activities"`
	Endpoint         *EndpointSummary       `json:"endpoint,omitempty"`
	RelatedAlerts    []*RelatedAlertSummary `json:"relatedAlerts,omitempty"`
	CollectionStatus *CollectionStatus      `json:"collectionStatus"`
}

// AlertSummary is the SentinelOne Unified Alert the investigation is about.
type AlertSummary struct {
	ID              string     `json:"id"`
	Name            string     `json:"name,omitempty"`
	Description     string     `json:"description,omitempty"`
	Severity        string     `json:"severity,omitempty"`
	Status          string     `json:"status,omitempty"`
	Result          string     `json:"result,omitempty"`
	AnalystVerdict  string     `json:"analystVerdict,omitempty"`
	Classification  string     `json:"classification,omitempty"`
	ConfidenceLevel string     `json:"confidenceLevel,omitempty"`
	ExternalID      string     `json:"externalId,omitempty"`
	StorylineID     string     `json:"storylineId,omitempty"`
	NoteExists      bool       `json:"noteExists,omitempty"`
	DetectedAt      *time.Time `json:"detectedAt,omitempty"`
	CreatedAt       *time.Time `json:"createdAt,omitempty"`
	UpdatedAt       *time.Time `json:"updatedAt,omitempty"`
	FirstSeenAt     *time.Time `json:"firstSeenAt,omitempty"`
	LastSeenAt      *time.Time `json:"lastSeenAt,omitempty"`

	DetectionSource *DetectionSource `json:"detectionSource,omitempty"`
	Asset           *AlertAsset      `json:"asset,omitempty"`
	Assignee        *Assignee        `json:"assignee,omitempty"`
	Process         *AlertProcess    `json:"process,omitempty"`

	// Raw is never populated over /api/ds — includeRawAlert is not wire-settable.
	Raw map[string]interface{} `json:"raw,omitempty"`
}

type DetectionSource struct {
	Product string `json:"product,omitempty"`
	Vendor  string `json:"vendor,omitempty"`
	Engine  string `json:"engine,omitempty"`
}

type AlertAsset struct {
	ID        string `json:"id,omitempty"`
	Name      string `json:"name,omitempty"`
	Type      string `json:"type,omitempty"`
	AgentUUID string `json:"agentUuid,omitempty"`
	OSType    string `json:"osType,omitempty"`
}

type Assignee struct {
	FullName string `json:"fullName,omitempty"`
	Email    string `json:"email,omitempty"`
}

type AlertProcess struct {
	CmdLine    string    `json:"cmdLine,omitempty"`
	ParentName string    `json:"parentName,omitempty"`
	File       *FileInfo `json:"file,omitempty"`
}

type FileInfo struct {
	Path   string `json:"path,omitempty"`
	SHA1   string `json:"sha1,omitempty"`
	SHA256 string `json:"sha256,omitempty"`
	MD5    string `json:"md5,omitempty"`
}

// Associations explains how the threats and activities were tied to the alert,
// so a caller can weigh the evidence instead of assuming it. Check Confidence
// before trusting Threats: see the S1Confidence* constants.
type Associations struct {
	Method               string     `json:"method"`
	Confidence           string     `json:"confidence"`
	MatchedOn            string     `json:"matchedOn,omitempty"`
	AttemptedMethods     []string   `json:"attemptedMethods,omitempty"`
	ExternalID           string     `json:"externalId,omitempty"`
	StorylineID          string     `json:"storylineId,omitempty"`
	AgentID              string     `json:"agentId,omitempty"`
	AgentUUID            string     `json:"agentUuid,omitempty"`
	WindowStart          *time.Time `json:"windowStart,omitempty"`
	WindowEnd            *time.Time `json:"windowEnd,omitempty"`
	CandidatesConsidered int        `json:"candidatesConsidered,omitempty"`
	AlertLookupStrategy  string     `json:"alertLookupStrategy,omitempty"`
}

type ThreatSummary struct {
	ID               string       `json:"id"`
	ThreatID         string       `json:"threatId,omitempty"`
	Name             string       `json:"name,omitempty"`
	Classification   string       `json:"classification,omitempty"`
	ConfidenceLevel  string       `json:"confidenceLevel,omitempty"`
	AnalystVerdict   string       `json:"analystVerdict,omitempty"`
	IncidentStatus   string       `json:"incidentStatus,omitempty"`
	MitigationStatus string       `json:"mitigationStatus,omitempty"`
	Storyline        string       `json:"storyline,omitempty"`
	InitiatedBy      string       `json:"initiatedBy,omitempty"`
	ProcessUser      string       `json:"processUser,omitempty"`
	DetectionEngines []string     `json:"detectionEngines,omitempty"`
	File             *FileInfo    `json:"file,omitempty"`
	CreatedAt        *time.Time   `json:"createdAt,omitempty"`
	UpdatedAt        *time.Time   `json:"updatedAt,omitempty"`
	IdentifiedAt     *time.Time   `json:"identifiedAt,omitempty"`
	AgentID          string       `json:"agentId,omitempty"`
	AgentUUID        string       `json:"agentUuid,omitempty"`
	ComputerName     string       `json:"computerName,omitempty"`
	Indicators       []*Indicator `json:"indicators,omitempty"`
	MatchedBy        string       `json:"matchedBy,omitempty"`
}

type Indicator struct {
	Category    string   `json:"category,omitempty"`
	Description string   `json:"description,omitempty"`
	Tactics     []string `json:"tactics,omitempty"`
}

// ActivitySummary is one endpoint activity from the window centered on the
// detection time. Newest first.
type ActivitySummary struct {
	ID                   string                 `json:"id"`
	ActivityType         int64                  `json:"activityType,omitempty"`
	PrimaryDescription   string                 `json:"primaryDescription,omitempty"`
	SecondaryDescription string                 `json:"secondaryDescription,omitempty"`
	Description          string                 `json:"description,omitempty"`
	CreatedAt            *time.Time             `json:"createdAt,omitempty"`
	AgentID              string                 `json:"agentId,omitempty"`
	ThreatID             string                 `json:"threatId,omitempty"`
	SiteID               string                 `json:"siteId,omitempty"`
	AccountID            string                 `json:"accountId,omitempty"`
	UserID               string                 `json:"userId,omitempty"`
	OSFamily             string                 `json:"osFamily,omitempty"`
	MatchedBy            string                 `json:"matchedBy,omitempty"`
	Data                 map[string]interface{} `json:"data,omitempty"`
}

// EndpointSummary is the agent record for the affected asset.
type EndpointSummary struct {
	AgentID              string     `json:"agentId,omitempty"`
	UUID                 string     `json:"uuid,omitempty"`
	ComputerName         string     `json:"computerName,omitempty"`
	Domain               string     `json:"domain,omitempty"`
	OSName               string     `json:"osName,omitempty"`
	OSType               string     `json:"osType,omitempty"`
	OSRevision           string     `json:"osRevision,omitempty"`
	MachineType          string     `json:"machineType,omitempty"`
	AgentVersion         string     `json:"agentVersion,omitempty"`
	GroupName            string     `json:"groupName,omitempty"`
	SiteName             string     `json:"siteName,omitempty"`
	AccountName          string     `json:"accountName,omitempty"`
	LastIPToMgmt         string     `json:"lastIpToMgmt,omitempty"`
	ExternalIP           string     `json:"externalIp,omitempty"`
	LastLoggedInUserName string     `json:"lastLoggedInUserName,omitempty"`
	NetworkStatus        string     `json:"networkStatus,omitempty"`
	ScanStatus           string     `json:"scanStatus,omitempty"`
	MitigationMode       string     `json:"mitigationMode,omitempty"`
	IsActive             bool       `json:"isActive"`
	IsDecommissioned     bool       `json:"isDecommissioned"`
	Infected             bool       `json:"infected"`
	ActiveThreats        int64      `json:"activeThreats"`
	RegisteredAt         *time.Time `json:"registeredAt,omitempty"`
	LastActiveDate       *time.Time `json:"lastActiveDate,omitempty"`
}

// RelatedAlertSummary is a sibling alert sharing the storyline or the endpoint.
// Only collected when IncludeRelatedAlerts is requested.
type RelatedAlertSummary struct {
	ID          string     `json:"id"`
	Name        string     `json:"name,omitempty"`
	Severity    string     `json:"severity,omitempty"`
	Status      string     `json:"status,omitempty"`
	DetectedAt  *time.Time `json:"detectedAt,omitempty"`
	StorylineID string     `json:"storylineId,omitempty"`
	AssetName   string     `json:"assetName,omitempty"`
	MatchedBy   string     `json:"matchedBy,omitempty"`
}

// CollectionStatus reports the per-section outcome using the S1Status*
// constants. A failed section does not invalidate the rest of the
// investigation.
type CollectionStatus struct {
	Threats       string   `json:"threats"`
	Activities    string   `json:"activities"`
	Endpoint      string   `json:"endpoint,omitempty"`
	RelatedAlerts string   `json:"relatedAlerts,omitempty"`
	Warnings      []string `json:"warnings"`
}
