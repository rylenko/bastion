package receivingchain

import (
	"fmt"

	"github.com/rylenko/bastion/pkg/ratchet/errors"
	"github.com/rylenko/bastion/pkg/ratchet/messagechaincommon"
)

const messageKeysSkipLimit = 1024

type config struct {
	crypto                    Crypto
	messageKeysSkipLimit      uint64
	skippedMessageKeysStorage SkippedMessageKeysStorage
}

func newConfig(options []Option) (config, error) {
	cfg := config{
		crypto:                    messagechaincommon.NewCrypto(),
		messageKeysSkipLimit:      messageKeysSkipLimit,
		skippedMessageKeysStorage: newSkippedMessageKeysStorage(),
	}
	if err := cfg.applyOptions(options); err != nil {
		return config{}, fmt.Errorf("%w: %w", errors.ErrOption, err)
	}

	return cfg, nil
}

func (cfg *config) applyOptions(options []Option) error {
	for _, option := range options {
		if err := option(cfg); err != nil {
			return err
		}
	}

	return nil
}
