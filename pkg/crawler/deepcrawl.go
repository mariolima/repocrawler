package crawler

import(
	_"github.com/mariolima/repocrawl/internal/entities"

	"gopkg.in/src-d/go-git.v4"								//It's def heavy but gets the job done - any alternatives for commit crawling?
	"gopkg.in/src-d/go-git.v4/plumbing/object"
	"gopkg.in/src-d/go-git.v4/storage/memory"

	log "github.com/sirupsen/logrus"

	"bufio"
	"strings"
)


/*
	Crawls Git repository and retrieves matches with a given `channel` 
	Head crawling only for now 
*/
func (c *crawler) DeepCrawl(giturl string, respChan chan Match) (error) {
	r, err := git.Clone(memory.NewStorage(), nil, &git.CloneOptions{
		URL: giturl,
	})
	if err != nil {
		log.Fatal("Error: ", err)
	}

	// ... retrieves the branch pointed by HEAD
	ref, err := r.Head()
	if err != nil {
		log.Fatal("Error: ", err)
	}

	// ... retrieves the commit history
	cIter, err := r.Log(&git.LogOptions{From: ref.Hash()})

	// ... just iterates over the commits, printing it
	err = cIter.ForEach(func(commit *object.Commit) error {
		parent, err:=commit.Parent(0)
		if err==nil{
			// log.Info(commit.Hash, ":",parent.Hash)
			patch, _ :=commit.Patch(parent)
			// files:=patch.FilePatches()[0].Files()
			scanner := bufio.NewScanner(strings.NewReader(patch.String()))
			for scanner.Scan() {
				line := scanner.Text()
				found := c.RegexLine(line)
				// dumb
				if len(found) > 0 {
					// log.Debug("Found:",found)
					for _, match := range found{
						// match.URL=result.FileURL
						// match.SearchResult=result
						respChan<-match
					}
				}
			}
			//log.Info(patch)
		}
		return nil
	})
	if err != nil {
		log.Fatal("Error: ", err)
	}
	log.Info(giturl)
	return nil
}
