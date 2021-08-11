package api

import (
	"crypto/x509/pkix"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-piv/piv-go/piv"
	"github.com/rs/zerolog/log"

	"github.com/feeltheajf/ztca/dto"
	"github.com/feeltheajf/ztca/errdefs"
	"github.com/feeltheajf/ztca/pki"
)

type YubiKeyRequest struct {
	Username string `form:"username" binding:"required"`
	Att      string `form:"att" binding:"required"`
	IntAtt   string `form:"intAtt" binding:"required"`
}

func yubikey(c *gin.Context) {
	q := new(YubiKeyRequest)
	if err := c.Bind(q); err != nil {
		handle(c, errdefs.InvalidParameter(err))
		return
	}

	att, err := pki.UnmarshalCertificate(q.Att)
	if err != nil {
		handle(c, errdefs.InvalidParameter("bad attestation certificate").CausedBy(err))
		return
	}

	intAtt, err := pki.UnmarshalCertificate(q.IntAtt)
	if err != nil {
		handle(c, errdefs.InvalidParameter("bad intermediate attestation certificate").CausedBy(err))
		return
	}

	meta, err := piv.Verify(intAtt, att)
	if err != nil {
		handle(c, errdefs.InvalidParameter("failed to verify attestation").CausedBy(err))
		return
	}

	ykSerial := dto.FormatYubiKeySerial(meta.Serial)
	t, err := pki.NewTemplate(pki.WithName(pkix.Name{
		CommonName:   q.Username,
		SerialNumber: ykSerial,
	}))
	if err != nil {
		handle(c, errdefs.Unknown("failed to create certificate template").CausedBy(err))
		return
	}

	crt, err := pki.NewCertificate(t, att.PublicKey)
	if err != nil {
		handle(c, errdefs.Unknown("failed to issue certificate").CausedBy(err))
		return
	}
	serial := fmt.Sprintf("%d", crt.SerialNumber)

	raw, err := pki.MarshalCertificate(crt)
	if err != nil {
		handle(c, errdefs.Unknown("failed to marshal certificate").CausedBy(err))
		return
	}

	log.Info().
		Str("serial", serial).
		Str("username", q.Username).
		Str("device_serial", ykSerial).
		Msg("new certificate")

	c.String(http.StatusCreated, raw)
}
