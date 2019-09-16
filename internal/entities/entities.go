/*
	Entities need to be as generic as possible in order to work with GitHub, Gitlab, Bitbucket, etc
*/

package entities

import (
	"fmt"
	// "time"
)

type Repository struct {
	GitURL string `json:"giturl"`
	Name   string `json:"name"`
	User   User   `json:"user"`
}

type User struct {
	Name    string `json:"name"`
	Company string `json:"company"`
	UUID    string `json:"uuid"`
}

type SearchResult struct {
	Repository  Repository
	FileURL     string
	FileContent string
}

func (sr SearchResult) String() string {
	return fmt.Sprintf("FileURL: %s\nRepository:\n\tGitURL: %s\n\tName:%s\n\tUser:%s", sr.FileURL, sr.Repository.GitURL, sr.Repository.Name, sr.Repository.User)
}

func (u User) GetName() string {
	return u.Name
}

func (u Repository) GetName() string {
	return u.Name
}
