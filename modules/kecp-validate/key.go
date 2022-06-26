package kecpvalidate

import (
	"math"
	"strings"
)

func IsAValidCryptoKey(s string) bool {
	return len(s) == 64 && entropy(s) > 2
}

func entropy(s string) (e float64) {
	m := make(map[rune]bool)
	for _, r := range s {
		if m[r] {
			continue
		}
		m[r] = true
		n := strings.Count(s, string(r))
		p := float64(n) / float64(len(s))
		e += p * math.Log(p) / math.Log(2)
	}
	return math.Abs(e)
}
