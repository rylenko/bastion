package ratchet

import (
	"fmt"

	"github.com/rylenko/bastion/pkg/ratchet/errors"
	"github.com/rylenko/bastion/pkg/ratchet/receivingchain"
	"github.com/rylenko/bastion/pkg/ratchet/rootchain"
	"github.com/rylenko/bastion/pkg/ratchet/sendingchain"
	"github.com/rylenko/bastion/pkg/ratchet/utils"
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

func WithMessageKeysSkipLimit(limit uint64) Option {
	return func(cfg *config) error {
		cfg.receivingOptions = append(cfg.receivingOptions, receivingchain.WithMessageKeysSkipLimit(limit))
		return nil
	}
}

func WithReceivingChainCrypto(crypto receivingchain.Crypto) Option {
	return func(cfg *config) error {
		cfg.receivingOptions = append(cfg.receivingOptions, receivingchain.WithCrypto(crypto))
		return nil
	}
}

func WithRootChainCrypto(crypto rootchain.Crypto) Option {
	return func(cfg *config) error {
		cfg.rootOptions = append(cfg.rootOptions, rootchain.WithCrypto(crypto))
		return nil
	}
}

func WithSendingChainCrypto(crypto sendingchain.Crypto) Option {
	return func(cfg *config) error {
		cfg.sendingOptions = append(cfg.sendingOptions, sendingchain.WithCrypto(crypto))
		return nil
	}
}

func WithSkippedMessageKeysStorage(storage receivingchain.SkippedMessageKeysStorage) Option {
	return func(cfg *config) error {
		cfg.receivingOptions = append(cfg.receivingOptions, receivingchain.WithSkippedMessageKeysStorage(storage))
		return nil
	}
}
