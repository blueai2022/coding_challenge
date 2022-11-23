package domain

import (
	"math/rand"
	"sort"
	"strings"
	"time"
)

// for generate random string
const alphabet = "abcdefghijklmnopqrstuvwxyz"

func init() {
	rand.Seed(time.Now().UnixNano())
}

// randomString generates a random string of length n
func randomString(n int) string {
	var sb strings.Builder
	k := len(alphabet)

	for i := 0; i < n; i++ {
		c := alphabet[rand.Intn(k)]
		sb.WriteByte(c)
	}

	return sb.String()
}

func contains(strs []string, match string) bool {
	for _, s := range strs {
		if s == match {
			return true
		}
	}
	return false
}

func equals(a []string, b []string, ignoreOrder bool) bool {
	if len(a) != len(b) {
		return false
	}

	if ignoreOrder {
		sort.Strings(a)
		sort.Strings(b)
	}

	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
