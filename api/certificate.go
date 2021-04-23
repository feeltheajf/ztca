package api

import (
	"crypto"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"

	"github.com/feeltheajf/ztca/errdefs"
	"github.com/feeltheajf/ztca/pki"
)

type certificateResponse struct {
	Crt []byte `json:"crt"`
}

func issueCertificate(c *gin.Context, pub crypto.PublicKey) {
	s := getSession(c)
	t, err := pki.NewTemplate(pki.WithCommonName(s.User.Name))
	if err != nil {
		handle(c, errdefs.Unknown("failed to create certificate template").CausedBy(err))
		return
	}

	crt, err := pki.NewCertificate(t, pub)
	if err != nil {
		handle(c, errdefs.Unknown("failed to create new certificate").CausedBy(err))
		return
	}

	b, err := pki.MarshalCertificate(crt)
	if err != nil {
		handle(c, errdefs.Unknown("failed to marshal certificate").CausedBy(err))
		return
	}

	log.Info().
		Str("serial", fmt.Sprintf("%d", crt.SerialNumber)).
		Str("user", s.User.Name).
		Msg("new certificate")

	c.JSON(http.StatusCreated, &certificateResponse{b})
}
