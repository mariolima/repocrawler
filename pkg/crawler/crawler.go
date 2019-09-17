package crawler

import (
	"github.com/mariolima/repocrawl/internal/entities"
	"github.com/mariolima/repocrawl/pkg/bitbucket"
	"github.com/mariolima/repocrawl/pkg/github"
	_ "github.com/mariolima/repocrawl/pkg/gitlab" //TODO

	"bufio"
	log "github.com/sirupsen/logrus"
	"regexp"
	"strings"

	"fmt" //TODO used just to format github lines #L01

	"encoding/json" //to get the rules file
	"io/ioutil"
	"os"
)

type MatchServer interface {
	Setup()
	PushMatch(Match) error        //pushes
	PushLogEntry(log.Entry) error //pushes logrus log entry
}

type crawler struct {
	Github       *github.GitHubCrawler
	Bitbucket    *bitbucket.BitbucketCrawler
	MatchRules   map[string]map[string]*regexp.Regexp
	MatchServers []MatchServer
	Opts         CrawlerOpts
}

type CrawlerOpts struct {
	GithubToken   string `json:"github_token,omitempty"`
	BitbucketHost string `json:"bitbucket_host,omitempty"`
	RulesFile     string `json:"rules_file,omitempty"`
	SlackWebhook  string `json:"webhook,omitempty"`
	NrThreads     int    `json:"nthreads,omitempty"`
}

// type Task struct{
// 	MatchesChannel	chan Match
// }

func NewRepoCrawler(opts CrawlerOpts) (*crawler, error) {
	//TODO depending on the opts - create only the required Clients
	c := crawler{
		Github:    github.NewCrawler(opts.GithubToken),
		Bitbucket: bitbucket.NewCrawler(opts.BitbucketHost, "test", "tesT"),
		Opts:      opts,
	}
	log.Info(c.Opts.RulesFile)
	err := c.compileRegexes() //Pre compile regexes for better performance
	if err != nil {
		log.Error("Couldn't load Regexes")
		return nil, err
	}
	log.Debug("Compiled Regexes Successfully")

	log.AddHook(&Logger{&c})

	return &c, nil
}

// HOOK for Logrus
type Logger struct {
	*crawler
}

func (x Logger) Fire(entry *log.Entry) error {
	for _, ms := range x.MatchServers {
		ms.PushLogEntry(*entry)
	}
	return nil //TODO
}

func (x Logger) Levels() []log.Level {
	return log.AllLevels //TODO remove Trace
}

// ---

func (c *crawler) AddMatchServer(ms MatchServer) {
	c.MatchServers = append(c.MatchServers, ms)
	log.Warn("Starting MatchServer")
	go ms.Setup()
	// if err != nil {
	// 	log.Fatal("Couldn't Match Server")
	// }
	// log.Warn("Started MatchServer")
}

func (c *crawler) PushMatch(match Match) {
	for _, ms := range c.MatchServers {
		ms.PushMatch(match)
	}
}

func (c *crawler) GithubCodeSearch(query string, response chan Match) {
	//make new task and setup multithreading with c.Opts etc
	var matches = make(map[string]Match)
	searchResultChan := make(chan entities.SearchResult)
	go c.Github.SearchCode(query, searchResultChan)
	for {
		select {
		case result := <-searchResultChan:
			log.Trace("Result received:", result)
			scanner := bufio.NewScanner(strings.NewReader(result.FileContent))
			i := 1
			for scanner.Scan() {
				line := scanner.Text()
				// if line already matched this session
				if _, ok := matches[line]; ok {
					continue
				}
				found := c.RegexLine(line)
				// dumb
				if len(found) > 0 {
					log.Debug("Found:", found)
					for _, match := range found {
						// avoid dups
						matches[line] = match
						match.URL = fmt.Sprintf("%s#L%d", result.FileURL, i)
						// match.URL=result.FileURL
						match.SearchResult = result
						match.LineNr = i
						response <- match
					}
				}
				i += 1
			}
		}
	}
}

