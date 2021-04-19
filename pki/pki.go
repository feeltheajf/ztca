package pki

import (
	"crypto"
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"fmt"
	"math/big"
	"time"
)

const (
	pemTypeCertificate  = "CERTIFICATE"
	pemTypeECPrivateKey = "EC PRIVATE KEY"
	pemTypePublicKey    = "PUBLIC KEY"
)

var (
	caCrt *x509.Certificate
	caKey crypto.PrivateKey

	defaultKeyUsageCA = x509.KeyUsageCertSign | x509.KeyUsageCRLSign

	defaultKeyUsageProxy    = x509.KeyUsageDigitalSignature
	defaultExtKeyUsageProxy = []x509.ExtKeyUsage{
		x509.ExtKeyUsageClientAuth,
		x509.ExtKeyUsageServerAuth,
	}

	defaultExpirationYears = 1
)

type Config struct {
	CA *CA `yaml:"ca"`
}

type CA struct {
	Certificate string `yaml:"certificate"`
	PrivateKey  string `yaml:"privateKey"`
}

func Setup(cfg *Config) (err error) {
	caCrt, err = ReadCertificate(cfg.CA.Certificate)
	if err != nil {
		return err
	}

	caKey, err = ReadPrivateKey(cfg.CA.PrivateKey)
	if err != nil {
		return err
	}

	return nil
}

func NewCertificate(template *x509.Certificate, pub crypto.PublicKey) (*x509.Certificate, error) {
	template.BasicConstraintsValid = true
	template.PermittedDNSDomainsCritical = true

	if template.NotBefore.IsZero() {
		template.NotBefore = time.Now()
	}

	if template.NotAfter.IsZero() {
		template.NotAfter = template.NotBefore.AddDate(defaultExpirationYears, 0, 0)
	}

	template.KeyUsage = defaultKeyUsageProxy
	template.ExtKeyUsage = defaultExtKeyUsageProxy

	if template.Equal(caCrt) {
		template.IsCA = true
		template.KeyUsage = defaultKeyUsageCA
		template.ExtKeyUsage = nil
	}

	serialNumber, err := rand.Int(rand.Reader, new(big.Int).Lsh(big.NewInt(1), 128))
	if err != nil {
		return nil, err
	}
	template.SerialNumber = serialNumber

	der, err := x509.CreateCertificate(rand.Reader, template, caCrt, pub, caKey)
	if err != nil {
		return nil, fmt.Errorf("failed to create certificate: %s", err)
	}

	return x509.ParseCertificate(der)
}

type TemplateOption func(*x509.Certificate) error

func NewTemplate(opts ...TemplateOption) (*x509.Certificate, error) {
	template := &x509.Certificate{}
	for _, opt := range opts {
		if err := opt(template); err != nil {
			return nil, err
		}
	}

	return template, nil
}

func WithName(name pkix.Name) TemplateOption {
	return func(template *x509.Certificate) error {
		template.Subject = name
		return nil
	}
}

func WithCommonName(commonName string) TemplateOption {
	return WithName(
		pkix.Name{
			CommonName: commonName,
		},
	)
}

func WithDNSDomains(dnsDomains ...string) TemplateOption {
	return func(template *x509.Certificate) error {
		template.PermittedDNSDomains = dnsDomains
		return WithCommonName(dnsDomains[0])(template)
	}
}
