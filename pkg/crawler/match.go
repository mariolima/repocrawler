package crawler

import (
	"github.com/mariolima/repocrawl/internal/entities"
)

// Match Secret found while crawling a repository
type Match struct { //Has to be generic - TODO move to other pkg
	Rule   MatchRule
	Line   string
	LineNr int
	Values []string
	//Repository struct // User struct and other generic stuff
	URL          string
	Entropy      float64
	SearchResult entities.SearchResult
}

// type MatchValue struct {
// 	Value string
// 	Entropy float64
// }

// MatchRule Used to identify secrets through Regexes
type MatchRule struct {
	Type  string
	Regex string
}
