package shield

import "errors"

var (
	ErrAdvanceChain  = errors.New("advance chain")
	ErrCrypto        = errors.New("crypto")
	ErrDiffieHellman = errors.New("diffie hellman")
	ErrForeignType   = errors.New("foreign type")
	ErrInvalidValue  = errors.New("invalid value")
)
