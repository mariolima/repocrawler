package crawler

import (
	"fmt"
	"github.com/mariolima/repocrawl/cmd/utils"
	// log "github.com/sirupsen/logrus"
)

func (c *crawler) Notify(match Match) {
	// Notify through Slack Webhook
	if c.Opts.SlackWebhook != "" {
		line := fmt.Sprintf("*%s* `%s` - %s\n", match.Rule.Regex, match.Values[0], match.URL)
		utils.SlackNotify(line, c.Opts.SlackWebhook)
	}
}
