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
	config := &config{
		crypto:               newCrypto(),
		receivingChainConfig: receivingchain.NewConfig(),
		rootChainConfig:      rootchain.NewConfig(),
		sendingChainConfig:   sendingchain.NewConfig(),
	}

	for _, option := range options {
		option(config)
	}

	return config
}

type ConfigOption func(config *config)

func WithCrypto(crypto Crypto) ConfigOption {
	return func(config *config) {
		config.crypto = crypto
	}
}

func WithMessageKeysSkipLimit(limit uint32) ConfigOption {
	return func(config *config) {
		config.receivingChainConfig.ApplyOptions(receivingchain.WithMessageKeysSkipLimit(limit))
	}
}

func WithReceivingChainCrypto(crypto receivingchain.Crypto) ConfigOption {
	return func(config *config) {
		config.receivingChainConfig.ApplyOptions(receivingchain.WithCrypto(crypto))
	}
}

func WithRootChainCrypto(crypto rootchain.Crypto) ConfigOption {
	return func(config *config) {
		config.rootChainConfig.ApplyOptions(rootchain.WithCrypto(crypto))
	}
}

func WithSendingChainCrypto(crypto sendingchain.Crypto) ConfigOption {
	return func(config *config) {
		config.sendingChainConfig.ApplyOptions(sendingchain.WithCrypto(crypto))
	}
}

func WithSkippedMessageKeys(keys receivingchain.SkippedMessageKeys) ConfigOption {
	return func(config *config) {
		config.receivingChainConfig.ApplyOptions(receivingchain.WithSkippedMessageKeys(keys))
	}
}
