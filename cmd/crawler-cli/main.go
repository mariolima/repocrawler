package main

import (
	log "github.com/sirupsen/logrus"

	"github.com/pkg/profile" //profiler for debugging performance

	"fmt"
	"strings"

	"github.com/mariolima/repocrawler/cmd/utils" // used to Highlight matches with colors
	"github.com/mariolima/repocrawler/cmd/utils/webserver"
	"github.com/mariolima/repocrawler/pkg/crawler"
)

var (
	GITHUB_ACCESS_TOKEN string
	SLACK_WEBHOOK       string
)

func main() {
	//Debug
	defer profile.Start().Stop()

	level, err := log.ParseLevel(getEnv("LOG_LEVEL", "info"))
	if err != nil {
		level = log.InfoLevel
	}
	log.SetLevel(level)

	cmdOpts, err := ParseOptions()
	if err != nil {
		log.Fatal("Cmd options Error: ", err)
		return
	}

	log.WithFields(log.Fields{
		"opts": cmdOpts,
	}).Info("Got Opts:")

	if GITHUB_ACCESS_TOKEN = getEnv("GITHUB_ACCESS_TOKEN", ""); GITHUB_ACCESS_TOKEN == "" {
		log.Fatal("Please 'export GITHUB_ACCESS_TOKEN' first")
		return
	}

	SLACK_WEBHOOK := getEnv("SLACK_WEBHOOK", "")
	if SLACK_WEBHOOK == "" {
		log.Error("SLACK_WEBHOOK not in Env - not doing notifications")
	}

	repoCrawler, err := crawler.NewRepoCrawler(
		crawler.Opts{
			NrThreads:     *cmdOpts.NrThreads,
			GithubToken:   GITHUB_ACCESS_TOKEN,
			BitbucketHost: *cmdOpts.BitbucketHost,
			RulesFile:     *cmdOpts.RulesFile,
			SlackWebhook:  SLACK_WEBHOOK,
		},
	)
	if err != nil {
		log.Fatal("Failed creating new Crawler: ", err)
		return
	}

	// repoCrawler.TestGraph()
	// return

	if *cmdOpts.WebServer {
		repoCrawler.AddMatchServer(&webserver.MatchServer{
			Port:     8090,
			Hostname: "gobh",
			CertFile: "configs/certs/",
		})
	}

	// Channel for Matches found
	matches := make(chan crawler.Match)

	if *cmdOpts.GitURL != "" {
		go repoCrawler.DeepCrawl(*cmdOpts.GitURL, matches)
	}

	if *cmdOpts.GithubSearchQuery != "" {
		go repoCrawler.GithubCodeSearch(*cmdOpts.GithubSearchQuery, matches)
	}

	if *cmdOpts.GithubRepo != "" && strings.Contains(*cmdOpts.GithubRepo, "/") {
		repo := strings.Split(*cmdOpts.GithubRepo, "/")
		go repoCrawler.DeepCrawlGithubRepo(repo[0], repo[1], matches)
	}

	if *cmdOpts.GithubUser != "" {
		go repoCrawler.DeepCrawlGithubUser(*cmdOpts.GithubUser, matches)
	}

	if *cmdOpts.GithubOrg != "" {
		go repoCrawler.DeepCrawlGithubOrg(*cmdOpts.GithubOrg, matches)
	}

	if *cmdOpts.BitbucketUser != "" {
		go repoCrawler.DeepCrawlBitbucketUser(*cmdOpts.BitbucketUser, matches)
	}

	if *cmdOpts.BitbucketRepo != "" {
		go repoCrawler.DeepCrawlBitbucketUser(*cmdOpts.BitbucketRepo, matches)
	}

	for match := range matches {
		matchLine := utils.HighlightWords(utils.TruncateString(match.Line, match.Values, 20, 500), match.Values)
		line := fmt.Sprintf("[%f] %-30s %-90s %s\n", match.Entropy, match.Rule.Regex, matchLine, match.URL)
		fmt.Print(line)
		if match.Rule.Type == "critical" {
			repoCrawler.Notify(match)
		}
		utils.SaveLineToFile(line, *cmdOpts.OutputFile)
	}
}
