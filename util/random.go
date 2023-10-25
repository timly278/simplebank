package util

import (
	"fmt"
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

// RandomOwner generates a random owner name
func RandomOwner() string {
	return RandomString(5)
}

// RandomMoney generates a random amount of money
func RandomMoney() int64 {
	return RandomInt(10, 1000)
}

// RandomCurrency generates a random currency code
func RandomCurrency() string {
	// "EUR", "USD", "CAD"
	concur := []string{EUR, USD, CAD}

	k := len(concur)

	return concur[rand.Intn(k)]
}

// RandomEmail generates a random email
func RandomEmail() string {
	return fmt.Sprintf("%s@gmail.com", RandomString(6))
}
