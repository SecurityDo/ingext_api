package model

import (
	"time"
)

type FluencyAccount struct {
	Disabled    bool     `json:"disabled"`
	MultiTenant bool     `json:"multiTenant"`
	Name        string   `json:"name"`
	Region      string   `json:"region"`
	Cluster     string   `json:"cluster"`
	VPC         string   `json:"vpc"`
	Description string   `json:"description"`
	Subnet      string   `json:"subnet"`
	SubnetID    string   `json:"subnetID"`
	ServerNode  string   `json:"serverNode"`
	WorkerNodes []string `json:"workerNodes"`
	//APIToken    string    `json:"APIToken"`
	Token     string    `json:"token"`
	URL       string    `json:"URL"`
	Oauth2    bool      `json:"oauth2"`
	CreatedOn time.Time `json:"createdOn"`
	UpdatedOn time.Time `json:"updatedOn"`
}

type ListFluencyAccountsResponse struct {
	Accounts []*FluencyAccount `json:"accounts"`
	Entries  []string          `json:"entries"`
}
type GridAddSaasAccountRequest struct {
	//Entry *FluencyAccount `json:"entry,omitempty"`
	Name        string `json:"name"`
	Region      string `json:"region"`
	Cluster     string `json:"cluster"`
	DisplayName string `json:"displayName"`
	SiteURL     string `json:"siteURL"`
	Token       string `json:"token"`
}

type GridDeleteSaasAccountRequest struct {
	//Entry *FluencyAccount `json:"entry,omitempty"`
	Name string `json:"name"`
	//Region      string `json:"region"`
	//Cluster string `json:"cluster"`
}
