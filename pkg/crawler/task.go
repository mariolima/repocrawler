package crawler

import (
	"github.com/mariolima/repocrawl/internal/entities"
	log "github.com/sirupsen/logrus"
)

/*
	TODO task system
*/
type CrawlerTask struct {
	//TODO move these to crawler struct instead
	respChan      chan Match
	AnalysedRepos []entities.Repository
	AnalysedUsers []entities.User
	*crawler
}

// //
// func (ct *CrawlerTask) DeepCrawl(giturl string) error {
// 	//does the crawling
// 	//blablabla
// 	ct.responseChan <- nil
// 	return nil
// }

// func (c *crawler) DeepCrawl(giturl string, respChan chan Match) error {
// 	// setup goroutines with c.Opts (nthreads)
// 	// adds task to the list of Tasks in Crawler
// 	c.AddTask(&CrawlerTask{
// 		responseChan: respChan,
// 	})
// 	return nil
// }

func (c *crawler) NewCrawlerTask(respChan chan Match) CrawlerTask {
	return CrawlerTask{
		// c.UsersChan: make(chan entities.Repository, c.Opts.NrThreads),
		// usersChan:   make(chan entities.User, c.Opts.NrThreads),
		respChan: respChan,
		crawler:  c,
	}
}

// AddRepo adds new Repository to the CrawlerTask
func (ct *CrawlerTask) AddRepo(repo entities.Repository) {
	ct.reposChan <- repo
	ct.reposWg.Add(1)
	log.Warn("Crawling ", repo.Name)
}

// AddUser adds new User to the CrawlerTask
func (ct *CrawlerTask) AddUser(user entities.User) {
	ct.usersChan <- user
	ct.usersWg.Add(1)
	log.Warn("Crawling user ", user.Name)
}

// DoneRepo marks 'repo` as completed in the CrawlerTask
func (ct *CrawlerTask) DoneRepo(repo entities.Repository) {
	<-ct.reposChan
	ct.AnalysedRepos = append(ct.AnalysedRepos, repo)
	log.Warn("DONE Crawling ", repo.Name)
}

// DoneUser marks 'user` as completed in the CrawlerTask
func (ct *CrawlerTask) DoneUser(user entities.User) {
	<-ct.usersChan
	ct.AnalysedUsers = append(ct.AnalysedUsers, user)
	log.Warn("DONE Crawling user ", user.Name)
}

// WaitRepos waits for all the Repos present in the CrawlerTask to be completed
func (ct *CrawlerTask) WaitRepos() {
	ct.reposWg.Wait()
}

// WaitUsers waits for all the Users present in the CrawlerTask to be completed
func (ct *CrawlerTask) WaitUsers() {
	ct.usersWg.Wait()
}
