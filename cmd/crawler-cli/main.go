package main

import (
	log "github.com/sirupsen/logrus"

	"github.com/pkg/profile" //profiler for debugging performance

	"fmt"
	"github.com/mariolima/repocrawl/cmd/utils" // used to Highlight matches with colors
	"github.com/mariolima/repocrawl/cmd/utils/webserver"
	"github.com/mariolima/repocrawl/pkg/crawler"
	"strings"
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

	cmd_opts, err := ParseOptions()
	if err != nil {
		log.Fatal("Cmd options Error: ", err)
		return
	}

	log.WithFields(log.Fields{
		"opts": cmd_opts,
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
		crawler.CrawlerOpts{
			NrThreads:     *cmd_opts.NrThreads,
			GithubToken:   GITHUB_ACCESS_TOKEN,
			BitbucketHost: *cmd_opts.BitbucketHost,
			RulesFile:     *cmd_opts.RulesFile,
			SlackWebhook:  SLACK_WEBHOOK,
		},
	)
	if err != nil {
		log.Fatal("Failed creating new Crawler: ", err)
		return
	}

	// repoCrawler.TestGraph()
	// return

	if *cmd_opts.WebServer {
		repoCrawler.AddMatchServer(&webserver.MatchServer{
			Port:     8090,
			Hostname: "gobh",
			CertFile: "configs/certs/",
		})
	}

	// Channel for Matches found
	matches := make(chan crawler.Match)

	if *cmd_opts.GitUrl != "" {
		go repoCrawler.DeepCrawl(*cmd_opts.GitUrl, matches)
	}

	if *cmd_opts.GithubSearchQuery != "" {
		go repoCrawler.GithubCodeSearch(*cmd_opts.GithubSearchQuery, matches)
	}

	if *cmd_opts.GithubRepo != "" && strings.Contains(*cmd_opts.GithubRepo, "/") {
		repo := strings.Split(*cmd_opts.GithubRepo, "/")
		go repoCrawler.DeepCrawlGithubRepo(repo[0], repo[1], matches)
	}

	if *cmd_opts.GithubUser != "" {
		go repoCrawler.DeepCrawlGithubUser(*cmd_opts.GithubUser, matches)
	}

	if *cmd_opts.GithubOrg != "" {
		go repoCrawler.DeepCrawlGithubOrg(*cmd_opts.GithubOrg, matches)
	}

	if *cmd_opts.BitbucketUser != "" {
		go repoCrawler.DeepCrawlBitbucketUser(*cmd_opts.BitbucketUser, matches)
	}

	if *cmd_opts.BitbucketRepo != "" {
		go repoCrawler.DeepCrawlBitbucketUser(*cmd_opts.BitbucketRepo, matches)
	}

	for match := range matches {
		match_line := utils.HighlightWords(utils.TruncateString(match.Line, match.Values, 20, 500), match.Values)
		line := fmt.Sprintf("[%f] %-30s %-90s %s\n", match.Entropy, match.Rule.Regex, match_line, match.URL)
		fmt.Print(line)
		if match.Rule.Type == "critical" {
			repoCrawler.Notify(match)
		}
		utils.SaveLineToFile(line, *cmd_opts.OutputFile)
	}
}
