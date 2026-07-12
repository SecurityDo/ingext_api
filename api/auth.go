package api

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/SecurityDo/ingext_api/client"
	"github.com/SecurityDo/ingext_api/model"
)

type AuthService struct {
	client *client.IngextClient
}

func NewAuthService(client *client.IngextClient) *AuthService {
	return &AuthService{
		client: client,
	}
}

type AddUserRequest struct {
	User *model.UserEntry `json:"user"`
}

func (s *AuthService) AddUser(req *AddUserRequest) error {
	_, err := s.client.GenericCall("api/auth", "userAdd", req)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error adding user: %v\n", err.Error())
		return err
	}
	fmt.Fprintln(os.Stderr, "User added successfully")
	return nil
}

type ListUserResponse struct {
	Users []*model.UserEntry `json:"users"`
}

func (s *AuthService) ListUser() (users []*model.UserEntry, err error) {
	res, err := s.client.GenericCall("api/auth", "userList", nil)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error listing users: %v\n", err.Error())
		return nil, err
	}
	var result ListUserResponse
	err = json.Unmarshal(res.GetBytes(), &result)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing user list response: %v\n", err.Error())
		return nil, err
	}
	return result.Users, nil
}

type GetUserRequest struct {
	Username string `json:"username"`
}

type GetUserResponse struct {
	User *model.UserEntry `json:"user"`
}

func (s *AuthService) GetUser(username string) (*model.UserEntry, error) {
	req := &GetUserRequest{
		Username: username,
	}

	res, err := s.client.GenericCall("api/auth", "getUser", req)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error getting user %s: %v\n", username, err.Error())
		return nil, err
	}
	var result GetUserResponse
	err = json.Unmarshal(res.GetBytes(), &result)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing get user response: %v\n", err.Error())
		return nil, err
	}
	return result.User, nil
}

type DeleteUserRequest struct {
	Username string `json:"username"`
}

func (s *AuthService) DeleteUser(username string) error {
	req := &DeleteUserRequest{
		Username: username,
	}

	_, err := s.client.GenericCall("api/auth", "userDelete", req)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error deleting user %s: %v\n", username, err.Error())
		return err
	}
	fmt.Println("User deleted successfully")
	return nil
}

type tokenRequestArgs struct {
	Name  string               `json:"id,omitempty"`
	Entry *model.ApiTokenEntry `json:"entry,omitempty"`
	Flag  bool                 `json:"flag"`
}
type tokenRequest struct {
	Action string            `json:"action"`
	Args   *tokenRequestArgs `json:"args,omitempty"`
}

type AddTokenResponse struct {
	Token string `json:"token"`
}

func (s *AuthService) AddToken(name, description, role string) (token string, err error) {
	req := &tokenRequest{

		Action: "add",
		Args: &tokenRequestArgs{

			Entry: &model.ApiTokenEntry{
				Name:        name,
				Description: description,
				Roles:       []string{role},
			},
		},
	}

	res, err := s.client.GenericCall("api/auth", "api_token", req)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error adding token %s: %v\n", name, err.Error())
		return "", err
	}
	var result AddTokenResponse
	err = json.Unmarshal(res.GetBytes(), &result)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing add token response: %v\n", err.Error())
		return "", err
	}
	return result.Token, nil
}

func (s *AuthService) DeleteToken(name string) (err error) {
	req := &tokenRequest{
		Action: "delete",
		Args: &tokenRequestArgs{
			Name: name,
		},
	}
	_, err = s.client.GenericCall("api/auth", "api_token", req)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error deleting token %s: %v\n", name, err.Error())
		return err
	}
	return nil
}

type ListTokenResponse struct {
	Entries []*model.ApiTokenEntry `json:"entries"`
}

func (s *AuthService) ListToken() (tokens []*model.ApiTokenEntry, err error) {
	req := &tokenRequest{
		Action: "list",
	}
	res, err := s.client.GenericCall("api/auth", "api_token", req)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error listing token: %v\n", err.Error())
		return nil, err
	}
	var result ListTokenResponse
	err = json.Unmarshal(res.GetBytes(), &result)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing list token response: %v\n", err.Error())
		return nil, err
	}
	return result.Entries, nil
}

type SetUserSitePolicyRequest struct {
	Username   string `json:"username"`
	PolicyName string `json:"policyName"`
}

func (s *AuthService) SetUserSitePolicy(username string, sitePolicy string) (err error) {
	req := &SetUserSitePolicyRequest{
		Username:   username,
		PolicyName: sitePolicy,
	}

	_, err = s.client.GenericCall("api/auth", "setUserSitePolicy", req)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error setting site policy for user %s: %v\n", username, err.Error())
		return err
	}
	return nil
}

// --- Role DAO (backed by the roleDao endpoint under api/auth) ---

