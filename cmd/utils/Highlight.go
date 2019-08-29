package utils

import(
	"strings"
	"fmt"
	"github.com/logrusorgru/aurora" //colors - why this? because it is simple to replace text with colored text (others are not)
)


func HighlightWords(line string, words []string) (res string) {
	if words == nil {
		return line
	}
	for _, word := range words {
		res=strings.ReplaceAll(line,word,fmt.Sprintf("%s",aurora.Red(word)))
	}
	return res
}

func HighlightWord(line string, word string) string {
	return strings.ReplaceAll(line,word,fmt.Sprintf("%s",aurora.Red(word)))
}
