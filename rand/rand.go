package rand

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
)

func Bytes(n int) ([]byte, error) {
	b := make([]byte, n)
	nRead, err := rand.Read(b)
	if err != nil {
		return nil, fmt.Errorf("bytes: %w", err)
	}
	if nRead != n {
		return nil, fmt.Errorf("bytes: read %d bytes, expected %d bytes", nRead, n)
	}

	return b, nil
}

// String generates a random string using crypto/rand
// n is the length of the string
func String(n int) (string, error) {
	b, err := Bytes(n)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}
