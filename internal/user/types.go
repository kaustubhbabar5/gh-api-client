package user

import "github.com/kaustubhbabar5/gh-api-client/pkg/github"

type GetUsersRequest struct {
	Usernames []string `json:"usernames" validate:"required,max=10,min=1,dive,max=39,min=1"`
}

type GetUsersResponse struct {
	Users         []github.User `json:"users"`
	NotFoundUsers []string      `json:"users_not_found"`
	Errors        []error       `json:"errors"`
}
