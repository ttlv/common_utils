package preauth

import (
	"math/rand"
	"time"
)

type RandMaker interface {
	Gen(int) string
}

type DefaultRandMaker struct {
}

func (m DefaultRandMaker) Gen(n int) string {
	var letter = []rune("0123456789")
	rand.Seed(time.Now().UnixNano())
	b := make([]rune, n)
	for i := range b {
		b[i] = letter[rand.Intn(len(letter))]
	}
	return string(b)
}

type MockRankMaker struct {
}

func (m MockRankMaker) Gen(n int) string {
	return "012345"
}
