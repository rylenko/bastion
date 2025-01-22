package shield

import "github.com/rylenko/sapphire/pkg/shield/keys"

// Provider provides access to the cryptographic part of the shield.
type Provider interface {
	// ComputeSharedSecretKey must compute a shared secret key based on the private and public keys. For example, this
	// could be the Diffie-Hellman algorithm.
	ComputeSharedSecretKey(privateKey *keys.Private, publicKey *keys.Public) (*keys.SharedSecret, error)

	// GeneratePrivateKey must generate a cryptographically secure private key.
	GeneratePrivateKey() (*keys.Private, error)
}
