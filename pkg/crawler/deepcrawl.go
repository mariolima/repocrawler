package crawler

import(
	_"github.com/mariolima/repocrawl/internal/entities"
	"github.com/mariolima/repocrawl/cmd/utils"

	"gopkg.in/src-d/go-git.v4"								//It's def heavy but gets the job done - any alternatives for commit crawling?
	"gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
	"gopkg.in/src-d/go-git.v4/storage/memory"

	log "github.com/sirupsen/logrus"

	"bufio"
	"strings"

	"fmt"													//rm later
)


/*
	Crawls Git repository and retrieves matches with a given `channel` 
	This code is trash - need to fix
*/
func (c *crawler) DeepCrawl(giturl string, respChan chan Match) (error) {
	r, err := git.Clone(memory.NewStorage(), nil, &git.CloneOptions{
		URL: giturl,
	})
	// log.Trace(r)
	if err != nil {
		log.Fatal("Error: ", err)
	}

	// ... retrieves the commit history

	// map to avoid repeated matches
	var matches = make(map[string]Match)

	// ... just iterates over the commits, printing it
	refIter, err := r.Branches()
	err = refIter.ForEach(func(cref *plumbing.Reference) error {
		log.Info("Current Branch ",cref)
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
					log.Trace("Going for patch ",p)
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
							found := c.RegexLine(line)
							// dumb
							if len(found) > 0 {
								outp:=chunk.Content()
								for _, match := range found{
									// match.URL=result.FileURL
									outp=utils.HighlightWord(outp, match.Value)
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
		log.Fatal("Error: ", err)
	}
	log.Info("Done with:",giturl)
	return nil
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
}

//same API 
func (c *crawler) DeepCrawlGithubOrg(org string, respChan chan Match) {
	c.DeepCrawlGithubUser(org, respChan)
}

func (c *crawler) DeepCrawlGithubUser(user string, respChan chan Match) {
	repos, _ := c.Github.GetUserRepositories(user)
	log.Info(fmt.Sprintf("Found %d repos on User %s",len(repos), user))
	for _, repo := range repos {
		log.Info("DeepCrawling repo ", repo.GitURL)
		c.DeepCrawl(repo.GitURL,respChan)
	}
}
