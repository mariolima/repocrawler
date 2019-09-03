package crawler

import (
	"github.com/mariolima/repocrawl/internal/entities"
)

type Match struct { //Has to be generic - TODO move to other pkg
	Rule   MatchRule
	Line   string
	LineNr int
	Values []string
	//Repository struct // User struct and other generic stuff
	URL          string
	SearchResult entities.SearchResult
}

type MatchRule struct {
	Type  string
	Regex string
}
