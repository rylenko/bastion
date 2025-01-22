package shield

import "github.com/rylenko/sapphire/pkg/shield/keys"

// Provider provides access to the cryptographic part of the shield.
type Provider interface {
	// GeneratePrivateKey must generate a cryptographically secure private key.
	GeneratePrivateKey() (*keys.Private, error)
}
