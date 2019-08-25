package main

import (
  "github.com/google/go-github/github" //multiple crawling methods need most of these functions
  log "github.com/sirupsen/logrus"
  "context"
  "os"

  "flag" //cli args
  "github.com/pkg/profile"

  "github.com/mariolima/repocrawl/pkg/crawler"
  "fmt"
)


const bucket_host string = "https://api.bitbucket.org"

var GITHUB_ACCESS_TOKEN string

func getEnv(key, fallback string) string {
    if value, ok := os.LookupEnv(key); ok {
        return value
    }
    return fallback
}

func main() {
	//Debug
	defer profile.Start().Stop()

	GITHUB_ACCESS_TOKEN = getEnv("GITHUB_ACCESS_TOKEN","")
	level , err := log.ParseLevel(getEnv("LOG_LEVEL","info"))
	if err != nil {
		level = log.InfoLevel
	}
	log.SetLevel(level)

	var query string
	flag.StringVar(&query, "ghq", "min-saude.pt", "GitHub Query to /search/code")
	flag.Parse()

	log.WithFields(log.Fields{
		"query": query,
	}).Info("Got Opts:")

	repoCrawler, _ := crawler.NewRepoCrawler(
		crawler.CrawlerOpts{
			GithubToken: GITHUB_ACCESS_TOKEN,
		},
	)

	//repoCrawl test

	matches := make(chan crawler.Match)

	go repoCrawler.DeepCrawlBitbucketUser("atlassian", matches)

	// go repoCrawler.DeepCrawl("https://bitbucket.org/atlassian/serverless-deploy", matches)
	// go repoCrawler.DeepCrawlGithubRepo("", "", matches)
	// go repoCrawler.DeepCrawlGithubOrg("TwilioDevEd",matches)
	// go repoCrawler.DeepCrawlGithubUser("",matches)

	// go repoCrawler.GithubCodeSearch(query, matches)

	for{
		select{
		case match:=<-matches:
			fmt.Printf("%-30s %-90s %s\n",match.Rule.Regex,match.Line,match.URL)
		}
	}
}

func FetchOrganizations(username string) ([]*github.Organization, error) {
	client := github.NewClient(nil)
	orgs, _, err := client.Organizations.List(context.Background(), username, nil)
	return orgs, err
}
