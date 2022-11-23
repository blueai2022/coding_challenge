package stats

import (
	"math/rand"
	"os"
	"testing"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func TestMain(m *testing.M) {
	os.Exit(m.Run())
}

// randomInt generates a random integer between min and max
func randomInt(min, max int64) int64 {
	return min + rand.Int63n(max-min+1)
}

// randomAge generates a random integer between 0 and 200
func randomAge() int64 {
	return randomInt(0, 200)
}

// shuffle the int slice in random order
func shuffle(s []int) {
	rand.Shuffle(len(s), func(i, j int) { s[i], s[j] = s[j], s[i] })
}
