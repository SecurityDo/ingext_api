package model

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
