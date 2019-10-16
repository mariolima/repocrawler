// Placeholder - Credit to: https://github.com/michenriksen/gitrob/blob/master/core/options.go ^_^

package main

import (
	"flag"
	"os"
)

// Options Used for the Flags given to the CLI application
type Options struct {
	RulesFile         *string `json:"-"`
	WebServer         *bool   `json:"-"`
	OutputFile        *string `json:"-"`
	GitURL            *string `json:"-"`
	GithubSearchQuery *string `json:"-"`
	GithubRepo        *string `json:"-"`
	GithubUser        *string `json:"-"`
	GithubOrg         *string `json:"-"`
	BitbucketHost     *string `json:"-"`
	BitbucketRepo     *string `json:"-"`
	BitbucketUser     *string `json:"-"`
	SlackWebhook      *string `json:"-"`
	NrThreads         *int    `json:"-"`
	BitbucketCreds    *BitbucketCreds
	WebServerOpts     *WebServer
}

// BitbucketCreds Used to pass Bitbucket credentials if needed
type BitbucketCreds struct {
	Username *string `json:"-"`
	Password *string `json:"-"`
}

// WebServer - Options used to setup the Web MatchServer
type WebServer struct {
	Hostname    *string `json:"-"`
	Port        *int    `json:"-"`
	CertsFolder *string `json:"-"`
	SSL         *bool   `json:"-"`
}

// ParseOptions Given Options struct, parses the CLI options and gives them default values
func ParseOptions() (Options, error) {
	options := Options{
		OutputFile:        flag.String("o", "output.txt", "File Output for raw matches stdout"),
		RulesFile:         flag.String("r", "rules.json", "Json file with all the regexes"),
		GitURL:            flag.String("git", "", "Crawls single repository given a .git Url"),
		GithubSearchQuery: flag.String("q", "", "Search GitHub for code containing specified query and match content for secrets"),
		GithubRepo:        flag.String("githubrepo", "", "DeepCrawls github repository and all repositories of it's contributors (format: user/repo)"),
		GithubUser:        flag.String("githubuser", "", "DeepCrawls all github repositories of given user"),
		GithubOrg:         flag.String("githuborg", "", "DeepCrawls github Org"),
		BitbucketHost:     flag.String("bitbuckethost", "https://api.bitbucket.org/2.0", "Bitbucket base API host"),
		BitbucketRepo:     flag.String("bitbucketrepo", "", "DeepCrawls bitbucket repository and all repositories of it's contributors (format: user/repo)"),
		BitbucketUser:     flag.String("bitbucketuser", "", "DeepCrawls all bitbucket repositories of given user"),
		NrThreads:         flag.Int("n", 5, "Number of threads to be used during DeepCrawling"),
		WebServer:         flag.Bool("w", false, "Run webserver page"),
	}

	options.BitbucketCreds = &BitbucketCreds{
		Username: flag.String("bitbucketusername", "", "Bitbucket Username"),
		Password: flag.String("bitbucketpassword", "", "Bitbucket Password"),
	}

	options.WebServerOpts = &WebServer{
		Hostname:    flag.String("h", "0.0.0.0", "Hostame used for Web MatchServer interface"),
		CertsFolder: flag.String("certs", "configs/certs/", "Location path of the certificate files to be used in the Web MatchServer (.key and .crt)"),
		Port:        flag.Int("p", 8090, "Port to be used for the Web Matchserver HTTP/WS connections"),
		SSL:         flag.Bool("ssl", false, "Run Web Matchserver with SSL"),
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
