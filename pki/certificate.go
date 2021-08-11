package pki

import (
	"crypto/x509"
	"encoding/pem"
	"errors"
	"io/ioutil"
)

// ReadCertificate loads certificate from file
func ReadCertificate(filename string) (*x509.Certificate, error) {
	b, err := ioutil.ReadFile(filename) // #nosec: G304
	if err != nil {
		return nil, err
	}
	return UnmarshalCertificate(string(b))
}

// UnmarshalCertificate parses certificate from PEM-encoded bytes
func UnmarshalCertificate(raw string) (*x509.Certificate, error) {
	block, _ := pem.Decode([]byte(raw))
	if block == nil {
		return nil, errors.New("failed to parse certificate: invalid PEM")
	}
	return x509.ParseCertificate(block.Bytes)
}

// MarshalCertificate returns PEM encoding of crt
func MarshalCertificate(crt *x509.Certificate) (string, error) {
	block := &pem.Block{
		Type:  pemTypeCertificate,
		Bytes: crt.Raw,
	}
	return string(pem.EncodeToMemory(block)), nil
}
