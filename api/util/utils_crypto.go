package util

import (
	"crypto/rand"
	"encoding/hex"
	"simple-commenting/app"
)

func randomHex(n int) (string, error) {
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		util.GetLogger().Errorf("cannot create %d-byte long random hex: %v\n", n, err)
		return "", app.ErrorInternal
	}

	return hex.EncodeToString(b), nil
}
