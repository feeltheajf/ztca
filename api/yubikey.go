package api

import (
	"github.com/gin-gonic/gin"
	"github.com/go-piv/piv-go/piv"

	"github.com/feeltheajf/ztca/errdefs"
	"github.com/feeltheajf/ztca/pki"
)

type yubikeyRequest struct {
	Att    []byte `json:"att" binding:"required"`
	IntAtt []byte `json:"intAtt" binding:"required"`
}

func yubikey(c *gin.Context) {
	q := new(yubikeyRequest)
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

	if _, err := piv.Verify(intAtt, att); err != nil {
		handle(c, errdefs.InvalidParameter("failed to verify attestation").CausedBy(err))
		return
	}

	issueCertificate(c, att.PublicKey)
}
