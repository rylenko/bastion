package ratchet

import (
	"github.com/rylenko/bastion/pkg/ratchet/receivingchain"
	"github.com/rylenko/bastion/pkg/ratchet/rootchain"
	"github.com/rylenko/bastion/pkg/ratchet/sendingchain"
)

type config struct {
	crypto               Crypto
	receivingChainConfig *receivingchain.Config
	rootChainConfig      *rootchain.Config
	sendingChainConfig   *sendingchain.Config
}

func newConfig(options ...ConfigOption) *config {
	cfg := &config{
		crypto:               newCrypto(),
		receivingChainConfig: receivingchain.NewConfig(),
		rootChainConfig:      rootchain.NewConfig(),
		sendingChainConfig:   sendingchain.NewConfig(),
	}

	for _, option := range options {
		option(cfg)
	}

	return cfg
}

type ConfigOption func(cfg *config)

func WithCrypto(crypto Crypto) ConfigOption {
	return func(cfg *config) {
		cfg.crypto = crypto
	}
}

func WithMessageKeysSkipLimit(limit uint32) ConfigOption {
	return func(cfg *config) {
		cfg.receivingChainConfig.ApplyOptions(receivingchain.WithMessageKeysSkipLimit(limit))
	}
}

func WithReceivingChainCrypto(crypto receivingchain.Crypto) ConfigOption {
	return func(cfg *config) {
		cfg.receivingChainConfig.ApplyOptions(receivingchain.WithCrypto(crypto))
	}
}

func WithRootChainCrypto(crypto rootchain.Crypto) ConfigOption {
	return func(cfg *config) {
		cfg.rootChainConfig.ApplyOptions(rootchain.WithCrypto(crypto))
	}
}

func WithSendingChainCrypto(crypto sendingchain.Crypto) ConfigOption {
	return func(cfg *config) {
		cfg.sendingChainConfig.ApplyOptions(sendingchain.WithCrypto(crypto))
	}
}

func WithSkippedMessageKeysStorage(keys receivingchain.SkippedMessageKeysStorage) ConfigOption {
	return func(cfg *config) {
		cfg.receivingChainConfig.ApplyOptions(receivingchain.WithSkippedMessageKeysStorage(keys))
	}
}
