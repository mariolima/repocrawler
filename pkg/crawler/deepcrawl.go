package crawler

import (
	_ "github.com/mariolima/repocrawl/cmd/utils"
	"github.com/mariolima/repocrawl/internal/entities"

	"gopkg.in/src-d/go-git.v4" //It's def heavy but gets the job done - any alternatives for commit crawling?
	"gopkg.in/src-d/go-git.v4/plumbing/object"
	"gopkg.in/src-d/go-git.v4/plumbing/transport/client"

	_ "gopkg.in/src-d/go-git.v4/storage/memory"
	"io/ioutil"
	"os"

	_ "gopkg.in/src-d/go-billy.v4/memfs" //???????????????????

	"crypto/tls"
	githttp "gopkg.in/src-d/go-git.v4/plumbing/transport/http"
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"

	_ "bufio"
	"strings"

	"fmt" //TODO move funcs that use these `Sprintf` to cmd/utils

	"sync"
	// "os"
)

/*
	Crawls Git repository and retrieves matches with a given `channel`
	This code is trash - need to fix
*/

func (ct *CrawlerTask) DeepCrawl(giturl string) error {
	// Saving repos to Disk as a test - Prevents MEM leak
	dir, err := ioutil.TempDir("", "repo")
	if err != nil {
		log.Fatal(err)
	}
	defer os.RemoveAll(dir)

	// Saving repos to MEM - better performance?
	// storer := memory.NewStorage()
	// r, err := git.Clone(storer, nil, &git.CloneOptions{

	r, err := git.PlainClone(dir, false, &git.CloneOptions{
		URL: giturl,
		// Progress:      os.Stdout,
	})
	if err != nil {
		log.Error("Git Clone Error: ", err)
		return err
	}
	log.Debug("Done cloning ", giturl)

	// map to avoid repeated matches
	var matches = make(map[string]Match)

	// // ... just iterates over the commits, printing it
	// refIter, err := r.Branches()
	// if err != nil {
	// 	log.Fatal("Error getting branches of : ", err)
	// }

	// get the commit object, pointed by ref
	cIter, _ := r.CommitObjects()
	err = cIter.ForEach(func(commit *object.Commit) error {
		log.Trace("Current commit ", commit)
		fIter, _ := commit.Files()
		err = fIter.ForEach(func(cb *object.File) error {
			if ib, _ := cb.IsBinary(); ib {
				return err
			}
			lines, _ := cb.Lines()
			for i, line := range lines {
				log.Trace(line)
				found := ct.RegexLine(line)
				// dumb
				if len(found) > 0 {
					for _, match := range found {
						if match.Rule.Type == "keys" && match.Entropy < 4.3 {
							log.Debug("Dismissed match ", match.Values[0], "due to entropy: ", match.Entropy)
							continue
						}
						// match.SearchResult=result
						log.Debug("Commit:", commit.Hash)
						log.Debug("From ", commit.Author, " ", commit.Message)
						match.URL = commitFileToUrl(giturl, commit.Hash.String(), cb.Name, i+1)
						match.LineNr = i
						// match.URL=fmt.Sprintf("%s/commit/%s",giturl,commit.Hash)
						if _, ok := matches[line]; !ok {
							matches[line] = match
							ct.respChan <- match
							ct.PushMatch(match)
						}
					}
				}
			}
			return err
		})
		return err
	})
	return nil
}

func (c *crawler) DeepCrawl(giturl string, respChan chan Match) error {
	ct := c.NewCrawlerTask(respChan)
	return ct.DeepCrawl(giturl)
}

func setupGitClient() {
	customClient := &http.Client{
		// accept any certificate (might be useful for testing)
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},

		// 15 second timeout
		Timeout: 40 * time.Second,

		// don't follow redirect
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	// Override http(s) default protocol to use our custom client
	client.InstallProtocol("https", githttp.NewClient(customClient))
}

func commitFileToUrl(giturl string, commitHash string, file string, line int) string {
	// why blame? because certain files don't render cleartext (i.e. .md)
	// return fmt.Sprintf("%s/blame/%s/%s#L%d",giturl,commitHash,file,line)
	return fmt.Sprintf("%s/blob/%s/%s#L%d", giturl, commitHash, file, line)
}

