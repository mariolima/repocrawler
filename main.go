package main

import (
  "fmt"
  _"github.com/mariolima/repocrawl/bitbucket"
  _"github.com/mariolima/repocrawl/github_own"
  _"github.com/bndr/gopencils"
  "github.com/google/go-github/github"
  "golang.org/x/oauth2" //retarded to use this just to inject a fucking header
  log "github.com/sirupsen/logrus"
  "context"
  "os"

  "regexp"
  "bufio"
  "strings"
  _"strconv"
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


	githubClient = setupGithubClient()

	//-----
	query := "fon.com garrafon"
	// fmt.Print("Enter GitHub code search query: ")
	// fmt.Scanf("%s", &query)

	log.WithFields(log.Fields{
		"query": query,
	}).Debug("Query was made")

	results, err := GithubCodeSearch(query)
	if err != nil {
		log.Fatal("Error: %v\n", err)
		return
	}

	//fmt.Printf("Total:%d\nIncomplete Results:%v\n",*results.Total,*results.IncompleteResults)
	log.WithFields(log.Fields{
		"total": *results.Total,
		"results": *results.IncompleteResults,
	}).Info("Results")

	for _, result := range results.CodeResults {
		log.WithFields(log.Fields{
			"URL": *result.HTMLURL,
			"Path":*result.Path,
			"Description":result.Repository.GetDescription(),
		}).Info("Result")
		GithubCrawlResult(result)
		log.Debug("Crawling Github file #TODO")
	}

}

func GatherRepositoriesFromBitbucketUsername(link string) {

}

func optionsMenu() {

}

func RegexLine(line string) (matches []string) {
	//TODO import rules from cfg
	rules := map[string]string{"SSID":"([A-Z])\\w+"}

	for _, regex := range rules {
		// matched, err := regexp.Match(regex, []byte(line))
		re := regexp.MustCompile(regex)
		matches=append(matches, fmt.Sprintf("%v",re.FindAll([]byte(line), 0)))
		log.Warn(matches)

		// matches = re.FindAll([]byte(line), -1)
		// if err != nil {
		// 	log.Fatal("Error: %v\n", err)
		// 	return false
		// }
		// if matched {
		// 	return true
		// }
	}
	return matches
}

func GithubCrawlResult(result github.CodeResult) { // https://godoc.org/github.com/google/go-github/github#CodeResult -- single file
	// TODO make print Matches funct
	for _, match := range result.TextMatches {
		for _, m := range match.Matches {
			log.Trace("Got match:",*m.Text)
			//TODO print matches
		}
	}

	// https://godoc.org/github.com/google/go-github/github#RepositoriesService.GetContents
	// TODO make getContents function
	log.WithFields(log.Fields{
		"*result.Repository.Name": *result.Repository.Name,
		"*result.Repository.Login": *result.Repository.Owner.Login,
	}).Debug("GetContents params")

	repo_content, _, _, err := githubClient.Repositories.GetContents(context.Background(), *result.Repository.Owner.Login, *result.Repository.Name, *result.Path, nil) 
	if err != nil {
		log.Fatal("Error: %v\n", err)
		return
	}
	//https://github.com/google/go-github/blob/master/github/repos_contents.go#L23
	file_content , err := repo_content.GetContent()
	// log.Trace(file_content)

	scanner := bufio.NewScanner(strings.NewReader(file_content))
	// Itrating through the file - //TODO fix this later in new funct?
	for scanner.Scan() {
		line := scanner.Text()
		log.Trace("Content:",line)
		found := RegexLine(line)
		log.Trace(found)
		// dumb
		// if found != nil {
		// 	log.Warning("Pattern found! ",found)
		// }
	}
}

func GithubCodeSearch(q string) (*github.CodeSearchResult, error) {
	res, _, err := githubClient.Search.Code(context.Background(), q, &github.SearchOptions{
		TextMatch: true,
		ListOptions: github.ListOptions{ 1, 1000 },
	})
	return res, err
}

func FetchOrganizations(username string) ([]*github.Organization, error) {
	client := github.NewClient(nil)
	orgs, _, err := client.Organizations.List(context.Background(), username, nil)
	return orgs, err
}
