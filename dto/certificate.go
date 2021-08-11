package dto

import (
	"fmt"
	"math/big"
	"time"

	"github.com/feeltheajf/ztca/pki"
)

type Certificate struct {
	Model
	Raw          string    `json:"raw"`
	SerialNumber string    `json:"serialNumber"`
	ExpiresAt    time.Time `json:"expiresAt"`

	Username     string `json:"username"`
	DeviceSerial string `json:"deviceSerial"`
}

func UnmarshalCertificate(raw string) (*Certificate, error) {
	x509, err := pki.UnmarshalCertificate(raw)
	if err != nil {
		return nil, err
	}

	crt := &Certificate{
		Raw:          raw,
		SerialNumber: FormatCertificateSerial(x509.SerialNumber),
		ExpiresAt:    x509.NotAfter,
		Username:     x509.Subject.CommonName,
		DeviceSerial: x509.Subject.SerialNumber,
	}

	return crt, nil
}

func FormatCertificateSerial(serial *big.Int) string {
	return fmt.Sprintf("%X", serial)
}

func FormatYubiKeySerial(serial uint32) string {
	return fmt.Sprintf("%d", serial)
}
