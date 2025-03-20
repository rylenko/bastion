package sendingchain

import (
	"fmt"

	"github.com/rylenko/bastion/pkg/ratchet/errlist"
	"github.com/rylenko/bastion/pkg/utils"
)

type Option func(cfg *config) error

func WithCrypto(crypto Crypto) Option {
	return func(cfg *config) error {
		if utils.IsNil(crypto) {
			return fmt.Errorf("%w: crypto is nil", errlist.ErrInvalidValue)
		}

		cfg.crypto = crypto

		return nil
	}
}
