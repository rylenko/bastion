package shieldprovider

import "crypto/ecdh"

// Provider is a crypto provider for shield module.
type Provider struct {
	curve ecdh.Curve
}

// New creates a new instance of provider.
func New() *Provider {
	return &Provider{curve: ecdh.X25519()}
}
