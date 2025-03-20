package ratchet

import (
	"fmt"

	"github.com/rylenko/bastion/pkg/ratchet/errors"
	"github.com/rylenko/bastion/pkg/ratchet/receivingchain"
	"github.com/rylenko/bastion/pkg/ratchet/rootchain"
	"github.com/rylenko/bastion/pkg/ratchet/sendingchain"
)

type config struct {
	crypto           Crypto
	receivingOptions []receivingchain.Option
	rootOptions      []rootchain.Option
	sendingOptions   []sendingchain.Option
}

func newConfig(options []Option) (config, error) {
	cfg := config{crypto: newDefaultCrypto()}
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
