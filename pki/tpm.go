package pki

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/asn1"
	"fmt"
	"hash"
	"math/big"
	"os"
	"path"

	"github.com/feeltheajf/ztca/errdefs"
)

var (
	tpmTrustedRoots []*x509.Certificate
)

func init() {
	roots := "testdata/tpm/bundle"
	files, _ := os.ReadDir(roots)
	for _, f := range files {
		crt, _ := ReadCertificate(path.Join(roots, f.Name()))
		tpmTrustedRoots = append(tpmTrustedRoots, crt)
	}

	ek, _ := ReadCertificate("testdata/tpm/ek.crt")
	if err := VerifyAttestation(ek); err != nil {
		fmt.Printf("%+v\n\n", err)
	}
}

type ecdsaSignature struct {
	R, S *big.Int
}

// VerifyAttestation verifies the given certificate chain against
// a list of trusted roots from known TPM manufacturers
func VerifyAttestation(crt *x509.Certificate) error {
	for _, root := range tpmTrustedRoots {
		var h hash.Hash
		var hashAlgorithm crypto.Hash
		switch crt.SignatureAlgorithm {
		case x509.SHA256WithRSA, x509.ECDSAWithSHA256:
			h = sha256.New()
			hashAlgorithm = crypto.SHA256
		default:
			return errdefs.InvalidParameter("unsupported signature algorithm: '%s'", crt.SignatureAlgorithm)
		}
		h.Write(crt.RawTBSCertificate)
		hsum := h.Sum(nil)

		switch pub := root.PublicKey.(type) {
		case *rsa.PublicKey:
			err := rsa.VerifyPKCS1v15(pub, hashAlgorithm, hsum, crt.Signature)
			if err == nil {
				return nil
			}
		case *ecdsa.PublicKey:
			// taken from
			// https://github.com/codelittinc/gobitauth/blob/master/sign.go
			raw := asn1.RawValue{}
			if _, err := asn1.Unmarshal(crt.Signature, &raw); err != nil {
				panic(err.Error())
			}
			// The format of DER string is 0x02 + rlen + r + 0x02 + slen + s
			rLen := raw.Bytes[1] // The entire length of R + offset of 2 for 0x02 and rlen
			r := big.NewInt(0).SetBytes(raw.Bytes[2 : rLen+2])
			// Ignore the next 0x02 and slen bytes and just take the start of S to the end of the byte array
			s := big.NewInt(0).SetBytes(raw.Bytes[rLen+4:])

			ok := ecdsa.Verify(pub, hsum, r, s)
			if ok {
				return nil
			}
		default:
			return errdefs.InvalidParameter("unknown public key type: '%T'", root.PublicKey)
		}
	}

	return errdefs.InvalidParameter("no matching root certificate")
}
