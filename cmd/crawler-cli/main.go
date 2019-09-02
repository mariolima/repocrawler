package main

import (
  log "github.com/sirupsen/logrus"

  "github.com/pkg/profile"							//profiler for debugging performance

  "strings"
  "fmt"
  "github.com/mariolima/repocrawl/pkg/crawler"
  "github.com/mariolima/repocrawl/cmd/utils"		// used to Highlight matches with colors
)

var (
	GITHUB_ACCESS_TOKEN string
	SLACK_WEBHOOK		string
)

func main() {
	//Debug
	defer profile.Start().Stop()

	level , err := log.ParseLevel(getEnv("LOG_LEVEL","info"))
	if err != nil {
		level = log.InfoLevel
	}
	log.SetLevel(level)


	cmd_opts, err := ParseOptions()
	if err != nil {
		log.Fatal("Cmd options Error: ", err)
		return
	}

	log.WithFields(log.Fields{
		"opts": cmd_opts,
	}).Info("Got Opts:")

	if GITHUB_ACCESS_TOKEN = getEnv("GITHUB_ACCESS_TOKEN",""); GITHUB_ACCESS_TOKEN==""{
		log.Fatal("Please 'export GITHUB_ACCESS_TOKEN' first")
		return
	}

	SLACK_WEBHOOK := getEnv("SLACK_WEBHOOK","");
	if GITHUB_ACCESS_TOKEN==""{
		log.Error("Could not import SLACK_WEBHOOK - not doing notifications")
		return
	}

	repoCrawler, err := crawler.NewRepoCrawler(
		crawler.CrawlerOpts{
			GithubToken: GITHUB_ACCESS_TOKEN,
			BitbucketHost: *cmd_opts.BitbucketHost,
			RulesFile: *cmd_opts.RulesFile,
			SlackWebhook: SLACK_WEBHOOK,
		},
	)
	if err != nil {
		log.Fatal("Failed creating new Crawler: ", err)
		return
	}

	// Channel for Matches found
	matches := make(chan crawler.Match)

	if *cmd_opts.GitUrl != "" {
		go repoCrawler.DeepCrawl(*cmd_opts.GitUrl, matches)
	}

	if *cmd_opts.GithubSearchQuery != "" {
		go repoCrawler.GithubCodeSearch(*cmd_opts.GithubSearchQuery, matches)
	}

	if *cmd_opts.GithubRepo != "" && strings.Contains(*cmd_opts.GithubRepo,"/") {
		repo:=strings.Split(*cmd_opts.GithubRepo,"/")
		go repoCrawler.DeepCrawlGithubRepo(repo[0], repo[1], matches)
	}

	if *cmd_opts.GithubUser != "" {
		go repoCrawler.DeepCrawlGithubUser(*cmd_opts.GithubUser,matches)
	}

	if *cmd_opts.GithubOrg != "" {
		go repoCrawler.DeepCrawlGithubOrg(*cmd_opts.GithubOrg,matches)
	}


	if *cmd_opts.BitbucketUser != "" {
		go repoCrawler.DeepCrawlBitbucketUser(*cmd_opts.BitbucketUser, matches)
	}

	if *cmd_opts.BitbucketRepo != "" {
		go repoCrawler.DeepCrawlBitbucketUser(*cmd_opts.BitbucketRepo, matches)
	}

	for{
		select{
		case match:=<-matches:
			fmt.Printf("%-30s %-90s %s\n",match.Rule.Regex,utils.HighlightWords(match.Line, match.Values),match.URL)
			if match.Rule.Type == "critical" {
				repoCrawler.Notify(match)
			}
		}
	}
}
