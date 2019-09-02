package crawler

import(
	"github.com/mariolima/repocrawl/cmd/utils"
	"github.com/mariolima/repocrawl/internal/entities"

	"gopkg.in/src-d/go-git.v4"								//It's def heavy but gets the job done - any alternatives for commit crawling?
	"gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
	"gopkg.in/src-d/go-git.v4/plumbing/transport/client"
	"gopkg.in/src-d/go-git.v4/storage/memory"

	"net/http"
	"crypto/tls"
	"time"
	githttp "gopkg.in/src-d/go-git.v4/plumbing/transport/http"

	log "github.com/sirupsen/logrus"

	"bufio"
	"strings"

	"fmt"													//TODO move funcs that use these `Sprintf` to cmd/utils
)


/*
	TODO task system
*/
// type CrawlerTask struct{
// 	reponseChan		chan Match
// 	nthreads		int
// }
//
// func (ct *CrawlerTask) DeepCrawl(giturl string) (error) {
// 	//does the crawling
// 	//blablabla
// 	ct.responseChan<-nil
// 	return nil
// }
//
// func (c *crawler) DeepCrawl(giturl string, respChan chan Match) (error) {
// 	// setup goroutines with c.Opts (nthreads)
// 	// adds task to the list of Tasks in Crawler
// 	c.AddTask(&CrawlerTask{
// 		responseChan: respChan,
// 	})
// 	return nil
// }

/*
	Crawls Git repository and retrieves matches with a given `channel` 
	This code is trash - need to fix
*/
func (c *crawler) DeepCrawl(giturl string, respChan chan Match) (error) {
	setupClient()	// in order to disable SSL checking and timeout
	r, err := git.Clone(memory.NewStorage(), nil, &git.CloneOptions{
		URL: giturl,
	})
	if err != nil {
		log.Error("Git Clone Error: ", err)
		return err
	}
	log.Debug("Done cloning ",giturl)

	// map to avoid repeated matches
	var matches = make(map[string]Match)

	// ... just iterates over the commits, printing it
	refIter, err := r.Branches()
	err = refIter.ForEach(func(cref *plumbing.Reference) error {
		log.Debug("Current Branch ",cref)
		cIter, _ := r.Log(&git.LogOptions{From: cref.Hash()})
		err = cIter.ForEach(func(commit *object.Commit) error {
			log.Trace("Current commit ",commit)
			parent, err:=commit.Parent(0)							//https://godoc.org/gopkg.in/src-d/go-git.v4/plumbing/object#Commit.Parent
			if err==nil{
				// stats, _ :=commit.Stats()
				// log.Info(commit.Hash, ":",parent.Hash)
				patch, _ :=commit.Patch(parent)						//https://godoc.org/gopkg.in/src-d/go-git.v4/plumbing/format/diff#Patch

				file_patches:=patch.FilePatches()
				for _, p := range file_patches{						//https://godoc.org/gopkg.in/src-d/go-git.v4/plumbing/format/diff#FilePatch
					// log.Trace("Going for patch ",p)
					if p.IsBinary() {
						log.Trace("Found binary file, skipping")
						continue									// might add this later with c.Opts
					}
					from, to :=p.Files()
					for _, chunk := range p.Chunks(){
						scanner := bufio.NewScanner(strings.NewReader(chunk.Content()))
						i:=1
						for scanner.Scan() {
							line := scanner.Text()
							log.Trace(line)
							found := c.RegexLine(line)
							// dumb
							if len(found) > 0 {
								outp:=chunk.Content()
								for _, match := range found{
									// match.URL=result.FileURL
									outp=utils.HighlightWords(outp, match.Values)
									// match.SearchResult=result
									log.Debug("Commit:",commit.Hash)
									log.Debug("From ",commit.Author, " ",commit.Message)
									if from != nil {
										match.URL=commitFileToUrl(giturl, commit.Hash.String(), from.Path(),i)
										match.Line=match.Line
									}else {
										match.URL=commitFileToUrl(giturl, commit.Hash.String(), to.Path(),i)
										match.Line=match.Line
									}
									match.LineNr=i
									// match.URL=fmt.Sprintf("%s/commit/%s",giturl,commit.Hash)
									if _, ok := matches[line]; !ok {
										matches[line]=match
										respChan<-match
									}
								}
								log.Trace(outp)
							}
							i+=1
						}
					}
				}

				// // /*
				// // 	Same as above but with diff contents output only 
				// // */
				// scanner := bufio.NewScanner(strings.NewReader(patch.String()))
				// for scanner.Scan() {
				// 	line := scanner.Text()
				// 	found := c.RegexLine(line)
				// 	// dumb
				// 	if len(found) > 0 {
				// 		// log.Debug("Found:",found)
				// 		outp:=patch.String()
				// 		for _, match := range found{
				// 			// match.URL=result.FileURL
				// 			outp=utils.HighlightWord(outp, match.Value)
				// 			// match.SearchResult=result
				// 			respChan<-match
				// 		}
				// 		log.Trace(outp) //-- too much
				// 	}
				// }
				// //log.Info(patch)

			}
			return nil
		})
		return nil
	})
	if err != nil {
		log.Error("Error Getting Repository: ", err)
	}
	log.Info("Done with:",giturl)
	return nil
}

