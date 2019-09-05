package github

import (
	"context"
	"github.com/google/go-github/github" //multiple crawling methods need most of these functions
	"golang.org/x/oauth2"                //retarded to use this just to inject a fucking Auth header

	log "github.com/sirupsen/logrus"

	"github.com/mariolima/repocrawl/internal/entities" //structs common in GitHub/GitLab/BitBucket - RepoData/UserData etc
)

type GitHubCrawler struct {
	API_KEY string
	client  *github.Client
}

func NewCrawler(api_key string) *GitHubCrawler {
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

func (c *GitHubCrawler) getResultContent(result github.CodeResult) (string, error) { // https://godoc.org/github.com/google/go-github/github#CodeResult -- single file
	// TODO make print Matches funct
	for _, match := range result.TextMatches {
		log.Debug("Got match FRAGMENT:", *match.Fragment)
		for _, m := range match.Matches {
			log.Trace("Got match:", *m.Text)
			//TODO print matches
		}
	}

	// https://godoc.org/github.com/google/go-github/github#RepositoriesService.GetContents
	// TODO make getContents function
	log.WithFields(log.Fields{
		"*result.Repository.Name":  *result.Repository.Name,
		"*result.Repository.Login": *result.Repository.Owner.Login,
		"*result.Path":             *result.Path,
	}).Debug("GetContents params")

	repo_content, _, _, err := c.client.Repositories.GetContents(context.Background(), *result.Repository.Owner.Login, *result.Repository.Name, *result.Path, nil)
	if err != nil {
		log.Fatal("Error: ", err)
		return "", err
	}
	log.WithFields(log.Fields{
		"Size": *repo_content.Size,
		"Name": *repo_content.Name,
		"Path": *repo_content.Path,
	}).Debug("Repo_content")
	//https://github.com/google/go-github/blob/master/github/repos_contents.go#L23
	file_content, err := repo_content.GetContent()
	if err != nil {
		log.Fatal("Error: %v\n", err)
		return "", err
	}
	return file_content, nil
}

func (c *GitHubCrawler) SearchCode(q string, resp chan entities.SearchResult) { //https://godoc.org/github.com/google/go-github/github#CodeResult
	//https://github.com/google/go-github/blob/master/github/repos_contents.go#L23
	//-----
	// query := "hx.spiderfoot.net"
	// fmt.Print("Enter GitHub code search query: ")
	// fmt.Scanf("%s", &query)

	log.WithFields(log.Fields{
		"query": q,
	}).Info("Query was made")

	// results, err := GithubCodeSearch(query)
	page := 1
	log.Info("Going for page: ", page)
	for {
		results, rsp, err := c.client.Search.Code(context.Background(), q, &github.SearchOptions{
			TextMatch:   true,
			ListOptions: github.ListOptions{Page: page, PerPage: 100}, //max per page is 100 - max pages is 10 - max Results is 1000 -.-
		})
		log.Info("Page [", page, "/", results.GetTotal()/100, "]")

		if err != nil {
			log.Fatal("Error: ", err)
			break
		}

		//fmt.Printf("Total:%d\nIncomplete Results:%v\n",*results.Total,*results.IncompleteResults)
		log.WithFields(log.Fields{
			"total":             *results.Total,
			"IncompleteResults": *results.IncompleteResults,
		}).Debug("Results:")

		for _, result := range results.CodeResults {
			log.WithFields(log.Fields{
				"URL":         *result.HTMLURL,
				"Path":        *result.Path,
				"Description": result.Repository.GetDescription(),
			}).Debug("Result:")
			// c.crawlResult(result, line_res)
			file_content, err := c.getResultContent(result)
			if err != nil {
				log.Fatal("Error: ", err)
				return
			}

			resp <- entities.SearchResult{
				Repository:  c.formatRepo(result.GetRepository()),
				FileURL:     *result.HTMLURL,
				FileContent: file_content,
			}
		}
		if rsp.NextPage == 0 {
			break
		}
		page = rsp.NextPage
		// page+=1
	}
}

func (c *GitHubCrawler) formatRepo(repository *github.Repository) entities.Repository {
	return entities.Repository{
		GitURL: *repository.HTMLURL,
		Name:   *repository.Name,
		User: entities.User{
			Name: *repository.Owner.Login,
		},
	}
}

func (c *GitHubCrawler) formatContributor(contributor *github.Contributor) entities.User {
	//For now
	return entities.User{
		Name: *contributor.Login,
		//TODO from Contributor get User data - 1 more request :(
		// Bio:*contributor.Bio,
	}
}

func (c *GitHubCrawler) formatUser(user *github.User) entities.User {
	//For now
	return entities.User{
		Name: *user.Login,
		Company:  user.GetCompany(),
	}
}

func (c *GitHubCrawler) GetUserRepositories(user string) (repos []entities.Repository, err error) {
	page := 0

	// Using Search
	// for {
	// 	//TODO holy shit this code is bad
	// 	log.Debug("Getting User repos page ",page)
	// 	results, rsp, err := c.client.Search.Repositories(context.Background(), "user:"+user, &github.SearchOptions{
	// 			TextMatch: true,
	// 			ListOptions: github.ListOptions{ Page:page, PerPage:1000 }, //max per page is 100 - max pages is 10 - max Results is 1000 -.-
	// 	})
	// 	if err != nil {
	// 		//TODO check error to see if Connection error or final page
	// 		log.Fatal("Error: ", err)
	// 		return repos, nil
	// 	}
	// 	if rsp.NextPage == 0 {
	// 		break
	// 	}
	// 	for _, repo := range results.Repositories {
	// 		repos=append(repos,c.formatRepo(&repo))
	// 	}
	// 	page=rsp.NextPage
	// 	// page+=1
	// }

	for {
		log.Debug("Getting User repos page ", page)
		//https://godoc.org/github.com/google/go-github/github#RepositoryListOptions
		results, rsp, err := c.client.Repositories.List(context.Background(), user, &github.RepositoryListOptions{
			ListOptions: github.ListOptions{Page: page, PerPage: 1000}, //max per page is 100 - max pages is 10 - max Results is 1000 -.-
		})
		if err != nil {
			log.Fatal("Error: ", err)
			return repos, nil
		}
		for _, repo := range results {
			if !*repo.Fork { //TODO c.Opt to Ignore/Not forks
				repos = append(repos, c.formatRepo(repo))
			}
		}
		if rsp.NextPage == 0 {
			break
		}
		page = rsp.NextPage
		// page+=1
	}
	return repos, nil
}

func (c *GitHubCrawler) GetFollowers(user string) (users []entities.User, err error) {
	results, _, err := c.client.Users.ListFollowers(context.Background(), user, &github.ListOptions{
		Page: 1,
	})
	if err != nil {
		log.Fatal("Error: ", err)
		return users, nil
	}
	for _, user := range results {
		users = append(users, c.formatUser(user))
	}
	return users, nil
}

func (c *GitHubCrawler) GetOrgMembers(org string) (users []entities.User, err error) {
	results, _, err := c.client.Organizations.ListMembers(context.Background(), org, &github.ListMembersOptions{
		PublicOnly: true,
	})
	if err != nil {
		log.Fatal("Error: ", err)
		return users, nil
	}
	for _, user := range results {
		users = append(users, c.formatUser(user))
	}
	page:=0
	search_results, _, err := c.client.Search.Users(context.Background(), org, &github.SearchOptions{
			ListOptions: github.ListOptions{Page: page, PerPage: 100}, //max per page is 100 - max pages is 10 - max Results is 1000 -.-
	})
	if err != nil {
		log.Fatal("Error: ", err)
		return users, nil
	}
	for _, user := range search_results.Users {
		users = append(users, c.formatUser(&user))
	}
	return users, nil
}

func (c *GitHubCrawler) GetRepoContributors(user, repo string) (users []entities.User, err error) {
	page := 1
	// Get all commits
	// commits, _, err := c.client.Repositories.ListCommits(context.Background(), user, repo, &github.CommitsListOptions{
	// 		ListOptions: github.ListOptions{ Page:page, PerPage:1000 }, //max per page is 100 - max pages is 10 - max Results is 1000 -.-
	// })

	//Github already has an API that lists all Contributors - https://developer.github.com/v3/repos/#list-contributors (no need to itereate commits)
	//Contributors LISTED by Number of contributions
	results, _, err := c.client.Repositories.ListContributors(context.Background(), user, repo, &github.ListContributorsOptions{
		ListOptions: github.ListOptions{Page: page, PerPage: 1000}, //max per page is 100 - max pages is 10 - max Results is 1000 -.-
	})
	if err != nil {
		log.Fatal("Error: ", err)
		return users, err
	}
	for _, user := range results {
		users = append(users, c.formatContributor(user))
	}
	return users, err
}

func (c *GitHubCrawler) GetUsersOrganizations(user string) {
	//TODO
}

func (c *GitHubCrawler) DeepCrawlUser(user string) {
	//TODO
}

func (c *GitHubCrawler) DeepCrawlOrganization(organization string) {
	//TODO
}
