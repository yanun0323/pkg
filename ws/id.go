package ws

import (
	"crypto/rand"
	"encoding/hex"
)

func newLogID() string {
	buf := make([]byte, 4)
	if _, err := rand.Read(buf); err != nil {
		return newLogID()
	}
	return hex.EncodeToString(buf)
}
