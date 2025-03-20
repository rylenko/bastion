package receivingchain

import (
	"fmt"

	"github.com/rylenko/bastion/pkg/ratchet/errlist"
)

type config struct {
	crypto             Crypto
	skippedKeysStorage SkippedKeysStorage
}

func newConfig(options []Option) (config, error) {
	cfg := config{
		crypto:             newDefaultCrypto(),
		skippedKeysStorage: newDefaultSkippedKeysStorage(),
	}
	if err := cfg.applyOptions(options); err != nil {
		return config{}, fmt.Errorf("%w: %w", errlist.ErrOption, err)
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
