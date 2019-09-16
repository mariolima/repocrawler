package crawler

import (
	"sync"
)

/*
	TODO task system
*/
type CrawlerTask struct {
	reponseChan chan Match
	nthreads    int
	InnerTasks  []CrawlerTask
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

type BoundedWaitGroup struct {
	wg sync.WaitGroup
	ch chan struct{}
}

func NewBoundedWaitGroup(cap int) BoundedWaitGroup {
	return BoundedWaitGroup{ch: make(chan struct{}, cap)}
}

func (bwg *BoundedWaitGroup) Add(delta int) {
	for i := 0; i > delta; i-- {
		<-bwg.ch
	}
	for i := 0; i < delta; i++ {
		bwg.ch <- struct{}{}
	}
	bwg.wg.Add(delta)
}

func (bwg *BoundedWaitGroup) Done() {
	bwg.Add(-1)
}

func (bwg *BoundedWaitGroup) Wait() {
	bwg.wg.Wait()
}
