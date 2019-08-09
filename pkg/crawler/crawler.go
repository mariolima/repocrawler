package crawler

import(
	"github.com/mariolima/repocrawl/pkg/github"
	_"github.com/mariolima/repocrawl/pkg/bitbucket"
	_"github.com/mariolima/repocrawl/pkg/gitlab"
	"github.com/mariolima/repocrawl/internal/entities"

	log "github.com/sirupsen/logrus"
	"regexp"
	"bufio"
	"strings"
	"fmt"
	"github.com/logrusorgru/aurora" //colors - why this? because it is simple to replace text with colored text (others are not)

)

type crawler struct{
	Github				*github.GitHubCrawler
	MatchRules			[]MatchRule
	Opts				CrawlerOpts
}

type CrawlerOpts struct{
	GITHUB_ACCESS_TOKEN		string
}

type Crawler interface{

}

type Task struct{
	MatchesChannel	chan Match
}

func NewRepoCrawler(opts CrawlerOpts) *crawler {
	return &crawler{
		Github: github.NewGitHubCrawler(opts.GITHUB_ACCESS_TOKEN),
		Opts: opts,
	}
}

type Match struct { //Has to be generic - TODO move to other pkg
	Rule		MatchRule
	Line		string
	//Repository struct // User struct and other generic stuff
	URL			string
}

type MatchRule struct {
	Type		string
	Regex		string
}

type CodeText string

func (c *crawler) GithubCodeSearch(query string, response chan Match) {
	//make new task and setup multithreading with c.Opts etc
	searchResultChan := make(chan entities.SearchResult)
	go c.Github.SearchCode(query,searchResultChan)
	for {
		select{
		case result:=<-searchResultChan:
			log.Info("Result received:",result)
			scanner := bufio.NewScanner(strings.NewReader(result.FileContent))
			for scanner.Scan() {
				line := scanner.Text()
				found := c.RegexLine(line)
				// dumb
				if len(found) > 0 {
					log.Info(found)
					for _, match := range found{
						match.URL=result.FileURL
						response<-match
					}
				}
			}
		}
	}
}

// func (c *crawler) CrawlFile(file_content string, response chan Match) {
// 	scanner := bufio.NewScanner(strings.NewReader(file_content))
//
// 	for scanner.Scan() {
// 		line := scanner.Text()
// 		log.Trace("Sending Content:",line)
// 		found := c.RegexLine(line)
// 	}
// }

func (c *crawler) DeepCrawl(git_url string) (error) {
	//TODO
	return nil
}

//TODO move this to internals/utils
func HighlightWords(line string, words []string) (res string) {
	if words == nil {
		return line
	}
	for _, word := range words {
		res=strings.ReplaceAll(line,word,fmt.Sprintf("%s",aurora.Green(word)))
	}
	return res
}

func HighlightWord(line string, word string) string {
	return strings.ReplaceAll(line,word,fmt.Sprintf("%s",aurora.Green(word)))
}

