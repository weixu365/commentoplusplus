package util

import (
	"crypto/rand"
	"encoding/hex"
)

func randomHex(n int) (string, error) {
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		util.GetLogger().Errorf("cannot create %d-byte long random hex: %v\n", n, err)
		return "", errorInternal
	}

	return hex.EncodeToString(b), nil
}
