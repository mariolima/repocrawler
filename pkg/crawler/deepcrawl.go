package crawler

import(
	log "github.com/sirupsen/logrus"
)


func (c *crawler) DeepCrawl(giturl string) (error) {
	log.Info(giturl)
	return nil
}
