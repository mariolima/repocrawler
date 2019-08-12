package main

import (
  "github.com/google/go-github/github" //multiple crawling methods need most of these functions
  "golang.org/x/oauth2" //retarded to use this just to inject a fucking Auth header
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
var githubClient *github.Client

func GatherRepositoriesFromBitbucketTeam(team string) {
}

func GatherRepositoriesFromBitbucketUrl(url string) {

}


func getEnv(key, fallback string) string {
    if value, ok := os.LookupEnv(key); ok {
        return value
    }
    return fallback
}

func setupGithubClient() *github.Client {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: GITHUB_ACCESS_TOKEN},
	)
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)
	return client
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

	repoCrawler, _ := crawler.NewRepoCrawler(crawler.CrawlerOpts{
		GithubToken: GITHUB_ACCESS_TOKEN,
	})

	//repoCrawl test
	repoCrawler.DeepCrawl("https://github.com/ptgmiguel/orkos.git")
	return 

	matches := make(chan crawler.Match)
	go repoCrawler.GithubCodeSearch(query, matches)
	for{
		select{
		case match:=<-matches:
			if match.Rule.Type == "keys" {
				coolPrint(match)
				// log.Warn(match)
			}
		}
	}
}

func coolPrint(m crawler.Match) {
	fmt.Printf("[MATCH %s] Line:\t%s\nLink:\t%s\nRepo:\tOwner:%s\tURL:%s\n",m.Rule.Regex,m.Line,m.URL,m.SearchResult.Repository.User.Name,m.SearchResult.Repository.GitURL)
}

func FetchOrganizations(username string) ([]*github.Organization, error) {
	client := github.NewClient(nil)
	orgs, _, err := client.Organizations.List(context.Background(), username, nil)
	return orgs, err
}
