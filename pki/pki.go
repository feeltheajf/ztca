package pki

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/asn1"
	"math/big"
	"os"
	"sync"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/feeltheajf/ztca/errdefs"
	"github.com/feeltheajf/ztca/fs"
)

var (
	ctx    zerolog.Logger
	lock   sync.RWMutex
	config *Config

	caCrt *x509.Certificate
	caKey *ecdsa.PrivateKey

	defaultKeyUsageCA        = x509.KeyUsageCertSign | x509.KeyUsageCRLSign
	defaultKeyUsageClient    = x509.KeyUsageDigitalSignature
	defaultExtKeyUsageClient = []x509.ExtKeyUsage{
		x509.ExtKeyUsageClientAuth,
	}

	defaultUpdateCRL = 4 * time.Hour

	oidCRLReason = asn1.ObjectIdentifier{2, 5, 29, 21}
)

// Config holds CA configuration
type Config struct {
	Certificate       string `yaml:"certificate"`
	PrivateKey        string `yaml:"privateKey"`
	ExpirationDays    int    `yaml:"expirationDays"`
	CertificateURL    string `yaml:"certificateUrl" bind:"required"`
	CRL               string `yaml:"crl"`
	CRLExpirationDays int    `yaml:"crlExpirationDays"`
	CRLURL            string `yaml:"crlUrl" bind:"required"`
	CRLNumber         string `yaml:"crlNumber" bind:"required"`
}

// Setup initializes CA
func Setup(cfg *Config) (err error) {
	ctx = log.With().Str("module", "ca").Logger()
	config = cfg

	caCrt, err = ReadCertificate(cfg.Certificate)
	if err != nil {
		return err
	}

	caKey, err = ReadPrivateKey(cfg.PrivateKey)
	if err != nil {
		return err
	}

	if err := NewRevocationList(); err != nil {
		return err
	}

	go updateCRL()
	return nil
}

// NewCertificate issues a new certificate using the given template
func NewCertificate(template *x509.Certificate, pub crypto.PublicKey) (*x509.Certificate, error) {
	template.BasicConstraintsValid = true

	if template.NotBefore.IsZero() {
		template.NotBefore = time.Now()
	}

	if template.NotAfter.IsZero() {
		template.NotAfter = template.NotBefore.AddDate(0, 0, config.ExpirationDays)
	}

	template.KeyUsage = defaultKeyUsageClient
	template.ExtKeyUsage = defaultExtKeyUsageClient
	template.SerialNumber = random()
	template.IssuingCertificateURL = []string{config.CertificateURL}
	template.CRLDistributionPoints = []string{config.CRLURL}

	if template.Equal(caCrt) {
		template.IsCA = true
		template.KeyUsage = defaultKeyUsageCA
		template.ExtKeyUsage = nil
	}

	b, err := x509.CreateCertificate(rand.Reader, template, caCrt, pub, caKey)
	if err != nil {
		return nil, errdefs.Unknown("failed to issue certificate").CausedBy(err)
	}

	crt, err := x509.ParseCertificate(b)
	if err != nil {
		return nil, errdefs.Unknown("failed to parse certificate").CausedBy(err)
	}

	ctx.Info().
		Str("serial", MarshalSerial(crt.SerialNumber)).
		Str("username", crt.Subject.CommonName).
		Str("device_serial", crt.Subject.SerialNumber).
		Msg("certificate issued")
	return crt, nil
}

// CertificateOption are used for easier template generation
type CertificateOption func(*x509.Certificate) error

// NewTemplate generates new x509 certificate with the given options
func NewTemplate(opts ...CertificateOption) (*x509.Certificate, error) {
	template := &x509.Certificate{}
	for _, opt := range opts {
		if err := opt(template); err != nil {
			return nil, err
		}
	}

	return template, nil
}

// WithName sets certificate subject to the given PKIX name
func WithName(name pkix.Name) CertificateOption {
	return func(template *x509.Certificate) error {
		template.Subject = name
		return nil
	}
}

// Revoke given certificate
func Revoke(crt *x509.Certificate, reason CRLReason, when time.Time) error {
	if reason == "" {
		reason = CRLReasonUnspecified
	}

	if when.IsZero() {
		when = time.Now()
	}

	revoke := pkix.RevokedCertificate{
		SerialNumber:   crt.SerialNumber,
		RevocationTime: when,
		Extensions: []pkix.Extension{
			{
				Id:    oidCRLReason,
				Value: []byte(reason),
			},
		},
	}

	if err := NewRevocationList(revoke); err != nil {
		return err
	}

	ctx.Info().
		Str("serial", MarshalSerial(crt.SerialNumber)).
		Str("username", crt.Subject.CommonName).
		Str("device_serial", crt.Subject.SerialNumber).
		Msg("certificate issued")
	return nil
}

// NewRevocationList issues a new certificate revocation list
func NewRevocationList(revoke ...pkix.RevokedCertificate) error {
	lock.Lock()
	defer lock.Unlock()

	crl := new(x509.RevocationList)
	if _, err := os.Stat(config.CRL); os.IsNotExist(err) {
		log.Warn().Str("crl", config.CRL).Msg("crl not found")
	} else {
		crl, err = ReadRevocationList(config.CRL)
		if err != nil {
			return errdefs.Unknown("failed to load crl").CausedBy(err)
		}
	}

	srl := new(big.Int)
	if _, err := os.Stat(config.CRLNumber); os.IsNotExist(err) {
		log.Warn().Str("crl", config.CRLNumber).Msg("crl number not found")
	} else {
		srl, err = ReadSerial(config.CRLNumber)
		if err != nil {
			return errdefs.Unknown("failed to load crl number").CausedBy(err)
		}
	}
	srl = srl.Add(srl, big.NewInt(1))

	now := time.Now()
	crl.Number = srl
	crl.ThisUpdate = now
	crl.NextUpdate = now.AddDate(0, 0, config.CRLExpirationDays)
	crl.RevokedCertificates = append(crl.RevokedCertificates, revoke...)

	b, err := x509.CreateRevocationList(rand.Reader, crl, caCrt, caKey)
	if err != nil {
		return errdefs.Unknown("failed to create crl").CausedBy(err)
	}

	if err := fs.Write(config.CRL, encode(PEMTypeRevocationList, b)); err != nil {
		return errdefs.Unknown("failed to save crl").CausedBy(err)
	}

	if err := WriteSerial(config.CRLNumber, srl); err != nil {
		return errdefs.Unknown("failed to save crl number").CausedBy(err)
	}

	ctx.Info().
		Str("number", MarshalSerial(srl)).
		Int("revoked", len(crl.RevokedCertificates)).
		Msg("crl issued")
	return nil
}

func updateCRL() {
	for range time.Tick(defaultUpdateCRL) {
		if err := NewRevocationList(); err != nil {
			log.Error().Err(err).Msg("failed to update crl")
		}
	}
}
