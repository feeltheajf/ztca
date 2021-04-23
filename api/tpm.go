package api

import (
	"crypto/subtle"
	"encoding/hex"
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/google/go-attestation/attest"
	"github.com/rs/zerolog/log"

	"github.com/feeltheajf/ztca/errdefs"
	"github.com/feeltheajf/ztca/x/rand"
)

var (
	// TODO replace with proper cache implementation
	cache = map[string]string{}
)

type nonceResponse struct {
	Nonce string `json:"nonce"`
}

func nonce(c *gin.Context) {
	c.JSON(http.StatusOK, &nonceResponse{rand.String()})
}

type generateRequest struct {
	// Only TPM 2.0 supported
	TPMVersion attest.TPMVersion `json:"version" binding:"required"`
	// Public key burned in TPM
	EndorsementKey []byte `json:"ek" binding:"required"`
	// TODO
	AttestationParameters attest.AttestationParameters `json:"params" binding:"required"`
}

type generateResponse struct {
	// Encrypted challenge
	Credential []byte `json:"credential"`
}

func generate(c *gin.Context) {
	q := new(generateRequest)
	if err := c.Bind(q); err != nil {
		handle(c, errdefs.InvalidParameter(err))
		return
	}

	ek, err := jwt.ParseRSAPublicKeyFromPEM(q.EndorsementKey)
	if err != nil {
		handle(c, errdefs.InvalidParameter("bad endorsement key").CausedBy(err))
		return
	}

	// TODO make sure that EK is in TPM
	// by comparing public keys?
	// by checking some signature?
	// research
	ap := attest.ActivationParameters{
		TPMVersion: q.TPMVersion,
		EK:         ek,
		AK:         q.AttestationParameters,
	}

	secret, ec, err := ap.Generate()
	if err != nil {
		handle(c, errdefs.Unknown("failed to generate activation challenge").CausedBy(err))
		return
	}

	secretString := hex.EncodeToString(secret)
	credential := hex.EncodeToString(ec.Credential)
	cache[credential] = secretString

	c.JSON(http.StatusCreated, &generateResponse{ec.Credential})
}

type activateRequest struct {
	// Encrypted challenge
	Credential []byte `json:"credential" binding:"required"`
	// Decrypted challenge
	Secret []byte `json:"secret" binding:"required"`
}

func activate(c *gin.Context) {
	q := new(activateRequest)
	if err := c.Bind(q); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	original, ok := cache[hex.EncodeToString(q.Credential)]
	if !ok {
		c.AbortWithStatusJSON(http.StatusBadRequest, "no such credential")
		return
	}

	received := hex.EncodeToString(q.Secret)
	log.Printf("original: %s", original)
	log.Printf("received: %s", received)

	if subtle.ConstantTimeCompare([]byte(original), []byte(received)) == 0 {
		handle(c, errdefs.InvalidParameter("challenge did not match"))
		return
	}

	// TODO get public key from cache and issue certificate
	// issueCertificate(c, att.PublicKey)
}
