package api

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-piv/piv-go/piv"
	"github.com/rs/zerolog/log"

	"github.com/feeltheajf/ztca/pki"
)

type certificateRequest struct {
	Att    []byte `json:"att"`
	IntAtt []byte `json:"intAtt"`
}

type certificateResponse struct {
	Crt []byte `json:"crt"`
}

func requestYubiKeyCertificate(c *gin.Context) {
	q := new(certificateRequest)
	if err := c.Bind(q); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
	}

	att, err := pki.UnmarshalCertificate(q.Att)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
	}

	intAtt, err := pki.UnmarshalCertificate(q.IntAtt)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
	}

	if _, err := piv.Verify(intAtt, att); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
	}

	s := getSession(c)
	t, err := pki.NewTemplate(pki.WithCommonName(s.User.Name))
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
	}

	crt, err := pki.NewCertificate(t, att.PublicKey)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
	}

	b, err := pki.MarshalCertificate(crt)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
	}

	log.Info().
		Str("serial", fmt.Sprintf("%d", crt.SerialNumber)).
		Str("user", s.User.Name).
		Msg("new certificate")

	c.JSON(http.StatusCreated, &certificateResponse{b})
}
