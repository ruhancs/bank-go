package util

import (
	"math/rand"
	"strings"
	"time"
)

const alphabet = "abcdefghijklmnopqrstuwvxyz"

func init() {
	rand.New(rand.NewSource(time.Now().UnixNano()))
}

func RandomInt(min,max int64) int64 {
	return min + rand.Int63n(max-min + 1)//min -> max
}

func RandomString(n int) string {
	var sb strings.Builder
	k := len(alphabet)

	for i := 0; i < n; i++ {
		c := alphabet[rand.Intn(k)]
		sb.WriteByte(c)
	}

	return sb.String()
}

func RandomOwner() string {
	return RandomString(6)
}

func RandomMoney() int64 {
	return RandomInt(10,1000)
}

func RandomCurrency() string {
	currencies := []string{"EUR","USD", "BRL"}
	n := len(currencies)
	return currencies[rand.Intn(n)]
}