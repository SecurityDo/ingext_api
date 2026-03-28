package model

type EndpointConfig struct {
	Name        string `json:"name"`
	Integration string `json:"integration"`
	Action      string `json:"action"`

	Email *EndpointEmailConfig `json:"email,omitempty"`
	Slack *EndpointSlackConfig `json:"slack,omitempty"`
}

type EndpointEmailConfig struct {
	To []string `json:"to,omitempty"`
	Cc []string `json:"cc,omitempty"`
}

type EndpointSlackConfig struct {
	Channel         string   `json:"channel,omitempty"`
	Channels        []string `json:"channels,omitempty"`
	IntegrationName string   `json:"integrationName,omitempty"`
}
