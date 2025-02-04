package ratchet

import "errors"

var (
	ErrCrypto       = errors.New("crypto")
	ErrInvalidValue = errors.New("invalid value")
)