func (c *crawler) RegexLine(line string) (matches []Match) {
	//TODO import rules from cfg
	// rules taken from https://github.com/dxa4481/truffleHogRegexes/blob/master/truffleHogRegexes/regexes.json :)
	rules := map[string]string{
		"secret":"(?i)(secret)\\W",
		"token":"(?i)(token)",
		// "password":"(?i)(password)",

		"jdbc":"(?i)(jdbc)",
		"priv_keys":"(?s)(-----BEGIN (RSA|DSA|PGP|EC|) PRIVATE KEY.*END (RSA|DSA|PGP|EC|) PRIVATE KEY-----)",
		"Slack Token": "(xox[p|b|o|a]-[0-9]{12}-[0-9]{12}-[0-9]{12}-[a-z0-9]{32})",
		"RSA private key": "-----BEGIN RSA PRIVATE KEY-----",
		"SSH (DSA) private key": "-----BEGIN DSA PRIVATE KEY-----",
		"SSH (EC) private key": "-----BEGIN EC PRIVATE KEY-----",
		"PGP private key block": "-----BEGIN PGP PRIVATE KEY BLOCK-----",
		"Amazon AWS Access Key ID": "AKIA[0-9A-Z]{16}",
		"Amazon MWS Auth Token": "amzn\\.mws\\.[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}",
		"AWS API Key": "AKIA[0-9A-Z]{16}",
		"Facebook Access Token": "EAACEdEose0cBA[0-9A-Za-z]+",
		"Facebook OAuth": "[f|F][a|A][c|C][e|E][b|B][o|O][o|O][k|K].*['|\"][0-9a-f]{32}['|\"]",
		// "GitHub": "[g|G][i|I][t|T][h|H][u|U][b|B].*['|\"][0-9a-zA-Z]{35,40}['|\"]", //giving me issues
		"Generic API Key": "[a|A][p|P][i|I][_]?[k|K][e|E][y|Y].*['|\"][0-9a-zA-Z]{32,45}['|\"]",
		"Generic Secret": "[s|S][e|E][c|C][r|R][e|E][t|T].*['|\"][0-9a-zA-Z]{32,45}['|\"]",
		"Google API Key": "AIza[0-9A-Za-z\\-_]{35}",
		"Google Cloud Platform API Key": "AIza[0-9A-Za-z\\-_]{35}",
		"Google Cloud Platform OAuth": "[0-9]+-[0-9A-Za-z_]{32}\\.apps\\.googleusercontent\\.com",
		"Google Drive API Key": "AIza[0-9A-Za-z\\-_]{35}",
		"Google Drive OAuth": "[0-9]+-[0-9A-Za-z_]{32}\\.apps\\.googleusercontent\\.com",
		"Google (GCP) Service-account": "\"type\": \"service_account\"",
		"Google Gmail API Key": "AIza[0-9A-Za-z\\-_]{35}",
		"Google Gmail OAuth": "[0-9]+-[0-9A-Za-z_]{32}\\.apps\\.googleusercontent\\.com",
		"Google OAuth Access Token": "ya29\\.[0-9A-Za-z\\-_]+",
		"Google YouTube API Key": "AIza[0-9A-Za-z\\-_]{35}",
		"Google YouTube OAuth": "[0-9]+-[0-9A-Za-z_]{32}\\.apps\\.googleusercontent\\.com",
		"Heroku API Key": "[h|H][e|E][r|R][o|O][k|K][u|U].*[0-9A-F]{8}-[0-9A-F]{4}-[0-9A-F]{4}-[0-9A-F]{4}-[0-9A-F]{12}",
		"MailChimp API Key": "[0-9a-f]{32}-us[0-9]{1,2}",
		"Mailgun API Key": "key-[0-9a-zA-Z]{32}",
		"Password in URL": "[a-zA-Z]{3,10}://[^/\\s:@]{3,20}:[^/\\s:@]{3,20}@.{1,100}[\"'\\s]",
		"PayPal Braintree Access Token": "access_token\\$production\\$[0-9a-z]{16}\\$[0-9a-f]{32}",
		"Picatic API Key": "sk_live_[0-9a-z]{32}",
		"Slack Webhook": "https://hooks.slack.com/services/T[a-zA-Z0-9_]{8}/B[a-zA-Z0-9_]{8}/[a-zA-Z0-9_]{24}",
		"Stripe API Key": "sk_live_[0-9a-zA-Z]{24}",
		"Stripe Restricted API Key": "rk_live_[0-9a-zA-Z]{24}",
		"Square Access Token": "sq0atp-[0-9A-Za-z\\-_]{22}",
		"Square OAuth Secret": "sq0csp-[0-9A-Za-z\\-_]{43}",
		"Twilio API Key": "SK[0-9a-fA-F]{32}",
		"Twitter Access Token": "[t|T][w|W][i|I][t|T][t|T][e|E][r|R].*[1-9][0-9]+-[0-9a-zA-Z]{40}",
		"Twitter OAuth": "[t|T][w|W][i|I][t|T][t|T][e|E][r|R].*['|\"][0-9a-zA-Z]{35,44}['|\"]",
		/*
			MY RULES ^_^
		*/
		// "Random Code": "[\"' ][0-9A-Za-z]{20,}[\"' ]",
		"Hardcoded Password": "(?i)(password).*[=:\\s][\"']\\S+\\s", //Slightly better regex for passwords
		"Fastly API Key": "\\W(Fastly-key)\\W*[A-Za-z0-9+=]{44,}\\W", //is it b64 tho?
		"Disqus API Key": "\\W(?i)(disqus).+\\w[K|k][E|e][Y|y]\\W+[A-Za-z0-9]{64}\\W",
		"Zoho Desk Token": "[0-9]{4}(.)[0-9a-f]{32}(.)[0-9a-f]{32}", //https://desk.zoho.com/DeskAPIDocument
		// "Auth": "\\W(Authorization:).+\\W",
		"Auth Bearer": "(Authorization: Bearer )[a-zA-z0-9]{20,}\\S",
		"Auth Basic": "(Basic )[a-zA-z0-9+=]{20,}\\S",
		"Hash Token": "\\w*[T|t][O|o][K|k][E|e][N|n]\\W*([0-9a-f]{32}|[0-9a-f]{40}|[0-9a-f]{64})\\W", //md5, sha1, sha256
		"Hash Key": "\\w*[K|k][E|e][Y|y]\\W*([0-9a-f]{32}|[0-9a-f]{40}|[0-9a-f]{64})\\W", //md5, sha1, sha256
		"Random API Key": "\\S*[K|k][E|e][Y|y]+\\W+[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}",
		"Meraki API Key": "[X|x]-[C|c][I|i][S|s][C|c][O|o]+-[M|m][E|e][R|r][A|a][K|k][I|i].+[:=]\\W+[0-9a-f]{40}\\W",
	}

	for rule, regex := range rules {
		// matched, err := regexp.Match(regex, []byte(line))
		re := regexp.MustCompile(regex)
		ms := re.FindAllString(line,-1) //https://golang.org/pkg/regexp/#Regexp.FindAllString
		if(len(ms)>0) {
			result:=line
			log.Debug("Found:",rule)
			result=HighlightWords(line,ms)
			matches=append(matches,Match{
				Rule: MatchRule{ "Key", rule },
				Line: result,
			})
		}

		// if(re.MatchString(line)) {
		// 	log.Warning("Found:",rule)
		// 	matches=append(matches,line)
		// }

		// matches=append(matches, fmt.Sprintf("%v",re.FindAll([]byte(line), 0)))
		// log.Warn(matches)

		// matches = re.FindAll([]byte(line), -1)
		// if err != nil {
		// 	log.Fatal("Error: %v\n", err)
		// 	return false
		// }
		// if matched {
		// 	return true
		// }
	}
	return matches
}
