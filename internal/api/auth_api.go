package api

import (
	"fmt"

	ingextAPI "github.com/SecurityDo/ingext_api/api"
	"github.com/SecurityDo/ingext_api/client"
	"github.com/SecurityDo/ingext_api/model"
	ingextModel "github.com/SecurityDo/ingext_api/model"
)

func (c *Client) AddUser(name, displayName, role, org, oauth string) error {

	// Use structured logging
	//c.Logger.Info("adding user",
	//	"name", name,
	//	"role", role,
	//)

	authService := ingextAPI.NewAuthService(c.ingextClient)

	user := &ingextModel.UserEntry{
		Username:     name,
		Email:        name,
		FirstName:    displayName,
		Roles:        []string{role},
		Organization: org,
	}

	switch oauth {
	case "Azure":
		user.OAuthProvider = "Azure AD"
		user.OAuthFlag = true
	case "Google":
		user.OAuthProvider = "Google"
		user.OAuthFlag = true
	}

	err := authService.AddUser(&ingextAPI.AddUserRequest{
		User: user,
	})
	if err != nil {
		c.Logger.Error("failed to add user", "error", err, "name", name, "role", role)
		return fmt.Errorf("failed to add user: %w", err)
	}
	return nil
}

func (c *Client) DeleteUser(username string) (err error) {

	// Use structured logging

	authService := ingextAPI.NewAuthService(c.ingextClient)

	err = authService.DeleteUser(username)
	if err != nil {
		c.Logger.Error("failed to delete user", "error", err, "name", username)
		return fmt.Errorf("failed to delete user: %w", err)
	}
	return nil
}

func (c *Client) ListUser() (users []*model.UserEntry, err error) {

	// Use structured logging

	authService := ingextAPI.NewAuthService(c.ingextClient)

	users, err = authService.ListUser()
	if err != nil {
		c.Logger.Error("failed to list user", "error", err)
		return nil, fmt.Errorf("failed to list user: %w", err)
	}
	return users, nil
}

// AddToken adds an API token on the configured site (ingextClient).
func (c *Client) AddToken(name, description, role string) (token string, err error) {
	authService := ingextAPI.NewAuthService(c.ingextClient)

	token, err = authService.AddToken(name, description, role)
	if err != nil {
		c.Logger.Error("failed to add token", "error", err, "name", name, "role", role)
		return "", fmt.Errorf("failed to add token: %w", err)
	}
	return token, nil
}

// DeleteToken removes an API token on the configured site (ingextClient).
func (c *Client) DeleteToken(name string) (err error) {
	authService := ingextAPI.NewAuthService(c.ingextClient)

	err = authService.DeleteToken(name)
	if err != nil {
		c.Logger.Error("failed to delete token", "error", err, "name", name)
		return fmt.Errorf("failed to delete token: %w", err)
	}
	return nil
}

// GPTAIAuthService builds an AuthService for the given GPT API base URL and bearer token (e.g. ingext ai --url --token).
// HTTP debug dumps are always disabled for this client. Call AddToken / DeleteToken on the returned service.
func (c *Client) GPTAIAuthService(baseURL, bearerToken string) *ingextAPI.AuthService {
	cli := client.NewIngextClient(baseURL, bearerToken, false, c.Logger)
	return ingextAPI.NewAuthService(cli)
}

func (c *Client) SetUserSitePolicy(username, policy string) error {
	authService := ingextAPI.NewAuthService(c.ingextClient)

	err := authService.SetUserSitePolicy(username, policy)
	if err != nil {
		c.Logger.Error("failed to set user site policy", "error", err, "username", username, "policy", policy)
		return fmt.Errorf("failed to set user site policy: %w", err)
	}
	return nil
}

func (c *Client) ListToken() (tokens []*model.ApiTokenEntry, err error) {

	// Use structured logging

	authService := ingextAPI.NewAuthService(c.ingextClient)

	tokens, err = authService.ListToken()
	if err != nil {
		c.Logger.Error("failed to list token", "error", err)
		return nil, fmt.Errorf("failed to list token: %w", err)
	}
	return tokens, nil
}
