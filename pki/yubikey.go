package pki

import (
	"encoding/asn1"
)

var (
	oidYubiKeySerial = asn1.ObjectIdentifier{1, 3, 6, 1, 4, 1, 41482, 3, 7}
)

type YubiKeyExtensions struct {
	//
}

//

// TODO convert to get YubiKeyExtensions(crt *x509.Certificate)
// GetYubikeySerial extracts Yubikey serial from the given attestation statement
// func GetYubikeySerial(att *x509.Certificate) (int64, error) {
// 	for _, e := range att.Extensions {
// 		if !e.Id.Equal(oidYubiKeySerial) {
// 			continue
// 		}
// 		var serial int64
// 		_, err := asn1.Unmarshal(e.Value, &serial)
// 		if err != nil {
// 			return 0, errdefs.InvalidParameter("failed to unmarshal Yubikey serial").CausedBy(err)
// 		}
// 		if serial < 0 {
// 			return 0, errdefs.InvalidParameter("negative Yubikey serial: %d", serial)
// 		}
// 		return serial, nil
// 	}
// 	return 0, errdefs.InvalidParameter("missing Yubikey serial extension")
// }
