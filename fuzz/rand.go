package fuzz

import (
	"math/rand"
	"time"
)

var asciiRanges = []characterRange{
	{'A', 'Z'},
	{'a', 'z'},
	{'0', '9'},
}

type characterRange struct {
	first, last rune
}

func (r *characterRange) choose(rand *rand.Rand) rune {
	count := int64(r.last - r.first)
	return r.first + rune(rand.Int63n(count))
}

func randAlphaString(r *rand.Rand, length int) string {
	runes := make([]rune, length)
	for i := range runes {
		runes[i] = asciiRanges[r.Intn(2)].choose(r)
	}
	return string(runes)
}

func randDateTime() string {
	min := time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Unix()
	sec := rand.Int63n(time.Now().Unix()-min) + min
	return time.Unix(sec, 0).Format(time.RFC3339)
}
