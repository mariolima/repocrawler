package crawler

import (
	"math"
)

// FindEntropy shannonEntropy - https://www.reddit.com/r/dailyprogrammer/comments/4fc896/20160418_challenge_263_easy_calculating_shannon/d2e1wr1/
func FindEntropy(input string) float64 {
	charMap := make(map[rune]int)
	for _, c := range input {
		if _, ok := charMap[c]; !ok {
			charMap[c] = 1
		} else {
			charMap[c]++
		}
	}

	var sum float64
	length := float64(len(input))
	for _, cnt := range charMap {
		tmp := float64(cnt) / length
		sum += tmp * math.Log2(tmp)
	}
	return -1 * sum
}
