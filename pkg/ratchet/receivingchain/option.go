package receivingchain

import (
	"fmt"

	"github.com/rylenko/bastion/pkg/ratchet/errors"
	"github.com/rylenko/bastion/pkg/utils"
)

type Option func(cfg *config) error

func WithCrypto(crypto Crypto) Option {
	return func(cfg *config) error {
		if utils.IsNil(crypto) {
			return fmt.Errorf("%w: crypto is nil", errors.ErrInvalidValue)
		}

		cfg.crypto = crypto

		return nil
	}
}

func WithSkippedKeysStorage(storage SkippedKeysStorage) Option {
	return func(cfg *config) error {
		if utils.IsNil(storage) {
			return fmt.Errorf("%w: storage is nil", errors.ErrInvalidValue)
		}

		cfg.skippedKeysStorage = storage

		return nil
	}
}
