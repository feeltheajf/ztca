package pki

import (
	"crypto/x509"

	"github.com/feeltheajf/ztca/fs"
)

// ReadRevocationList loads certificate revocation list from file
func ReadRevocationList(filename string) (*x509.RevocationList, error) {
	d, err := fs.Read(filename)
	if err != nil {
		return nil, err
	}
	return UnmarshalRevocationList(d)
}

// UnmarshalRevocationList parses certificate revocation list from PEM-encoded string
func UnmarshalRevocationList(raw string) (*x509.RevocationList, error) {
	crl, err := x509.ParseCRL([]byte(raw))
	if err != nil {
		return nil, err
	}
	template := &x509.RevocationList{
		RevokedCertificates: crl.TBSCertList.RevokedCertificates,
	}
	return template, nil
}
