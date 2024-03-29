# RepoCrawler
Crawl GitHub/Bitbucket/Gitlab/Git repositories in search for unsafely stored secrets. Completely written in Go

## Overview
Inspired by other tools like [GitGot](https://github.com/BishopFox/GitGot/), [Trufflehog](https://github.com/dxa4481/truffleHog/), [GitRob](https://github.com/michenriksen/gitrob/), [Git-All-Secrets](https://github.com/anshumanbh/git-all-secrets) and many others, I decided to develop a tool that takes the best of these tools and impletements optimized crawling to most `Git` compatible services.

This tool crawls repositories on various Git services using a variety of methods that ultimately search for secrets/api tokens/sessions tokens/passwords/private keys that should otherwise be private.

### Available Crawling methods
- **DeepCrawl** (Inspired by [Trufflehog](https://github.com/dxa4481/truffleHog/))
  - Given a Git reporitory:
    + Enumerate all commits and analyse `diff` contents for submitted secrets
    + Enumerate all users that participated in the repo (commits) -> Look for public repositories belonging to these users -> Enumerate all commits and analyse `diff` contents for submitted secrets
    ```
                    +-------------------------------------+
                    | Users that participated in the repo |
                    +------------------+------------------+
                                       |
                     +-----------------+
                  +--v---+         +---v--+
                  |User 1|         |User 2|
          +-------------------+    +------+
          |                   |
    +-----v------+     +------v-----+
    |Repository 1|     |Repository 2|
    +-----+------+     +------+-----+
          |                   |
     +----v----+         +----v----+
     |DeepCrawl|         |DeepCrawl|
     +---------+         +---------+
    ```
- Github
  - Given a `string` search parameter (Inspired by [GitGot](https://github.com/BishopFox/GitGot/) functionality):
    + Search the entirety of Github in search of code containing the submitted string `/search/code API call`
    + Loop through results and search for secrets in matched files
  - Given a Github repository / `User` / `Organization` do a `deepcrawl`
- Bitbucket
  - Given a Bitbucket repository / `User` / `Group` do a `deepcrawl`
- Gitlab
  - Given a Gitlab repository / `User` / `Organization` do a `deepcrawl`

## Instalation
### From source
```sht
go get github.com/mariolima/repocrawler
cd ~/go/src/github.com/mariolima/repocrawler/cmd/crawler-cli
go build .
export LOG_LEVEL=info
export GITHUB_ACCESS_TOKEN=TOKEN
./crawler-cli -h
```
### Using Docker
```sh
git clone github.com/mariolima/repocrawler
cd repocrawler
docker build . -t repocrawler
docker run -it -e 'GITHUB_ACCESS_TOKEN=TOKEN' -e 'SLACK_WEBHOOK=YOURWEBHOOK' repocrawler -h
```


## Structure
![](docs/diag.png)

## Packages used and why

| Package                   |     |
|---------------------------|-----------|
| [logrusorgru/aurora](github.com/logrusorgru/aurora) | Color highlights in CLI matches - theres others but this one works the best |
| [sirupsen/logrus](github.com/sirupsen/logrus) | Logs used for debugging - easy to setup and allows multiple verbose levels |
| [src-d/go-git.v4](gopkg.in/src-d/go-git.v4) | Clone and get data for Git servives - Heavy since most functions aren't used but no viable alternatives |
| [bndr/gopencils](github.com/bndr/gopencils) | Bitbucket/Gitlab API requests - Why not just use `net` ? Easier to code requests with multiple parameters/queries - No API bindings of Bitbucket/Gitlab for Go that support the calls needed |
| [google/go-github](github.com/google/go-github) | API calls to Github (Get Repositories/Users/...)
| [pkg/profile](github.com/pkg/profile) | Debuging and benchmarking with `pprof` |
| [x/oauth2](golang.org/x/oauth2) | Used for the `go-github` Github Token authentication request - def. overkill  |

## FAQ
-

## Credits
- mariolima
- [GitGot](https://github.com/BishopFox/GitGot/)
- [Trufflehog](https://github.com/dxa4481/truffleHog/)
- [GitRob](https://github.com/michenriksen/gitrob/)
