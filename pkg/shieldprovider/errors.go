package shieldprovider

import "errors"

var (
	ErrInvalidValue  = errors.New("invalid value")
	ErrDiffieHellman = errors.New("Diffie-Hellman")
)
