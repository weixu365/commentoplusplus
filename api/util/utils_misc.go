package util

import (
	"bytes"
)

func concat(a bytes.Buffer, b bytes.Buffer) []byte {
	return append(a.Bytes(), b.Bytes()...)
}
