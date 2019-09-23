package main

import (
	"testing"

	log "github.com/sirupsen/logrus"

	"github.com/pkg/profile" //profiler for debugging performance

	"fmt"

	"github.com/mariolima/repocrawl/cmd/utils" // used to Highlight matches with colors
	"github.com/mariolima/repocrawl/cmd/utils/webserver"
	"github.com/mariolima/repocrawl/pkg/crawler"
	"github.com/stretchr/testify/assert"
)

func TestMain(t *testing.T) {
	//Debug
	defer profile.Start().Stop()

	log.SetLevel(log.DebugLevel)

	gitURL := "https://github.com/fsubal/TwitKJ_New"

	if GITHUB_ACCESS_TOKEN = getEnv("GITHUB_ACCESS_TOKEN", ""); GITHUB_ACCESS_TOKEN == "" {
		log.Fatal("Please 'export GITHUB_ACCESS_TOKEN' first")
		return
	}

	repoCrawler, err := crawler.NewRepoCrawler(
		crawler.Opts{
			NrThreads:     5,
			GithubToken:   GITHUB_ACCESS_TOKEN,
			BitbucketHost: "",
			RulesFile:     "rules.json",
			SlackWebhook:  "",
		},
	)
	if err != nil {
		log.Fatal("Failed creating new Crawler: ", err)
		return
	}

	repoCrawler.AddMatchServer(&webserver.MatchServer{
		Port:     8090,
		Hostname: "gobh",
		CertFile: "configs/certs/",
	})

	// Channel for Matches found
	matches := make(chan crawler.Match)
	go repoCrawler.DeepCrawl(gitURL, matches)

	for match := range matches {
		matchLine := utils.HighlightWords(utils.TruncateString(match.Line, match.Values, 20, 500), match.Values)
		line := fmt.Sprintf("[%f] %-30s %-90s %s\n", match.Entropy, match.Rule.Regex, matchLine, match.URL)
		fmt.Print(line)
		if match.Rule.Type == "critical" {
			repoCrawler.Notify(match)
		}
		fmt.Println(match.Line)
		assert.Equal(t, match.Line, `define("CONSUMER_SECRET","lEFZGpHMI1OpVVHH02mJQCIsvnMfjikVr7L3l2TnBo");`)
		return
		// utils.SaveLineToFile(line, *cmd_opts.OutputFile)
	}
}
