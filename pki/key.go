package pki

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"encoding/base64"
	"fmt"

	"golang.org/x/crypto/ssh"

	"github.com/feeltheajf/ztca/fs"
)

var (
	// EllipticCurve is the default curve used for key generation
	EllipticCurve = elliptic.P256()
)

// NewPrivateKey generates new private key using `EllipticCurve`
func NewPrivateKey() (*ecdsa.PrivateKey, error) {
	return ecdsa.GenerateKey(EllipticCurve, rand.Reader)
}

// ReadPrivateKey loads private key from file
func ReadPrivateKey(filename string) (crypto.PrivateKey, error) {
	raw, err := fs.Read(filename)
	if err != nil {
		return nil, err
	}
	return UnmarshalPrivateKey(raw)
}

// UnmarshalPrivateKey parses private key from PEM-encoded string
func UnmarshalPrivateKey(raw string) (crypto.PrivateKey, error) {
	block, err := decode(raw)
	if err != nil {
		return nil, err
	}
	return x509.ParseECPrivateKey(block.Bytes)
}

// WritePrivateKey saves private key to file
func WritePrivateKey(filename string, key crypto.PrivateKey) error {
	raw, err := MarshalPrivateKey(key)
	if err != nil {
		return err
	}
	return fs.Write(filename, raw)
}

// MarshalPrivateKey returns PEM encoding of key
func MarshalPrivateKey(key crypto.PrivateKey) (string, error) {
	var b []byte
	var err error
	var t PEMType

	switch k := key.(type) {
	case *ecdsa.PrivateKey:
		t = PEMTypeECPrivateKey
		b, err = x509.MarshalECPrivateKey(k)
	default:
		return "", fmt.Errorf("unsupported private key type: %T", key)
	}
	if err != nil {
		return "", err
	}

	return encode(t, b), nil
}

// WritePublicKey saves public key to file
func WritePublicKey(filename string, key crypto.PublicKey) error {
	raw, err := MarshalPublicKey(key)
	if err != nil {
		return err
	}
	return fs.Write(filename, raw)
}

// MarshalPublicKey returns PEM encoding of key
func MarshalPublicKey(key crypto.PublicKey) (string, error) {
	b, err := x509.MarshalPKIXPublicKey(key)
	if err != nil {
		return "", err
	}
	return encode(PEMTypePublicKey, b), nil
}

// WritePublicKeySSH saves public key to file in OpenSSH format
func WritePublicKeySSH(filename string, key crypto.PublicKey) error {
	raw, err := MarshalPublicKeySSH(key)
	if err != nil {
		return err
	}
	return fs.Write(filename, raw)
}

// MarshalPublicKeySSH returns OpenSSH encoding of key
func MarshalPublicKeySSH(key crypto.PublicKey) (string, error) {
	pub, err := ssh.NewPublicKey(key)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(pub.Marshal()), nil
}
