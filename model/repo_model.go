package model

import (
	"time"
)

type GithubRepo struct {
	Owner       string    `json:"owner"`
	Repo        string    `json:"repo"`
	Branch      string    `json:"branch"`
	Token       string    `json:"token"`
	Description string    `json:"description"`
	Mode        string    `json:"mode"`
	DisplayName string    `json:"displayName"`
	ID          string    `json:"id"`
	CreatedOn   time.Time `json:"createdOn"`
	UpdatedOn   time.Time `json:"updatedOn"`
}
