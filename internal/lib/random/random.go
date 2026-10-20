package random

import "math/rand/v2"

var chars = []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ" +
	"abcdefghijklmnopqrstuvwxyz" +
	"0123456789")

// NewRandomString generates random string with given size.
func NewRandomString(size int) string {
	b := make([]rune, size)
	for i := range b {
		b[i] = chars[rand.IntN(len(chars))]
	}

	return string(b)
}
