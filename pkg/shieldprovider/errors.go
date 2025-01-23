package shieldprovider

import "errors"

var (
	ErrDiffieHellman      = errors.New("Diffie-Hellman")
	ErrInvalidValue       = errors.New("invalid value")
	ErrKDF                = errors.New("KDF")
	ErrMapFromForeignType = errors.New("map from foreign type")
	ErrMapToForeignType   = errors.New("map to foreign type")
	ErrNewHash            = errors.New("new hash")
)
