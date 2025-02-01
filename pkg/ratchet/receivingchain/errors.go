package receivingchain

import "errors"

var (
	ErrCrypto       = errors.New("crypto")
	ErrInvalidValue = errors.New("invalid value")
	ErrMAC          = errors.New("MAC")
	ErrNewHash      = errors.New("new hash")
)
