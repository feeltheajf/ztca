package pki

import (
	"fmt"
	"math/big"
	"strings"
)

const (
	defaultSerialBase = 16
)

// UnmarshalSerial converts string to certificate serial
func UnmarshalSerial(serial string) *big.Int {
	srl := new(big.Int)
	srl.SetString(serial, defaultSerialBase)
	return srl
}

// MarshalSerial converts certificate serial to string
func MarshalSerial(serial *big.Int) string {
	return strings.ToUpper(serial.Text(defaultSerialBase))
}

// MarshalYubiKeySerial converts YubiKey serial to string
func MarshalYubiKeySerial(serial uint32) string {
	return fmt.Sprintf("%d", serial)
}
