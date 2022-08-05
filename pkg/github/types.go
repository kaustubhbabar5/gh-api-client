package github

import (
	"sort"
	"strings"
)

type Users []User

type User struct {
	Name        string `json:"name,omitempty"`
	Login       string `json:"login,omitempty"`
	Company     string `json:"company,omitempty"`
	Followers   int    `json:"followers,omitempty"`
	PublicRepos int    `json:"public_repos,omitempty"`
} // @name User

func (users Users) SortByLogin() Users {
	sort.SliceStable(users, func(i, j int) bool {
		return strings.ToLower(users[i].Login) < strings.ToLower(users[j].Login)
	})
	return users
}
