package api

import (
	"errors"

	"github.com/SecurityDo/ingext_api/client"
	"github.com/SecurityDo/ingext_api/model"
)

// SentinelOneService provides helpers for the SentinelOne endpoints under
// /api/ds.
type SentinelOneService struct {
	client *client.IngextClient
}

// NewSentinelOneService constructs a SentinelOneService instance backed by the provided client.
func NewSentinelOneService(client *client.IngextClient) *SentinelOneService {
	return &SentinelOneService{client: client}
}

func (s *SentinelOneService) call(function string, payload interface{}, out interface{}) error {
	return ApiCall(s.client, function, payload, out)
}

// InvestigateAlertOptions are the optional collection controls for
// investigate_sentinelone_alert. Every field has a platform-side default.
//
// The include flags are pointers on purpose: three of them default to true, so
// an omitted field and an explicit false must stay distinguishable. Leave a
// pointer nil to keep the platform default and set it to point at false to turn
// a section off — with plain bools, a request that only set IncludeRelatedAlerts
// would silently disable everything else.
//
// A zero int likewise means "platform default". Out-of-range integers are
// clamped by the platform, not rejected.
type InvestigateAlertOptions struct {
	// IncludeThreats correlates SentinelOne threats to the alert. Default true.
	IncludeThreats *bool `json:"includeThreats,omitempty"`
	// IncludeActivities collects endpoint activities around the detection time.
	// Default true.
	IncludeActivities *bool `json:"includeActivities,omitempty"`
	// IncludeEndpoint includes the agent record for the affected asset. Default true.
	IncludeEndpoint *bool `json:"includeEndpoint,omitempty"`
	// IncludeRelatedAlerts includes other alerts sharing the storyline or
	// endpoint. Default false — it is the slowest section.
	IncludeRelatedAlerts *bool `json:"includeRelatedAlerts,omitempty"`
	// ActivityWindowMinutes is the HALF-WIDTH in minutes of the window centered
	// on the detection time — mitigation and quarantine activity lands after
	// detection, so the window reaches both ways. Default 60, clamped to 1..10080.
	ActivityWindowMinutes int `json:"activityWindowMinutes,omitempty"`
	// MaxActivities caps the activities returned, newest first. Default 200,
	// clamped to 1..1000.
	MaxActivities int `json:"maxActivities,omitempty"`
}

// InvestigateAlertRequest is the kargs payload for investigate_sentinelone_alert.
type InvestigateAlertRequest struct {
	// AlertID is the Unified Alert's own id — not the threat id. Required.
	AlertID string `json:"alertId"`
	// IntegrationName names the SentinelOneEvents integration to query. It may
	// be omitted when exactly one is configured; otherwise the platform's error
	// message enumerates the configured names. Note the integration type is
	// "SentinelOneEvents", not "SentinelOne".
	IntegrationName string `json:"integrationName,omitempty"`
	// Options is optional; omit it to accept every platform default.
	Options *InvestigateAlertOptions `json:"options,omitempty"`
}

// InvestigateAlert calls /api/ds/investigate_sentinelone_alert and returns one
// SentinelOne Unified Alert investigated end to end: the alert itself, the
// threats correlated to it, the endpoint activity around the detection, the
// affected endpoint's agent record, and optionally sibling alerts.
//
// Only the alert lookup is fatal. A returned AlertInvestigation may carry
// failed or partial sections, reported per section in CollectionStatus along
// with a Warnings list — read the other sections normally and report the
// failure. Check Associations.Confidence before trusting the correlated
// threats; Associations.Method == model.S1MethodNone means nothing correlated,
// which is a real answer rather than an error.
//
// The whole call is bounded by a 150s platform-side budget.
func (s *SentinelOneService) InvestigateAlert(request *InvestigateAlertRequest) (*model.AlertInvestigation, error) {
	if request == nil || request.AlertID == "" {
		return nil, errors.New("alertId is required")
	}
	var resp model.AlertInvestigation
	if err := s.call("investigate_sentinelone_alert", request, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// InvestigateAlertByID is InvestigateAlert with the platform defaults for every
// option. Pass an empty integrationName when exactly one SentinelOneEvents
// integration is configured.
func (s *SentinelOneService) InvestigateAlertByID(alertID string, integrationName string) (*model.AlertInvestigation, error) {
	return s.InvestigateAlert(&InvestigateAlertRequest{
		AlertID:         alertID,
		IntegrationName: integrationName,
	})
}
