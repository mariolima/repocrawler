package main

import (
  _"fmt"
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
	// rules taken from https://github.com/dxa4481/truffleHogRegexes/blob/master/truffleHogRegexes/regexes.json :)
	rules := map[string]string{
		"secret":"(?i)(secret)",
		"token":"(?i)(token)",
		"password":"(?i)(password)",
		"jdbc":"(?i)(jdbc)",
		"priv_keys":"(?s)(-----BEGIN (RSA|DSA|PGP|EC|) PRIVATE KEY.*END (RSA|DSA|PGP|EC|) PRIVATE KEY-----)",
		"Slack Token": "(xox[p|b|o|a]-[0-9]{12}-[0-9]{12}-[0-9]{12}-[a-z0-9]{32})",
		"RSA private key": "-----BEGIN RSA PRIVATE KEY-----",
		"SSH (DSA) private key": "-----BEGIN DSA PRIVATE KEY-----",
		"SSH (EC) private key": "-----BEGIN EC PRIVATE KEY-----",
		"PGP private key block": "-----BEGIN PGP PRIVATE KEY BLOCK-----",
		"Amazon AWS Access Key ID": "AKIA[0-9A-Z]{16}",
		"Amazon MWS Auth Token": "amzn\\.mws\\.[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}",
		"AWS API Key": "AKIA[0-9A-Z]{16}",
		"Facebook Access Token": "EAACEdEose0cBA[0-9A-Za-z]+",
		"Facebook OAuth": "[f|F][a|A][c|C][e|E][b|B][o|O][o|O][k|K].*['|\"][0-9a-f]{32}['|\"]",
		"GitHub": "[g|G][i|I][t|T][h|H][u|U][b|B].*['|\"][0-9a-zA-Z]{35,40}['|\"]",
		"Generic API Key": "[a|A][p|P][i|I][_]?[k|K][e|E][y|Y].*['|\"][0-9a-zA-Z]{32,45}['|\"]",
		"Generic Secret": "[s|S][e|E][c|C][r|R][e|E][t|T].*['|\"][0-9a-zA-Z]{32,45}['|\"]",
		"Google API Key": "AIza[0-9A-Za-z\\-_]{35}",
		"Google Cloud Platform API Key": "AIza[0-9A-Za-z\\-_]{35}",
		"Google Cloud Platform OAuth": "[0-9]+-[0-9A-Za-z_]{32}\\.apps\\.googleusercontent\\.com",
		"Google Drive API Key": "AIza[0-9A-Za-z\\-_]{35}",
		"Google Drive OAuth": "[0-9]+-[0-9A-Za-z_]{32}\\.apps\\.googleusercontent\\.com",
		"Google (GCP) Service-account": "\"type\": \"service_account\"",
		"Google Gmail API Key": "AIza[0-9A-Za-z\\-_]{35}",
		"Google Gmail OAuth": "[0-9]+-[0-9A-Za-z_]{32}\\.apps\\.googleusercontent\\.com",
		"Google OAuth Access Token": "ya29\\.[0-9A-Za-z\\-_]+",
		"Google YouTube API Key": "AIza[0-9A-Za-z\\-_]{35}",
		"Google YouTube OAuth": "[0-9]+-[0-9A-Za-z_]{32}\\.apps\\.googleusercontent\\.com",
		"Heroku API Key": "[h|H][e|E][r|R][o|O][k|K][u|U].*[0-9A-F]{8}-[0-9A-F]{4}-[0-9A-F]{4}-[0-9A-F]{4}-[0-9A-F]{12}",
		"MailChimp API Key": "[0-9a-f]{32}-us[0-9]{1,2}",
		"Mailgun API Key": "key-[0-9a-zA-Z]{32}",
		"Password in URL": "[a-zA-Z]{3,10}://[^/\\s:@]{3,20}:[^/\\s:@]{3,20}@.{1,100}[\"'\\s]",
		"PayPal Braintree Access Token": "access_token\\$production\\$[0-9a-z]{16}\\$[0-9a-f]{32}",
		"Picatic API Key": "sk_live_[0-9a-z]{32}",
		"Slack Webhook": "https://hooks.slack.com/services/T[a-zA-Z0-9_]{8}/B[a-zA-Z0-9_]{8}/[a-zA-Z0-9_]{24}",
		"Stripe API Key": "sk_live_[0-9a-zA-Z]{24}",
		"Stripe Restricted API Key": "rk_live_[0-9a-zA-Z]{24}",
		"Square Access Token": "sq0atp-[0-9A-Za-z\\-_]{22}",
		"Square OAuth Secret": "sq0csp-[0-9A-Za-z\\-_]{43}",
		"Twilio API Key": "SK[0-9a-fA-F]{32}",
		"Twitter Access Token": "[t|T][w|W][i|I][t|T][t|T][e|E][r|R].*[1-9][0-9]+-[0-9a-zA-Z]{40}",
		"Twitter OAuth": "[t|T][w|W][i|I][t|T][t|T][e|E][r|R].*['|\"][0-9a-zA-Z]{35,44}['|\"]",
	}

	for rule, regex := range rules {
		// matched, err := regexp.Match(regex, []byte(line))
		re := regexp.MustCompile(regex)
		if(re.MatchString(line)) {
			log.Warning("Found:",rule)
			matches=append(matches,line)
		}
		// matches=append(matches, fmt.Sprintf("%v",re.FindAll([]byte(line), 0)))
		// log.Warn(matches)

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
		// dumb
		if len(found) > 0 {
			log.Warning("Pattern found! ",found)
		}
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
