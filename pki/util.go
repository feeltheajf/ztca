package pki

import (
	"crypto/rand"
	"encoding/pem"
	"errors"
	"fmt"
	"math/big"
)

// MarshalCertificateSerial converts certificate serial to string
func MarshalCertificateSerial(serial *big.Int) string {
	return fmt.Sprintf("%X", serial)
}

// MarshalCertificateSerial converts string to certificate serial
func UnmarshalCertificateSerial(serial string) *big.Int {
	srl := new(big.Int)
	srl.SetString(serial, 16)
	return srl
}

// MarshalYubiKeySerial converts YubiKey serial to string
func MarshalYubiKeySerial(serial uint32) string {
	return fmt.Sprintf("%d", serial)
}

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
