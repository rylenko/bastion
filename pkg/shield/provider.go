package shield

import "github.com/rylenko/sapphire/pkg/shield/keys"

// Provider provides access to the cryptographic part of the shield.
type Provider interface {
	// ComputeSharedSecretKey must compute a shared secret key based on the private and public keys.
	//
	// For example, this could be the Diffie-Hellman algorithm.
	ComputeSharedSecretKey(privateKey *keys.Private, publicKey *keys.Public) (*keys.SharedSecret, error)

	// ForwardRootChain moves the root chain forward. In other words, creating a new root key, a new master key for the
	// sending or receiving chain, and the next header encryption key.
	ForwardRootChain(
		rootKey *keys.Root,
		sharedSecretKey *keys.SharedSecret,
	) (*keys.Root, *keys.Master, *keys.Header, error)

	// GeneratePrivateKey must generate a cryptographically secure private key.
	GeneratePrivateKey() (*keys.Private, error)
}
