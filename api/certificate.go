package api

import (
	"crypto/x509/pkix"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-piv/piv-go/piv"

	"github.com/feeltheajf/ztca/dto"
	"github.com/feeltheajf/ztca/errdefs"
	"github.com/feeltheajf/ztca/pki"
)

type CertificateRequest struct {
	IntAtt string `form:"intAtt" binding:"required"`
	Att    string `form:"att" binding:"required"`
}

func issueCertificate(c *gin.Context) {
	q := new(CertificateRequest)
	if err := c.Bind(q); err != nil {
		handle(c, errdefs.InvalidParameter(err))
		return
	}

	username := c.Param("username")
	_, err := dto.Certificates.Load(username)
	if err == nil {
		handle(c, errdefs.Conflict("certificate already issued"))
		return
	}

	intAtt, err := pki.UnmarshalCertificate(q.IntAtt)
	if err != nil {
		handle(c, errdefs.InvalidParameter("bad intermediate attestation certificate").CausedBy(err))
		return
	}

	att, err := pki.UnmarshalCertificate(q.Att)
	if err != nil {
		handle(c, errdefs.InvalidParameter("bad attestation certificate").CausedBy(err))
		return
	}

	meta, err := piv.Verify(intAtt, att)
	if err != nil {
		handle(c, errdefs.InvalidParameter("failed to verify attestation").CausedBy(err))
		return
	}

	template, err := pki.NewTemplate(pki.WithName(pkix.Name{
		CommonName:   username,
		SerialNumber: pki.MarshalYubiKeySerial(meta.Serial),
	}))
	if err != nil {
		handle(c, errdefs.Unknown("failed to create certificate template").CausedBy(err))
		return
	}

	x509, err := pki.NewCertificate(template, att.PublicKey)
	if err != nil {
		handle(c, errdefs.Unknown("failed to issue certificate").CausedBy(err))
		return
	}

	raw, err := pki.MarshalCertificate(x509)
	if err != nil {
		handle(c, errdefs.Unknown("failed to marshal certificate").CausedBy(err))
		return
	}

	crt := &dto.Certificate{
		Raw:          raw,
		SerialNumber: pki.MarshalSerial(x509.SerialNumber),
		ExpiresAt:    x509.NotAfter,
		Username:     x509.Subject.CommonName,
		DeviceSerial: x509.Subject.SerialNumber,
	}

	if err := dto.Certificates.Save(crt); err != nil {
		handle(c, errdefs.Unknown("failed to save certificate").CausedBy(err))
		return
	}

	c.String(http.StatusCreated, crt.Raw)
}

type RevokeRequest struct {
	Reason pki.CRLReason `form:"reason" binding:"omitempty,oneof=KeyCompromise AffiliationChanged Superseded Unspecified"`
}

func revokeCertificate(c *gin.Context) {
	q := new(RevokeRequest)
	if err := c.Bind(q); err != nil {
		handle(c, errdefs.InvalidParameter(err))
		return
	}

	username := c.Param("username")
	crt, err := dto.Certificates.Load(username)
	if err == nil {
		if err := pki.Revoke(crt.X509(), q.Reason, time.Now()); err != nil {
			handle(c, errdefs.Unknown("failed to revoke certificate").CausedBy(err))
			return
		}

		if err := dto.Certificates.Delete(crt); err != nil {
			handle(c, errdefs.Unknown("failed to delete certificate").CausedBy(err))
			return
		}
	}

	c.String(http.StatusNoContent, "")
}
