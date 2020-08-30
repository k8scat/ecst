package utils

import "math/rand"

func RandomString(source string, n int) string {
	rs := []rune(source)
	b := make([]byte, n)
	for i := range b {
		b[i] = source[rand.Intn(len(rs))]
	}
	return string(b)
}
