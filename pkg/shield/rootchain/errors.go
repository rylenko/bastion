package rootchain

import "errors"

var (
	ErrCrypto       = errors.New("crypto")
	ErrInvalidValue = errors.New("invalid value")
	ErrKDF          = errors.New("KDF")
	ErrNewHash      = errors.New("new hash")
)
