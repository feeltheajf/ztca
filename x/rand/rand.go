package rand

import (
	"crypto/rand"
	"encoding/hex"
)

const (
	defaultRandomLength = 20
)

func Bytes() []byte {
	b := make([]byte, defaultRandomLength)
	_, err := rand.Read(b)
	if err != nil {
		panic(err)
	}
	return b
}

func String() string {
	return hex.EncodeToString(Bytes())
}
