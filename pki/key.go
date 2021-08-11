package pki

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"io/ioutil"
)

var (
	// EllipticCurve is the default curve used for key generation
	EllipticCurve = elliptic.P256()
)

// NewPrivateKey generates new private key using `EllipticCurve`
func NewPrivateKey() (crypto.PrivateKey, error) {
	return ecdsa.GenerateKey(EllipticCurve, rand.Reader)
}

// ReadPrivateKey loads private key from file
func ReadPrivateKey(filename string) (crypto.PrivateKey, error) {
	b, err := ioutil.ReadFile(filename) // #nosec: G304
	if err != nil {
		return nil, err
	}
	return UnmarshalPrivateKey(string(b))
}

// UnmarshalPrivateKey parses private key from PEM-encoded bytes
func UnmarshalPrivateKey(raw string) (crypto.PrivateKey, error) {
	block, _ := pem.Decode([]byte(raw))
	if block == nil {
		return nil, errors.New("failed to parse private key: invalid PEM")
	}
	return x509.ParseECPrivateKey(block.Bytes)
}

// MarshalPublicKey returns PEM encoding of key
func MarshalPublicKey(key crypto.PublicKey) (string, error) {
	b, err := x509.MarshalPKIXPublicKey(key)
	if err != nil {
		return "", err
	}

	block := &pem.Block{
		Type:  pemTypePublicKey,
		Bytes: b,
	}
	return string(pem.EncodeToMemory(block)), nil
}
