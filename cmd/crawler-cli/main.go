package main

import (
  "github.com/google/go-github/github" //multiple crawling methods need most of these functions
  "golang.org/x/oauth2" //retarded to use this just to inject a fucking Auth header
  log "github.com/sirupsen/logrus"
  "context"
  "os"

  "flag" //cli args

  "github.com/mariolima/repocrawl/pkg/crawler"
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
	GITHUB_ACCESS_TOKEN = getEnv("GITHUB_ACCESS_TOKEN","")
	level , err := log.ParseLevel(getEnv("LOG_LEVEL","info"))
	if err != nil {
		level = log.InfoLevel
	}
	log.SetLevel(level)

	query := *flag.String("ghq", "min-saude.pt", "GitHub Query to /search/code")
	flag.Parse()

	log.WithFields(log.Fields{
		"query": query,
	}).Debug("Got Opts:")

	repoCrawler := crawler.NewRepoCrawler(crawler.CrawlerOpts{ GITHUB_ACCESS_TOKEN })

	matches := make(chan crawler.Match)
	go repoCrawler.GithubCodeSearch(query, matches)
	for{
		select{
		case match:=<-matches:
			log.Warn(match.Rule)
		}
	}

	// Create New Api with our auth
	//bitbucket
	// api := gopencils.Api("https://api.bitbucket.org/2.0/")
    //
	// resp := &bitbucket.Repositories{}
    //
	// raw, _ := api.Res("repositories").Res("atlassian",resp).Get()
	// fmt.Printf("%s",raw)

	// raw, _ := api.Res("search").Res("atlassian",resp).Get()

	// Github Code Search
	// api := gopencils.Api("https://api.github.com/")

	// resp := &github.GithubCode{}
	// querystring := map[string]string{"q": "kgg", "per_page": "1000"}
	// api.Res("search").Res("users",resp).Get(querystring)
	// fmt.Printf("%v",resp)
	// fmt.Printf("%s",raw.Raw.Header)



	//-----
	// query := "hx.spiderfoot.net"
	// fmt.Print("Enter GitHub code search query: ")
	// fmt.Scanf("%s", &query)


	// githubClient = setupGithubClient()
    //
    //
}

func FetchOrganizations(username string) ([]*github.Organization, error) {
	client := github.NewClient(nil)
	orgs, _, err := client.Organizations.List(context.Background(), username, nil)
	return orgs, err
}
