# RepoCrawl
Crawl GitHub/Bitbucket/Gitlab/Git repositories in search for unsafely stored secrets. Completely written in Go

## Overview
Inspired by other tools like [GitGot](https://github.com/BishopFox/GitGot/), [Trufflehog](https://github.com/dxa4481/truffleHog/), [GitRob](https://github.com/michenriksen/gitrob/) and the like, I decided to develop a tool that takes the best of these tools and impletements their methods on most `Git` compatible services.

This tool crawls repositories on various Git services using a variety of methods that ultimately search for secrets/api tokens/sessions tokens/passwords/private keys that should otherwise be private.

### Available Crawling methods
- DeepSearch (Inspired by [Trufflehog](https://github.com/dxa4481/truffleHog/))
  - Given a Git reporitory:
    + Enumerate all commits and analyse `diff` contents for submitted secrets
    + Enumerate all users that participated in the repo (commits) -> Look for public repositories belonging to these users -> Enumerate all commits and analyse `diff` contents for submitted secrets #**TODO** asciiflow chart
- Github
  - Given a `string` search paramter (Inspired by [GitGot](https://github.com/BishopFox/GitGot/) functionality):
    + Search the entirety of Github in search of the submitted string `/search/code`
    + Loop through results and search for secrets in matched files 
    + Give user an interactive view of each result until trying another one · giving the possibility to ignore `repo` or `user` or do a `deepsearch` on the current repository
  - Given a Github repository / `User` / `Organization` do a `deepsearch`
- Bitbucket
  - Given a Bitbucket repository / `User` / `Group` do a `deepsearch`
- Gitlab
  - Given a Gitlab repository / `User` / `Organization` do a `deepsearch`

## Instalation

```go get github.com/mariolima/RepoCrawl``` TODO

## Credits
- mariolima
- [GitGot](https://github.com/BishopFox/GitGot/)
- [Trufflehog](https://github.com/dxa4481/truffleHog/)
- [GitRob](https://github.com/michenriksen/gitrob/)
