package github

import(
  "github.com/google/go-github/github" //multiple crawling methods need most of these functions
  "golang.org/x/oauth2" //retarded to use this just to inject a fucking Auth header
  _"bufio"
  _"strings"
  "context"

  log "github.com/sirupsen/logrus"
)

type GitHubCrawler struct{
	API_KEY			string
	client			*github.Client
}

func NewGitHubCrawler(api_key string) *GitHubCrawler {
	return &GitHubCrawler{api_key, setupClient(api_key)}
}

func setupClient(GITHUB_ACCESS_TOKEN string) *github.Client {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: GITHUB_ACCESS_TOKEN},
	)
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)
	return client
}

func (c *GitHubCrawler) getResultContent(result github.CodeResult) (string, error){ // https://godoc.org/github.com/google/go-github/github#CodeResult -- single file
	// TODO make print Matches funct
	for _, match := range result.TextMatches {
		log.Debug("Got match FRAGMENT:",*match.Fragment)
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

	repo_content, _, _, err := c.client.Repositories.GetContents(context.Background(), *result.Repository.Owner.Login, *result.Repository.Name, *result.Path, nil)
	if err != nil {
		log.Fatal("Error: %v\n", err)
		return "", err
	}
	//https://github.com/google/go-github/blob/master/github/repos_contents.go#L23
	file_content , err := repo_content.GetContent()
	if err != nil {
		log.Fatal("Error: %v\n", err)
		return "", err
	}
	return file_content, err
}


func (c *GitHubCrawler) crawlResult(result github.CodeResult) { // https://godoc.org/github.com/google/go-github/github#CodeResult -- single file
	// TODO make print Matches funct
	for _, match := range result.TextMatches {
		log.Debug("Got match FRAGMENT:",*match.Fragment)
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

	// repo_content, _, _, err := c.client.Repositories.GetContents(context.Background(), *result.Repository.Owner.Login, *result.Repository.Name, *result.Path, nil)
	// if err != nil {
	// 	log.Fatal("Error: %v\n", err)
	// 	return
	// }
	//https://github.com/google/go-github/blob/master/github/repos_contents.go#L23
	// file_content , err := repo_content.GetContent()
	// log.Trace(file_content)

	// scanner := bufio.NewScanner(strings.NewReader(file_content))

	// Itrating through the file - //TODO fix this later in new funct?
	// for scanner.Scan() {
	// 	line := scanner.Text()
	// 	log.Trace("Content:",line)
	// 	found := RegexLine(line)
	// 	// dumb
	// 	if len(found) > 0 {
	// 		// log.WithFields(log.Fields{
	// 		// 	"url": *result.HTMLURL,
	// 		// 	"Pattern": found,
	// 		// }).Info("Pattern found")
	// 		log.Info(found," ",*result.HTMLURL)
	// 	}
	// }
}

func (c *GitHubCrawler) SearchCode(q string) {
	//-----
	// query := "hx.spiderfoot.net"
	// fmt.Print("Enter GitHub code search query: ")
	// fmt.Scanf("%s", &query)

	log.WithFields(log.Fields{
		"query": q,
	}).Debug("Query was made")

	// results, err := GithubCodeSearch(query)
	results, _, err := c.client.Search.Code(context.Background(), q, &github.SearchOptions{
		TextMatch: true,
		ListOptions: github.ListOptions{ 1, 1000 },
	})

	if err != nil {
		log.Fatal("Error: %v\n", err)
		return
	}

	//fmt.Printf("Total:%d\nIncomplete Results:%v\n",*results.Total,*results.IncompleteResults)
	log.WithFields(log.Fields{
		"total": *results.Total,
		"IncompleteResults": *results.IncompleteResults,
	}).Debug("Results")

	for _, result := range results.CodeResults {
		log.WithFields(log.Fields{
			"URL": *result.HTMLURL,
			"Path":*result.Path,
			"Description":result.Repository.GetDescription(),
		}).Debug("Result")
		c.crawlResult(result)
	}
}

func (c *GitHubCrawler) GetUsersRepositories(user string){
	//TODO
}

func (c *GitHubCrawler) GetUsersOrganizations(user string){
	//TODO
}

// Deepcrawl
func (c *GitHubCrawler) DeepCrawlRepository(repository string){
	//TODO
}

func (c *GitHubCrawler) DeepCrawlUser(user string){
	//TODO
}

func (c *GitHubCrawler) DeepCrawlOrganization(organization string){
	//TODO
}
