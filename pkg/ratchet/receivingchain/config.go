package receivingchain

import (
	"fmt"

	"github.com/rylenko/bastion/pkg/ratchet/errors"
)

type config struct {
	crypto             Crypto
	skippedKeysStorage SkippedKeysStorage
}

func newConfig(options []Option) (config, error) {
	cfg := config{
		crypto:             newCrypto(),
		skippedKeysStorage: newSkippedKeysStorage(),
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

func (cfg config) clone() config {
	cfg.skippedKeysStorage = cfg.skippedKeysStorage.Clone()
	return cfg
}
