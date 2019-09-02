package utils

import(
	"strings"
	"fmt"
	. "github.com/logrusorgru/aurora" //colors - why this? because it is simple to replace text with colored text (others are not)
)

const MAX_LINE_SIZE = 400

func HighlightWords(line string, words []string) (res string) {
	if words == nil {
		return line
	}
	for _, word := range words {
		res=strings.ReplaceAll(line,word,fmt.Sprintf("%s",Red(word)))
		// res=strings.ReplaceAll(line,word,fmt.Sprintf("%s",Black(word).BgGreen()))
	}
	return res
}

func HighlightWord(line string, word string) string {
	// return strings.ReplaceAll(line,word,fmt.Sprintf("%s",Black(word).BgGreen()))
	return strings.ReplaceAll(line,word,fmt.Sprintf("%s",Green(word)))
}


func TruncateString(line string, matches []string, buffer, max int) (truncated string) {
	if len(line) < max {
		return line
	}
	padding := buffer						  // Data to be shown arround the matches
	loc := strings.Index(line,matches[0]) // Index location of 1st match
	if (loc<padding) {
		padding = loc
	}
	if (loc+padding+len(matches[0]) > len(line)-1) {
		padding = 0
	}
	truncated = fmt.Sprintf("...%s...",line[loc-padding:loc+padding+len(matches[0])])
	return truncated
}
