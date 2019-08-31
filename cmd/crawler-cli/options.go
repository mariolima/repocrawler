// Placeholder - Credit to: https://github.com/michenriksen/gitrob/blob/master/core/options.go ^_^

package main

import (
	"flag"
	"os"
)

type Options struct {
	RulesFile				*string		`json:"-"`
	GitUrl					*string		`json:"-"`
	GithubSearchQuery		*string		`json:"-"`
	GithubRepo				*string		`json:"-"`
	GithubUser				*string		`json:"-"`
	GithubOrg				*string		`json:"-"`
	BitbucketHost			*string		`json:"-"`
	BitbucketRepo			*string		`json:"-"`
	BitbucketCreds			*BitbucketCreds
}

type BitbucketCreds struct{
	Username				*string `json:"-"`
	Password				*string `json:"-"`
}

func ParseOptions() (Options, error) {
  options := Options{
		RulesFile:			flag.String("rulesfile", "rules.json", "Json file with all the regexes"),
		GitUrl:				flag.String("git", "", "Crawls single repository given a .git Url"),
		GithubSearchQuery:	flag.String("ghq", "", "Search all of GitHub for specified query and match for secrets"),
		GithubRepo:			flag.String("githubrepo", "", "DeepCrawls github repository and all repositories of it's contributors (format: user/repo)"),
		GithubUser:			flag.String("githubuser", "", "DeepCrawls all github repositories of given user"),
		GithubOrg:			flag.String("githuborg", "", "DeepCrawls github Org"),
		BitbucketHost:		flag.String("bitbuckethost", "https://api.bitbucket.org/2.0", "Bitbucket base API host"),
		BitbucketRepo:		flag.String("bitbucketrepo", "", "DeepCrawls bitbucket repository and all repositories of it's contributors (format: user/repo)"),
  }

  options.BitbucketCreds = &BitbucketCreds{
	  Username:				flag.String("bitbucketusername", "", "Bitbucket Username"),
	  Password:				flag.String("bitbucketpassword", "", "Bitbucket Password"),
  }

  flag.Parse()
  // options.Logins = flag.Args()

  return options, nil
}

func getEnv(key, fallback string) string {
    if value, ok := os.LookupEnv(key); ok {
        return value
    }
    return fallback
}
