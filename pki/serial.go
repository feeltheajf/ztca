package pki

import (
	"fmt"
	"math/big"
	"strings"

	"github.com/feeltheajf/ztca/fs"
)

const (
	defaultSerialBase = 16
)

// ReadSerial loads serial from file
func ReadSerial(filename string) (*big.Int, error) {
	raw, err := fs.Read(filename)
	if err != nil {
		return nil, err
	}
	return UnmarshalSerial(raw), nil
}

// UnmarshalSerial converts string to certificate serial
func UnmarshalSerial(serial string) *big.Int {
	srl := new(big.Int)
	srl.SetString(serial, defaultSerialBase)
	return srl
}

// WriteSerial saves serial to file
func WriteSerial(filename string, serial *big.Int) error {
	return fs.Write(filename, MarshalSerial(serial))
}

// MarshalSerial converts certificate serial to string
func MarshalSerial(serial *big.Int) string {
	return strings.ToUpper(serial.Text(defaultSerialBase))
}

// MarshalYubiKeySerial converts YubiKey serial to string
func MarshalYubiKeySerial(serial uint32) string {
	return fmt.Sprintf("%d", serial)
}
