package chains

import (
	"fmt"

	"github.com/rylenko/sapphire/pkg/shield/keys"
)

// Root is the root chain of the participant of conversation. The main task of the root chain is to generate new keys
// of the sending and receiving chains and their next header encryption keys.
type Root struct {
	key *keys.Root
}

// NewRoot creates a new instance of root chain.
func NewRoot(key *keys.Root) *Root {
	return &Root{key: key}
}

// Forward moves the root chain forward. In other words, creating a new root key, a new message master key for the
// sending or receiving chain, and the next header encryption key.
//
// This method is a wrapper around the provider that sets a new root key into the current chain.
func (chain *Root) Forward(
	provider Provider,
	sharedSecretKey *keys.SharedSecret,
) (*keys.MessageMaster, *keys.Header, error) {
	if provider == nil {
		return nil, nil, fmt.Errorf("%w: provider is nil", ErrInvalidValue)
	}

	newRootKey, messageMasterKey, nextHeaderKey, err := provider.ForwardRootChain(chain.key, sharedSecretKey)
	if err != nil {
		return nil, nil, fmt.Errorf("%w: forward root chain: %w", ErrProvider, err)
	}

	chain.key = newRootKey

	return messageMasterKey, nextHeaderKey, nil
}
