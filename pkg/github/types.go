package github

type User struct {
	Name        string `json:"name,omitempty"`
	Login       string `json:"login,omitempty"`
	Company     string `json:"company,omitempty"`
	Followers   int    `json:"followers,omitempty"`
	PublicRepos int    `json:"public_repos,omitempty"`
} // @name User
