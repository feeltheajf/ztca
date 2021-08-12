package dto

import (
	"time"
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
