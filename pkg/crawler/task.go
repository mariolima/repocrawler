package crawler

import (
	"github.com/mariolima/repocrawl/internal/entities"
	log "github.com/sirupsen/logrus"
	"sync"
)

/*
	TODO task system
*/
type CrawlerTask struct {
	repos_wg  sync.WaitGroup
	reposChan chan entities.Repository //control maximum number of concurrent repos being crawled
	usersChan chan entities.User       //control maximum number of concurrent users being crawled
	users_wg  sync.WaitGroup
	respChan  chan Match
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
		reposChan: make(chan entities.Repository, c.Opts.NrThreads),
		usersChan: make(chan entities.User, c.Opts.NrThreads),
		respChan:  respChan,
		crawler:   c,
	}
}

func (ct *CrawlerTask) AddRepo(repo entities.Repository) {
	ct.reposChan <- repo
	ct.repos_wg.Add(1)
	log.Warn("Crawling ", repo.Name)
}

func (ct *CrawlerTask) AddUser(user entities.User) {
	ct.usersChan <- user
	ct.users_wg.Add(1)
	log.Warn("Crawling user ", user.Name)
}

func (ct *CrawlerTask) DoneRepo(repo entities.Repository) {
	<-ct.reposChan
	log.Warn("DONE Crawling ", repo.Name)
}

func (ct *CrawlerTask) DoneUser(repo entities.User) {
	<-ct.usersChan
	log.Warn("DONE Crawling user ", repo.Name)
}

func (ct *CrawlerTask) WaitRepos() {
	ct.repos_wg.Wait()
}

func (ct *CrawlerTask) WaitUsers() {
	ct.users_wg.Wait()
}
