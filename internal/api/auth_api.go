package api

import (
	"fmt"

	ingextAPI "github.com/SecurityDo/ingext_api/api"
	"github.com/SecurityDo/ingext_api/model"
	ingextModel "github.com/SecurityDo/ingext_api/model"
)

func (c *Client) AddUser(name, displayName, role, org string) error {

	// Use structured logging
	//c.Logger.Info("adding user",
	//	"name", name,
	//	"role", role,
	//)

	authService := ingextAPI.NewAuthService(c.ingextClient)

	err := authService.AddUser(&ingextAPI.AddUserRequest{
		User: &ingextModel.UserEntry{
			Username:     name,
			Email:        name,
			FirstName:    displayName,
			Roles:        []string{role},
			Organization: org,
		},
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
