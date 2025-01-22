package shieldprovider

import "errors"

var (
	ErrInvalidValue         = errors.New("invalid value")
	ErrConvertToForeignType = errors.New("convert to foreign type")
	ErrDiffieHellman        = errors.New("Diffie-Hellman")
)
