/*
	Entities need to be as generic as possible in order to work with GitHub, Gitlab, Bitbucket, etc
*/

package entities

import (
	"fmt"
	// "time"
)

type Repository struct {
	GitURL string
	Name   string
	User   User
}

type User struct {
	Name    string
	Company string
	UUID    string
}

type SearchResult struct {
	Repository  Repository
	FileURL     string
	FileContent string
}

func (sr SearchResult) String() string {
	return fmt.Sprintf("FileURL: %s\nRepository:\nGitURL: %s\nName:%s\nUser:%s", sr.FileURL, sr.Repository.GitURL, sr.Repository.Name, sr.Repository.User)
}

func (u User) GetName() string {
	return u.Name
}

func (u Repository) GetName() string {
	return u.Name
}
