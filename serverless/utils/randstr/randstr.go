package randstr

import (
	"crypto/rand"
	"math/big"
)

const (
	LetterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789_"
	DigitBytes  = "0123456789"
)

// randBytes generates a random string of a specified length
// using the provided character set.
func randBytes(count int, charSet string) (string, error) {
	maxVal := big.NewInt(int64(len(charSet)))
	b := make([]byte, count)

	for i := range b {
		val, err := rand.Int(rand.Reader, maxVal)

		if err != nil {
			return "", err
		}

		b[i] = charSet[val.Int64()]
	}

	return string(b), nil
}

// RandStringBytes generates a random string of a specified length composed
// of numerals, upper and lowercase English characters, and underscores.
func RandStringBytes(count int) (string, error) {
	return randBytes(count, LetterBytes)
}

// RandDigitBytes generates a random string of numerals of a specified length.
func RandDigitBytes(count int) (string, error) {
	return randBytes(count, DigitBytes)
}
