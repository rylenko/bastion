package shieldprovider

import "errors"

var (
	ErrConvertToForeignType = errors.New("convert to foreign type")
	ErrDiffieHellman        = errors.New("Diffie-Hellman")
	ErrInvalidValue         = errors.New("invalid value")
)
