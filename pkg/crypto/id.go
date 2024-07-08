package crypto

import (
	"crypto/rand"
	"encoding/hex"
	"io"
	mathrand "math/rand"
	"time"
)

func GenerateId(prefix string, size int) string {
	if size <= 0 || size >= 100 {
		size = 32
	}
	b := make([]byte, size)
	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		return ""
	}
	return prefix + "_" + hex.EncodeToString(b)
}

func GenerateRandomDuration(duration time.Duration) time.Duration {
	n := mathrand.Int63n(duration.Milliseconds())
	return time.Duration(n) * time.Millisecond
}