func (c *crawler) compileRegexes() error {
	// some of the rules taken from https://github.com/dxa4481/truffleHogRegexes/blob/master/truffleHogRegexes/regexes.json :)
	// but some were changed for better accuracy
	// List of popular API services and their keyworks from https://github.com/streaak/keyhacks
	rules := map[string]map[string]string{
		"keywords": {
			// "secret":"(?i)(secret)\\W",
			// "key":"(?i)(key)\\W",
			// "token":"(?i)(token)",
			// "password":"(?i)(password)",

			// "jdbc": "(?i)(jdbc)",
		},
		"keys": {
			// "GitHub": "[g|G][i|I][t|T][h|H][u|U][b|B].*['|\"][0-9a-zA-Z]{35,40}['|\"]", //giving me issues
			//TODO FIX - Google API Keys not working properly
			// "Google Cloud Platform API Key": "AIza[0-9A-Za-z\\-_]{35}",
			// "Google Cloud Platform OAuth": "[0-9]+-[0-9A-Za-z_]{32}\\.apps\\.googleusercontent\\.com",
			// "Google Drive API Key": "AIza[0-9A-Za-z\\-_]{35}",
			// "Google Drive OAuth": "[0-9]+-[0-9A-Za-z_]{32}\\.apps\\.googleusercontent\\.com",
			// "Google (GCP) Service-account": "\"type\": \"service_account\"",
			// "Google Gmail API Key": "AIza[0-9A-Za-z\\-_]{35}",
			// "Google Gmail OAuth": "[0-9]+-[0-9A-Za-z_]{32}\\.apps\\.googleusercontent\\.com",
			// "Google OAuth Access Token": "ya29\\.[0-9A-Za-z\\-_]+",
			// "Google YouTube API Key": "AIza[0-9A-Za-z\\-_]{35}",
			// "Google YouTube OAuth": "[0-9]+-[0-9A-Za-z_]{32}\\.apps\\.googleusercontent\\.com",
			/*
				MY RULES ^_^
			*/
			// "Generic Key/Secret": "(?i)(key)\\s{0,1}[=:]+\\s{0,1}['\"].{1,}['\"]+\\W",		//flawed /w big oneliners - key:"down"asdasdasdas" fix
			"Generic Secret":     "[s|S][e|E][c|C][r|R][e|E][t|T].*['|\"][0-9a-zA-Z]{32,45}['|\"]",
			"Generic Key/Secret": "(?i)(key|secret)\\s{0,1}[=:]+\\s{0,1}['\"][^\"']{5,}['\"]",
			// "Generic Env Key/Secret": "(?i)(key|secret)\\s{0,1}[=:]+\\s{0,1}[^'\"\\n]{4,}",
		},
		"critical": {
			// "Hardcoded Env Password": "(?i)(password)\\s{0,1}[=:]+\\s{0,1}[^'\"\\n]{4,}",
			"Slack Token":              "(xox[p|b|o|a]-[0-9]{12}-[0-9]{12}-[0-9]{12}-[a-z0-9]{32})",
			"Amazon AWS Access Key ID": "AKIA[0-9A-Z]{16}",
			"Amazon MWS Auth Token":    "amzn\\.mws\\.[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}",
			"AWS API Key":              "AKIA[0-9A-Z]{16}",
			"Facebook Access Token":    "EAACEdEose0cBA[0-9A-Za-z]+",
			"Facebook OAuth":           "[f|F][a|A][c|C][e|E][b|B][o|O][o|O][k|K].*['|\"][0-9a-f]{32}['|\"]",
			"priv_keys":                "(?s)(-----BEGIN (RSA|DSA|PGP|EC|) PRIVATE KEY.*END (RSA|DSA|PGP|EC|) PRIVATE KEY-----)",
			"RSA private key":          "-----BEGIN RSA PRIVATE KEY-----",
			"SSH (DSA) private key":    "-----BEGIN DSA PRIVATE KEY-----",
			"SSH (EC) private key":     "-----BEGIN EC PRIVATE KEY-----",
			"PGP private key block":    "-----BEGIN PGP PRIVATE KEY BLOCK-----",
			// "Auth": "\\W(Authorization:).+\\W",
			"Google API Key":                "AIza[0-9A-Za-z\\-_]{35}",
			"Heroku API Key":                "[h|H][e|E][r|R][o|O][k|K][u|U].*[0-9A-F]{8}-[0-9A-F]{4}-[0-9A-F]{4}-[0-9A-F]{4}-[0-9A-F]{12}",
			"MailChimp API Key":             "[0-9a-f]{32}-us[0-9]{1,2}",
			"Mailgun API Key":               "key-[0-9a-zA-Z]{32}",
			"Password in URL":               "[a-zA-Z]{3,10}://[^/\\s:@]{3,20}:[^/\\s:@]{3,20}@.{1,100}[\"'\\s]",
			"PayPal Braintree Access Token": "access_token\\$production\\$[0-9a-z]{16}\\$[0-9a-f]{32}",
			"Picatic API Key":               "sk_live_[0-9a-z]{32}",
			"Slack Webhook":                 "https://hooks.slack.com/services/T[a-zA-Z0-9_]{8}/B[a-zA-Z0-9_]{8}/[a-zA-Z0-9_]{24}",
			"Stripe API Key":                "sk_live_[0-9a-zA-Z]{24}",
			"Stripe Restricted API Key":     "rk_live_[0-9a-zA-Z]{24}",
			"Square Access Token":           "sq0atp-[0-9A-Za-z\\-_]{22}",
			"Square OAuth Secret":           "sq0csp-[0-9A-Za-z\\-_]{43}",
			"Twilio API Key":                "SK[0-9a-fA-F]{32}",
			"Twitter Access Token":          "[t|T][w|W][i|I][t|T][t|T][e|E][r|R].*[1-9][0-9]+-[0-9a-zA-Z]{40}",
			"Twitter OAuth":                 "[t|T][w|W][i|I][t|T][t|T][e|E][r|R].*['|\"][0-9a-zA-Z]{35,44}['|\"]",
			"Fastly API Key":                "\\W(Fastly-key)\\W*[A-Za-z0-9+=]{44,}\\W", //is it b64 tho?
			"Disqus API Key":                "\\W(?i)(disqus).+\\w[K|k][E|e][Y|y]\\W+[A-Za-z0-9]{64}\\W",
			"Zoho Desk Token":               "[0-9]{4}[.]+[0-9a-f]{32}[.]+[0-9a-f]{32}", //https://desk.zoho.com/DeskAPIDocument
			"Auth Bearer":                   "(Authorization: Bearer )[a-zA-z0-9]{20,}\\S",
			"Auth Basic":                    "(Basic )[a-zA-z0-9+=]{20,}\\S",
			"Generic Hash Token":            "\\w+[TtOoKkEeNn]\\s{0,1}[=:]\\s{0,1}([0-9a-f]{32}|[0-9a-f]{40}|[0-9a-f]{56}|[0-9a-f]{64})\\W", //md5, sha1, sha224, sha256
			"Generic Hash Key":              "\\w+[KkEeYy]\\s{0,1}[=:]\\s{0,1}([0-9a-f]{32}|[0-9a-f]{40}|[0-9a-f]{56}|[0-9a-f]{64})\\W",     //md5, sha1, sha224, sha256
			// "Generic UUID":                  "[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}\\W",
			"API Key":        "\\S*[K|k][E|e][Y|y]+\\W+[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}",
			"Meraki API Key": "[X|x]-[C|c][I|i][S|s][C|c][O|o]+-[M|m][E|e][R|r][A|a][K|k][I|i].+[:=]\\W+[0-9a-f]{40}\\W",
			"Mapbox Api Key": "['\"](sk.)[0-9A-za-z]{58,}(.)[0-9A-za-z]{22}['\"]",
			//TODO FIX these / prob caps
			"Hardcoded Password":        "(?i)(password)\\s{0,1}[=:]+\\s{0,1}['\"][^\"']{4,}['\"]+\\W", //Slightly better regex for passwords
			"Trello API Key":            "(?i)['\"]{0,1}(trello_api_key)['\"]{0,1}\\s{0,1}[:=]\\s{0,1}['\"]{0,1}[0-9a-f]{32}['\"]{0,1}\\W",
			"Algolia API Key":           "(?i)['\"]{0,1}(x-algolia-api-key)['\"]{0,1}\\s{0,1}[:=]\\s{0,1}['\"]{0,1}[0-9a-f]{32}['\"]{0,1}\\W",
			"Pendo Integration API Key": "(?i)['\"]{0,1}(x-pendo-integration-key)['\"]{0,1}\\s{0,1}[:=]\\s{0,1}['\"]{0,1}[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}['\"]{0,1}\\W",
		},
	}

	// save regexes as a test - TODO delete later
	file, _ := json.MarshalIndent(rules, "", " ")
	_ = ioutil.WriteFile("rules.json", file, 0644)

	// read rules File
	jsonFile, err := os.Open(c.Opts.RulesFile)
	defer jsonFile.Close()
	if err != nil {
		log.Error("Couldn't Open regexes file ", c.Opts.RulesFile)
		return err
	}
	byteValue, _ := ioutil.ReadAll(jsonFile)
	json.Unmarshal(byteValue, &rules)

	//Retarded code - fix later ?
	c.MatchRules = map[string]map[string]*regexp.Regexp{}
	for rule_type, regexes := range rules {
		c.MatchRules[rule_type] = make(map[string]*regexp.Regexp)
		for rule, regex := range regexes {
			c.MatchRules[rule_type][rule] = regexp.MustCompile(regex)
			// log.Trace(c.MatchRules[rule_type][rule])
		}
	}
	return nil
}

func (c *crawler) RegexLine(line string) (matches []Match) {
	results := make(map[string]Match) // map of matches in order to avoid duplicates with the same Rules
	// prioratizes later Regexes (more critical)
	for rule_type, regexes := range c.MatchRules {
		for rule, re := range regexes {
			// matched, err := regexp.Match(regex, []byte(line))
			ms := re.FindAllString(line, -1) //https://golang.org/pkg/regexp/#Regexp.FindAllString
			if len(ms) > 0 {
				if len(ms[0]) > 500 {
					continue
				}
				result := line
				results[line] = Match{
					Rule:    MatchRule{rule_type, rule},
					Entropy: FindEntropy(ms[0]), //TODO setup LineEntropy and Values Entropy
					Values:  ms,
					Line:    strings.TrimSpace(result), //TODO trim length
				}
				// matches=append(matches,Match{
				// 	Rule: MatchRule{ rule_type, rule },
				// 	Value: ms[0],
				// 	Line: result,
				// })
			}
		}
	}
	for _, m := range results {
		matches = append(matches, m)
	}
	return matches
}

func (c *crawler) PushLog(msg string) {

}
