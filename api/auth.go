package api

import (
	"encoding/json"
	"fmt"

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
		fmt.Printf("Error adding user: %v\n", err.Error())
		return err
	}
	fmt.Println("User added successfully")
	return nil
}

type ListUserResponse struct {
	Users []*model.UserEntry `json:"users"`
}

func (s *AuthService) ListUser() (users []*model.UserEntry, err error) {
	res, err := s.client.GenericCall("api/auth", "userList", nil)
	if err != nil {
		fmt.Printf("Error listing users: %v\n", err.Error())
		return nil, err
	}
	var result ListUserResponse
	err = json.Unmarshal(res.GetBytes(), &result)
	if err != nil {
		fmt.Printf("Error parsing user list response: %v\n", err.Error())
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
		fmt.Printf("Error getting user %s: %v\n", username, err.Error())
		return nil, err
	}
	var result GetUserResponse
	err = json.Unmarshal(res.GetBytes(), &result)
	if err != nil {
		fmt.Printf("Error parsing get user response: %v\n", err.Error())
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
		fmt.Printf("Error deleting user %s: %v\n", username, err.Error())
		return err
	}
	fmt.Println("User deleted successfully")
	return nil
}
