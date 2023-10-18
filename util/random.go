package util

import (
	"math/rand"
	"strings"
)

// write to generate random int and string value

const (
	alphabet = "asdfghjklzxcvbnmqwertyuiop"
)

// RandomInt generate a random integer within [min,max]
func RandomInt(min, max int64) int64 {
	return min + rand.Int63n(max-min+1)
}

// RandomString generate a random string with length n
func RandomString(n int) string {
	k := len(alphabet)
	var str strings.Builder

	for i := 0; i < n; i++ {
		c := alphabet[rand.Intn(k)]
		str.WriteByte(c)
	}

	return str.String()
}

func RandomOwner() string {
	return RandomString(5)
}

func RandomMoney() int64 {
	return RandomInt(10, 1000)
}

func RandomCurrency() string {
	// "EUR", "USD", "CAD"
	concur := []string{EUR, USD, CAD}

	k := len(concur)

	return concur[rand.Intn(k)]
}
