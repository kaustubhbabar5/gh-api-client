package user

import (
	"github.com/kaustubhbabar5/gh-api-client/pkg/errors"
	"github.com/kaustubhbabar5/gh-api-client/pkg/github"
)

type GetUsersRequest struct {
	Usernames []string `json:"usernames" validate:"required,max=10,min=1,dive,max=39,min=1" minLength:"1" maxLength:"39"`
} // @name GetUsersRequest

type GetUsersResponse struct {
	Users         github.Users    `json:"users"`
	NotFoundUsers []string        `json:"users_not_found"`
	Errors        errors.JSONErrs `json:"errors"`
} // @name GetUsersResponse
