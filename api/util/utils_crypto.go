package util

import (
	"crypto/rand"
	"encoding/hex"
)

func RandomHex(n int) (string, error) {
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		GetLogger().Errorf("cannot create %d-byte long random hex: %v\n", n, err)
		return "", err
	}

	return hex.EncodeToString(b), nil
}