func (c *crawler) DeepCrawlGithubRepo(user, repo string, respChan chan Match) {
	users, _ := c.Github.GetRepoContributors(user, repo)
	log.Info("Found ", len(users), " users for repo ", repo)
	for _, user := range users {
		c.DeepCrawlGithubUser(user.Name, respChan)
	}
	log.Warn(":::: DONE crawling users of repo ", user, "/", repo)
}

func (c *crawler) DeepCrawlBitbucketRepo(user, repo string, respChan chan Match) {
	users, _ := c.Bitbucket.GetRepoContributors(user, repo)
	log.Info("Found ", len(users), " users for repo ", repo)
	for _, user := range users {
		log.Trace(user.UUID)
		c.DeepCrawlBitbucketUser(user.UUID, respChan)
	}
	log.Warn(":::: DONE crawling users in repo ", repo)
}

func (c *crawler) DeepCrawlGithubOrg(org string, respChan chan Match) {
	ct := c.NewCrawlerTask(respChan)

	var crawled_users = make(map[string]entities.User)
	var mutex = &sync.Mutex{}

	users, _ := c.Github.GetOrgMembers(org)
	log.Info("Found ", len(users), " members for org ", org)

	ct.DeepCrawlGithubUser(org)

	for _, user := range users {
		ct.AddUser(user)
		go func(user entities.User, respChan chan Match) {
			defer ct.DoneUser(user)
			if strings.Contains(strings.ToUpper(user.Company), strings.ToUpper(org)) {
				log.Warn("User ", user.Company, " has ", org, " in his Bio")
			}
			mutex.Lock()
			if _, ok := crawled_users[user.Name]; ok {
				mutex.Unlock()
				return // avoid deepcrawling same User twice
			} else {
				crawled_users[user.Name] = user
				mutex.Unlock()
			}
			ct.DeepCrawlGithubUser(user.Name)
		}(user, respChan)
	}
	ct.WaitUsers()
	log.Warn("Done with users")

	repos, _ := c.Github.GetUserRepositories(org)
	log.Info(fmt.Sprintf("Found %d repos on Org %s", len(repos), org))
	for i, repo := range repos {
		ct.AddRepo(repo)
		go func(repo entities.Repository) {
			ct.DeepCrawl(repo.GitURL)
			ct.DoneRepo(repo)
		}(repo)

		log.Info("Crawling users of repo [", i, "/", len(repos), "] ", repo.Name)
		users, _ := c.Github.GetRepoContributors(repo.User.Name, repo.Name)
		log.Info("Found ", len(users), " users for repo ", repo.Name)

		// guard := make(chan struct{}, c.Opts.NrThreads)
		for _, user := range users {
			ct.AddUser(user)
			go func(user entities.User, respChan chan Match) {
				defer ct.DoneUser(user)
				if strings.Contains(strings.ToUpper(user.Company), strings.ToUpper(org)) {
					log.Warn("User ", user.Name, " has ", org, " in his Bio")
				}
				mutex.Lock()
				if _, ok := crawled_users[user.Name]; ok {
					// <-guard
					mutex.Unlock()
					return // avoid deepcrawling same User twice
				} else {
					crawled_users[user.Name] = user
					mutex.Unlock()
				}
				ct.DeepCrawlGithubUser(user.Name)
			}(user, respChan)
		}
	}
	ct.WaitUsers()
	ct.WaitRepos()
	log.Warn("---- Done with ORG")
}

func (c *crawler) DeepCrawlBitbucketUser(user string, respChan chan Match) {
	repos, _ := c.Bitbucket.GetUserRepositories(user)
	log.Info(fmt.Sprintf("Found %d repos on User %s", len(repos), user))
	for _, repo := range repos {
		log.Info("DeepCrawling repo ", repo.GitURL)
		c.DeepCrawl(repo.GitURL, respChan)
	}
}

func (ct *CrawlerTask) DeepCrawlGithubUser(user string) {
	repos, _ := ct.Github.GetUserRepositories(user)
	log.Info(fmt.Sprintf("Found %d repos on User %s", len(repos), user))
	for _, repo := range repos {
		ct.AddRepo(repo)
		go func(repo entities.Repository) {
			ct.DeepCrawl(repo.GitURL)
			ct.DoneRepo(repo)
		}(repo)
	}
}

func (c *crawler) DeepCrawlGithubUser(user string, respChan chan Match) {
	ct := c.NewCrawlerTask(respChan)
	ct.DeepCrawlGithubUser(user)
}
