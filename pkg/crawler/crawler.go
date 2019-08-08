package crawler

import(
	"github.com/mariolima/repocrawl/pkg/github"
	"github.com/mariolima/repocrawl/pkg/bitbucket"
	"github.com/mariolima/repocrawl/pkg/gitlab"
)

type crawler struct{
	CrawlerOpts
	GithubCrawler	*github.GitHubCrawler
}

type CrawlerOpts struct{
	GITHUB_ACCESS_TOKEN		string
}

type Crawler interface{

}

func NewRepoCrawler(opts CrawlerOpts) *crawler {
	return &crawler{ 
		opts
	}
}

func (c *crawler) GithubCodeSearch(q string) (error) {
	//TODO
}
