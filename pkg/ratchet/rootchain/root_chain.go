package rootchain

import (
	"fmt"

	"github.com/rylenko/bastion/pkg/ratchet/errors"
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

func (rc *RootChain) Advance(sharedSecretKey *keys.SharedSecret) (*keys.MessageMaster, *keys.Header, error) {
	if rc.config == nil {
		return nil, nil, fmt.Errorf("%w: config is nil", errors.ErrInvalidValue)
	}

	if rc.config.crypto == nil {
		return nil, nil, fmt.Errorf("%w: config crypto is nil", errors.ErrInvalidValue)
	}

	newRootKey, messageMasterKey, nextHeaderKey, err := rc.config.crypto.AdvanceChain(rc.key, sharedSecretKey)
	if err != nil {
		return nil, nil, fmt.Errorf("%w: advance: %w", errors.ErrCrypto, err)
	}

	rc.key = newRootKey

	return messageMasterKey, nextHeaderKey, nil
}

func (rc *RootChain) Clone() *RootChain {
	return New(rc.key.Clone(), rc.config)
}
