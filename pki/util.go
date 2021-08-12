package pki

import (
	"crypto/rand"
	"encoding/pem"
	"errors"
	"math/big"
)

func encode(t PEMType, b []byte) string {
	block := &pem.Block{
		Type:  string(t),
		Bytes: b,
	}
	return string(pem.EncodeToMemory(block))
}

func decode(raw string) (*pem.Block, error) {
	block, _ := pem.Decode([]byte(raw))
	if block == nil {
		return nil, errors.New("invalid PEM")
	}
	return block, nil
}

func random() *big.Int {
	serial, err := rand.Int(rand.Reader, new(big.Int).Lsh(big.NewInt(1), 128))
	if err != nil {
		panic(err)
	}
	return serial
}
