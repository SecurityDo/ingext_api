package model

import "time"

type UserEntry struct {
	Username     string `json:"username,omitempty" bson:"username"`
	Email        string `json:"email,omitempty" bson:"email"`
	FirstName    string `json:"firstName,omitempty" bson:"firstname"`
	LastName     string `json:"lastName,omitempty" bson:"lastname"`
	Organization string `json:"organization,omitempty" bson:"organization"`
	//Password      string   `json:"password" bson:"password"`

	// admin or analyst
	Roles []string `json:"roles"`

	// "Azure AD" or "Google"
	OAuthProvider string `json:"oauthProvider" bson:"oauthProvider"`
	OAuthFlag     bool   `json:"oauthFlag" bson:"oauthFlag"`

	//MFAProvider string `json:"mfaProvider"`
	//MFAFlag     bool   `json:"mfaFlag"`

	//Restricted bool   `json:"restricted"`
	//Provider   string `json:"provider"`

	DataPolicies []string `json:"dataPolicies,omitempty"`
	APIPolicies  []string `json:"APIPolicies,omitempty"`

	//Profile json.RawMessage `json:"profile"`
}

// Role is an RBAC role that bundles a set of data and API policies.
// Backed by the roleDao endpoint under api/auth.
type Role struct {
	Name         string   `json:"name"`
	Description  string   `json:"description"`
	DataPolicies []string `json:"dataPolicies"`
	APIPolicies  []string `json:"apiPolicies"`
}

// ResourcePolicy grants a set of privileges (e.g. "read", "write") on a resource.
type ResourcePolicy struct {
	Resource   string   `json:"resource"`
	Privileges []string `json:"privileges"`
}

// ApiPolicy is a named collection of resource privileges.
// Backed by the apiPolicyDao endpoint under api/auth.
type ApiPolicy struct {
	Name        string            `json:"name"`
	Description string            `json:"description"`
	Resources   []*ResourcePolicy `json:"resources"`
}

// SitePolicy scopes which accounts a set of users and tokens may access.
// Backed by the sitePolicyDao endpoint under api/auth.
type SitePolicy struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Accounts    []string `json:"accounts"`
	Users       []string `json:"users"`
	Tokens      []string `json:"tokens"`
}

type ApiTokenEntry struct {
	ID          string `json:"id,omitempty"`
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
	Disabled    bool   `json:"disabled"`

	// only apply to smb sites
	// Account string `json:"account,omitempty"`

	Token        string   `json:"token,omitempty"`
	Roles        []string `json:"roles"`
	DataPolicies []string `json:"dataPolicies"`
	APIPolicies  []string `json:"APIPolicies"`

	LastAccess int64     `json:"lastAccess"`
	CreatedOn  time.Time `json:"createdOn"`
	UpdatedOn  time.Time `json:"updatedOn"`
}
