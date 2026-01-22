package model

import (
	"time"
)

type InstanceRole struct {
	ID          string    `json:"id"`
	DisplayName string    `json:"displayName"`
	Description string    `json:"description"`
	ExternalID  string    `json:"externalID"`
	RoleARN     string    `json:"roleARN"`
	CreatedOn   time.Time `json:"createdOn"`
}

type InstanceRoleTestRequest struct {
	Names []string `json:"names"`
	// for new role without name
	RoleARN    string `json:"roleARN"`
	ExternalID string `json:"externalID"`
}