func setupClient() {
	customClient := &http.Client{
		// accept any certificate (might be useful for testing)
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},

		// 15 second timeout
		Timeout: 10 * time.Second,

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
	return fmt.Sprintf("%s/blob/%s/%s#L%d",giturl,commitHash,file,line)
}

func (c *crawler) DeepCrawlGithubRepo(user, repo string, respChan chan Match) {
	users, _ := c.Github.GetRepoContributors(user, repo)
	log.Info("Found ",len(users), " users for repo ", repo)
	for _, user := range users {
		c.DeepCrawlGithubUser(user.Name, respChan)
	}
	log.Warn(":::: DONE crawling users of repo ",user,"/",repo)
}

func (c *crawler) DeepCrawlBitbucketRepo(user, repo string, respChan chan Match) {
	users, _ := c.Bitbucket.GetRepoContributors(user, repo)
	log.Info("Found ",len(users), " users for repo ", repo)
	for _, user := range users {
		log.Trace(user.UUID)
		c.DeepCrawlBitbucketUser(user.UUID, respChan)
	}
	log.Warn(":::: DONE crawling users in repo ",repo)
}

//same API 
func (c *crawler) DeepCrawlGithubOrg(org string, respChan chan Match) {
	// Also works for Orgs
	repos, _ := c.Github.GetUserRepositories(org)
	log.Info(fmt.Sprintf("Found %d repos on Org %s",len(repos), org))
	var crawled_users = make(map[string]entities.User)
	for _, repo := range repos {
		users, _ := c.Github.GetRepoContributors(repo.User.Name, repo.Name)
		log.Info("Found ",len(users), " users for repo ", repo.Name)
		for _, user := range users {
			maxGoroutines := 3
			guard := make(chan struct{}, maxGoroutines)
			go func(user entities.User,respChan chan Match) {
					if strings.Contains(strings.ToUpper(user.Bio),strings.ToUpper(org)) {
						log.Warn("User ", user.Name," has ",org," in his Bio")
					}
					if _, ok := crawled_users[user.Name]; ok{
						return // avoid deepcrawling same User twice
					}
					log.Info("DeepCrawling user ", user.Name)
					c.DeepCrawlGithubUser(user.Name, respChan)
					crawled_users[user.Name]=user
					<-guard
			}(user,respChan)
		}
	}

	// c.DeepCrawlGithubUser(org, respChan)
	log.Warn(":::: DONE crawling Org ",org)
}

func (c *crawler) DeepCrawlBitbucketUser(user string, respChan chan Match) {
	repos, _ := c.Bitbucket.GetUserRepositories(user)
	log.Info(fmt.Sprintf("Found %d repos on User %s",len(repos), user))
	for _, repo := range repos {
		log.Info("DeepCrawling repo ", repo.GitURL)
		c.DeepCrawl(repo.GitURL,respChan)
	}

	// maxGoroutines := 3
    // guard := make(chan struct{}, maxGoroutines)
    //
	// for _, repo := range repos {
	// 	guard <- struct{}{}
	// 	go func(repoUrl string,respChan chan Match) {
	// 		log.Info("DeepCrawling repo ", repoUrl)
	// 		c.DeepCrawl(repoUrl,respChan)
	// 		<-guard
	// 	}(repo.GitURL,respChan)
	// }

	log.Warn(":::: DONE crawling repos of user ",user)
}

func (c *crawler) DeepCrawlGithubUser(user string, respChan chan Match) {
	repos, _ := c.Github.GetUserRepositories(user)
	log.Info(fmt.Sprintf("Found %d repos on User %s",len(repos), user))
	for _, repo := range repos {
		log.Info("DeepCrawling repo ", repo.GitURL)
		c.DeepCrawl(repo.GitURL,respChan)
	}
	log.Warn(":::: DONE crawling repos of user ",user)
}