// AddRole creates a new RBAC role. A role must reference at least one data or
// API policy.
func (s *AuthService) AddRole(entry *model.Role) error {
	req := &GenericDAORequest[model.Role]{
		Action: "create",
		Args:   &GenericDAORequestArgs[model.Role]{Entry: entry},
	}
	_, err := s.client.GenericCall("api/auth", "roleDao", req)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error adding role: %v\n", err.Error())
		return err
	}
	return nil
}

// ListRole returns all RBAC roles.
func (s *AuthService) ListRole() (roles []*model.Role, err error) {
	req := &GenericDAORequest[model.Role]{Action: "list"}
	res, err := s.client.GenericCall("api/auth", "roleDao", req)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error listing roles: %v\n", err.Error())
		return nil, err
	}
	var result GenericDaoListResponse[model.Role]
	if err = json.Unmarshal(res.GetBytes(), &result); err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing list role response: %v\n", err.Error())
		return nil, err
	}
	return result.Entries, nil
}

// DeleteRole removes the RBAC role identified by name.
func (s *AuthService) DeleteRole(name string) error {
	req := &GenericDAORequest[model.Role]{
		Action: "delete",
		Args:   &GenericDAORequestArgs[model.Role]{Id: name},
	}
	_, err := s.client.GenericCall("api/auth", "roleDao", req)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error deleting role %s: %v\n", name, err.Error())
		return err
	}
	return nil
}

// --- API policy DAO (backed by the apiPolicyDao endpoint under api/auth) ---

// AddApiPolicy creates a new API policy. A policy must define at least one
// resource.
func (s *AuthService) AddApiPolicy(entry *model.ApiPolicy) error {
	req := &GenericDAORequest[model.ApiPolicy]{
		Action: "create",
		Args:   &GenericDAORequestArgs[model.ApiPolicy]{Entry: entry},
	}
	_, err := s.client.GenericCall("api/auth", "apiPolicyDao", req)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error adding api policy: %v\n", err.Error())
		return err
	}
	return nil
}

// ListApiPolicy returns all API policies.
func (s *AuthService) ListApiPolicy() (policies []*model.ApiPolicy, err error) {
	req := &GenericDAORequest[model.ApiPolicy]{Action: "list"}
	res, err := s.client.GenericCall("api/auth", "apiPolicyDao", req)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error listing api policies: %v\n", err.Error())
		return nil, err
	}
	var result GenericDaoListResponse[model.ApiPolicy]
	if err = json.Unmarshal(res.GetBytes(), &result); err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing list api policy response: %v\n", err.Error())
		return nil, err
	}
	return result.Entries, nil
}

// DeleteApiPolicy removes the API policy identified by name.
func (s *AuthService) DeleteApiPolicy(name string) error {
	req := &GenericDAORequest[model.ApiPolicy]{
		Action: "delete",
		Args:   &GenericDAORequestArgs[model.ApiPolicy]{Id: name},
	}
	_, err := s.client.GenericCall("api/auth", "apiPolicyDao", req)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error deleting api policy %s: %v\n", name, err.Error())
		return err
	}
	return nil
}

// --- Site policy DAO (backed by the sitePolicyDao endpoint under api/auth) ---

// AddSitePolicy creates a new site policy. A policy must reference at least one
// user or token.
func (s *AuthService) AddSitePolicy(entry *model.SitePolicy) error {
	req := &GenericDAORequest[model.SitePolicy]{
		Action: "create",
		Args:   &GenericDAORequestArgs[model.SitePolicy]{Entry: entry},
	}
	_, err := s.client.GenericCall("api/auth", "sitePolicyDao", req)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error adding site policy: %v\n", err.Error())
		return err
	}
	return nil
}

// ListSitePolicy returns all site policies.
func (s *AuthService) ListSitePolicy() (policies []*model.SitePolicy, err error) {
	req := &GenericDAORequest[model.SitePolicy]{Action: "list"}
	res, err := s.client.GenericCall("api/auth", "sitePolicyDao", req)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error listing site policies: %v\n", err.Error())
		return nil, err
	}
	var result GenericDaoListResponse[model.SitePolicy]
	if err = json.Unmarshal(res.GetBytes(), &result); err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing list site policy response: %v\n", err.Error())
		return nil, err
	}
	return result.Entries, nil
}

// DeleteSitePolicy removes the site policy identified by name.
func (s *AuthService) DeleteSitePolicy(name string) error {
	req := &GenericDAORequest[model.SitePolicy]{
		Action: "delete",
		Args:   &GenericDAORequestArgs[model.SitePolicy]{Id: name},
	}
	_, err := s.client.GenericCall("api/auth", "sitePolicyDao", req)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error deleting site policy %s: %v\n", name, err.Error())
		return err
	}
	return nil
}
