package rootchain

import (
	"fmt"

	"github.com/rylenko/bastion/pkg/ratchet/keys"
)

type RootChain struct {
	key    *keys.Root
	config *Config
}

func New(key *keys.Root, config *Config) *RootChain {
	return &RootChain{
		key:    key,
		config: config,
	}
}

func (chain *RootChain) Advance(sharedSecretKey *keys.SharedSecret) (*keys.MessageMaster, *keys.Header, error) {
	if chain.config == nil {
		return nil, nil, fmt.Errorf("%w: config is nil", ErrInvalidValue)
	}

	if chain.config.crypto == nil {
		return nil, nil, fmt.Errorf("%w: config crypto is nil", ErrInvalidValue)
	}

	newRootKey, messageMasterKey, nextHeaderKey, err := chain.config.crypto.AdvanceChain(chain.key, sharedSecretKey)
	if err != nil {
		return nil, nil, fmt.Errorf("%w: advance via crypto: %w", ErrCrypto, err)
	}

	chain.key = newRootKey

	return messageMasterKey, nextHeaderKey, nil
}
