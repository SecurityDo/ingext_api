package api

import (
	"fmt"

	ingextAPI "github.com/SecurityDo/ingext_api/api"
	"github.com/SecurityDo/ingext_api/model"
)

func (c *Client) GetPodRole() (role, arn string, err error) {

	platformService := ingextAPI.NewPlatformService(c.ingextClient)

	role, arn, err = platformService.GetPodRole()

	if err != nil {
		c.Logger.Error("failed to get pod role", "error", err)
		return "", "", fmt.Errorf("failed to get pod role: %w", err)
	}

	return role, arn, nil
}

func (c *Client) TestAssumedRole(roleARN, roleExternalID string) (err error) {

	platformService := ingextAPI.NewPlatformService(c.ingextClient)

	err = platformService.TestAssumedRole(roleARN, roleExternalID)

	if err != nil {
		c.Logger.Error("failed to test assumed role", "error", err, "role", roleARN)
		return fmt.Errorf("failed to add user: %w", err)
	}
	return nil
}
func (c *Client) AddAssumedRole(roleName, roleARN, roleExternalID string) (id string, err error) {

	platformService := ingextAPI.NewPlatformService(c.ingextClient)

	id, err = platformService.AddAssumedRole(roleName, roleARN, roleExternalID)

	if err != nil {
		c.Logger.Error("failed to add assumed role", "error", err, "name", roleName, "role", roleARN)
		return "", fmt.Errorf("failed to add user: %w", err)
	}
	return id, nil
}

func (c *Client) DeleteAssumedRole(roleID string) (err error) {

	platformService := ingextAPI.NewPlatformService(c.ingextClient)

	err = platformService.DeleteAssumedRole(roleID)

	if err != nil {
		c.Logger.Error("failed to delete assumed role", "error", err, "id", roleID)
		return fmt.Errorf("failed to delete user: %w", err)
	}
	return nil
}

func (c *Client) ListAssumedRole() (roles []*model.InstanceRole, err error) {

	platformService := ingextAPI.NewPlatformService(c.ingextClient)

	roles, err = platformService.ListAssumedRole()

	if err != nil {
		c.Logger.Error("failed to list assumed role", "error", err)
		return nil, fmt.Errorf("failed to list assumed role: %w", err)
	}
	return roles, nil
}
