package pki

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"io/ioutil"

	"github.com/feeltheajf/ztca/x/fs"
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
	return UnmarshalPrivateKey(b)
}

// UnmarshalPrivateKey parses private key from PEM-encoded bytes
func UnmarshalPrivateKey(raw []byte) (crypto.PrivateKey, error) {
	block, _ := pem.Decode(raw)
	if block == nil {
		return nil, errors.New("failed to parse private key: invalid PEM")
	}
	return x509.ParseECPrivateKey(block.Bytes)
}

// WritePrivateKey saves private key to file
func WritePrivateKey(filename string, key crypto.PrivateKey) error {
	raw, err := MarshalPrivateKey(key)
	if err != nil {
		return err
	}
	return fs.WriteFile(filename, raw)
}

// MarshalPrivateKey returns PEM encoding of key
func MarshalPrivateKey(key crypto.PrivateKey) ([]byte, error) {
	var b []byte
	var err error
	var pemType string

	switch k := key.(type) {
	case *ecdsa.PrivateKey:
		pemType = pemTypeECPrivateKey
		b, err = x509.MarshalECPrivateKey(k)
	default:
		return nil, fmt.Errorf("unsupported private key type: %T", key)
	}

	if err != nil {
		return nil, err
	}

	block := &pem.Block{
		Type:  pemType,
		Bytes: b,
	}
	return pem.EncodeToMemory(block), nil
}

// WritePublicKey saves public key to file
func WritePublicKey(filename string, key crypto.PublicKey) error {
	b, err := MarshalPublicKey(key)
	if err != nil {
		return err
	}
	return fs.WriteFile(filename, b)
}

// MarshalPublicKey returns PEM encoding of key
func MarshalPublicKey(key crypto.PublicKey) ([]byte, error) {
	b, err := x509.MarshalPKIXPublicKey(key)
	if err != nil {
		return nil, err
	}

	block := &pem.Block{
		Type:  pemTypePublicKey,
		Bytes: b,
	}
	return pem.EncodeToMemory(block), nil
}
