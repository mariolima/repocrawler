package main

import (
  log "github.com/sirupsen/logrus"

  "github.com/pkg/profile"			  //profiler for debugging

  "fmt"
  "github.com/mariolima/repocrawl/pkg/crawler"
)

var GITHUB_ACCESS_TOKEN string

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

	repoCrawler, err := crawler.NewRepoCrawler(
		crawler.CrawlerOpts{
			GithubToken: GITHUB_ACCESS_TOKEN,
			BitbucketHost: *cmd_opts.BitbucketHost,
			RulesFile: *cmd_opts.RulesFile,
		},
	)
	if err != nil {
		log.Fatal("Failed creating new Crawler: ", err)
		return
	}

	// Channel for Matches found
	matches := make(chan crawler.Match)

	// go repoCrawler.DeepCrawlBitbucketUser("", matches)
	// go repoCrawler.DeepCrawlBitbucketRepo("openncp","tsl-utils", matches)

	// go repoCrawler.DeepCrawl("https://bitbucket.org/atlassian/serverless-deploy", matches)
	// go repoCrawler.DeepCrawlGithubRepo("khypponen", "openncp", matches)
	// go repoCrawler.DeepCrawlGithubOrg("TwilioDevEd",matches)
	// go repoCrawler.DeepCrawlGithubUser("",matches)

	go repoCrawler.GithubCodeSearch(*cmd_opts.GithubSearchQuery, matches)

	for{
		select{
		case match:=<-matches:
			fmt.Printf("%-30s %-90s %s\n",match.Rule.Regex,match.Line,match.URL)
		}
	}
}
