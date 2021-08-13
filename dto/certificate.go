package dto

import (
	"crypto/x509"
	"time"

	"github.com/feeltheajf/ztca/pki"
)

var (
	Certificates = &certificateService{}
)

type Certificate struct {
	Model
	Raw          string    `json:"raw"`
	SerialNumber string    `json:"serialNumber"`
	ExpiresAt    time.Time `json:"expiresAt"`
	// user metadata
	Username     string `json:"username"`
	DeviceSerial string `json:"deviceSerial"`
}

func (crt *Certificate) X509() *x509.Certificate {
	x509, err := pki.UnmarshalCertificate(crt.Raw)
	if err != nil {
		panic(err)
	}
	return x509
}

type certificateService service

func (cs *certificateService) Save(crt *Certificate) error {
	return db.Create(crt).Error
}

func (cs *certificateService) Load(username string) (*Certificate, error) {
	crt := new(Certificate)
	return crt, db.Where("username = ?", username).First(crt).Error
}

func (cs *certificateService) Delete(crt *Certificate) error {
	return db.Delete(crt).Error
}
