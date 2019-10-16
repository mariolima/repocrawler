package crawler

import (
	"github.com/mariolima/repocrawler/internal/entities"
	log "github.com/sirupsen/logrus"
)

// Task used to deepcrawl a set of Users/Git urls
type Task struct {
	respChan chan Match
	Graph    ItemGraph
	State    TaskState
	*crawler
}

// TaskState contains data defining current Task state
type TaskState struct {
	AnalysedRepos []entities.Repository
	AnalysedUsers []entities.User
	Crawling      []entities.Repository
}

func (c *crawler) NewTask(respChan chan Match, root string) *Task {
	t := &Task{
		respChan: respChan,
		crawler:  c,
	}
	c.Tasks = append(c.Tasks, t)
	return t
}

// TaskDone marks current Task as done and removes it from the Task list
func (ct *Task) TaskDone() {

}

// AddRepo adds new Repository to the Task
func (ct *Task) AddRepo(repo entities.Repository) {
	ct.reposChan <- repo
	ct.reposWg.Add(1)
	ct.State.Crawling = append(ct.State.Crawling, repo)
	ct.PushState() //Broadcast the state change
	log.Warn("Crawling ", repo.Name)
}

// AddUser adds new User to the Task
func (ct *Task) AddUser(user entities.User) {
	ct.usersChan <- user
	ct.usersWg.Add(1)
	log.Warn("Crawling user ", user.Name)
}

// DoneRepo marks 'repo` as completed in the Task
func (ct *Task) DoneRepo(repo entities.Repository) {
	<-ct.reposChan
	ct.State.AnalysedRepos = append(ct.State.AnalysedRepos, repo)
	// Shitty way to remove repo from `Currently Crawling`
	for i, r := range ct.State.Crawling {
		if r == repo {
			if len(ct.State.Crawling) > 0 {
				ct.State.Crawling = append(ct.State.Crawling[:i], ct.State.Crawling[i+1:]...)
			} else {
				ct.State.Crawling = nil
			}
		}
	}
	ct.PushState() //Broadcast the state change
	log.Warn("DONE Crawling ", repo.Name)
}

// DoneUser marks 'user` as completed in the Task
func (ct *Task) DoneUser(user entities.User) {
	<-ct.usersChan
	ct.State.AnalysedUsers = append(ct.State.AnalysedUsers, user)
	log.Warn("DONE Crawling user ", user.Name)
}

// Done Checks if current CrawlerTask is done
func (ct *Task) Done() bool {
	return true
}

// WaitRepos waits for all the Repos present in the Task to be completed
func (ct *Task) WaitRepos() {
	ct.reposWg.Wait()
}

// WaitUsers waits for all the Users present in the Task to be completed
func (ct *Task) WaitUsers() {
	ct.usersWg.Wait()
}
