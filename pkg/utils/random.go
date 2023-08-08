package utils

import (
	crypt_rand "crypto/rand"
	"io"
	"math/rand"
	"time"
)

var randSrc = rand.NewSource(time.Now().UnixNano())

const alpha = `abcdefghijklmnopqrstuvwxyz` +
	`ABCDEFGHIJKLMNOPQRSTUVWXYZ`
const num = `0123456789`
const charset = alpha + num

// create a random string of specified character length
func GenerateRandStr(l int) string {
	r := rand.New(randSrc)
	b := make([]byte, l)
	for i := range b {
		if i == 0 {
			b[i] = alpha[r.Intn(len(alpha))]
			continue
		}
		b[i] = charset[r.Intn(len(charset))]
	}
	return string(b)
}

// create random secret string made up of all digits
func GenerateRandDigitStr(length int) string {
	b := make([]byte, length)
	n, err := io.ReadAtLeast(crypt_rand.Reader, b, length)
	if n != length {
		panic(err)
	}
	for i := 0; i < len(b); i++ {
		b[i] = num[int(b[i])%len(num)]
	}
	return string(b)

}
