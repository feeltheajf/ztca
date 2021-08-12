package pki

// PEMType is used for encoding objects
type PEMType string

// Known PEM types
const (
	PEMTypeCertificate        PEMType = "CERTIFICATE"
	PEMTypeECPrivateKey       PEMType = "EC PRIVATE KEY"
	PEMTypePublicKey          PEMType = "PUBLIC KEY"
	PEMTypeRevocationList     PEMType = "X509 CRL"
	PEMTypeCertificateRequest PEMType = "CERTIFICATE REQUEST"
)

// CRLReason is used for designating certificate revocation reason.
// See
type CRLReason string

// Known CRL reasons
const (
	CRLReasonKeyCompromise      CRLReason = "KeyCompromise"
	CRLReasonAffiliationChanged CRLReason = "AffiliationChanged"
	CRLReasonSuperseded         CRLReason = "Superseded"
	CRLReasonUnspecified        CRLReason = "Unspecified"
)
