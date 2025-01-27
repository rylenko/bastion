package shield

import (
	"github.com/rylenko/sapphire/pkg/shield/receivingchain"
	"github.com/rylenko/sapphire/pkg/shield/rootchain"
	"github.com/rylenko/sapphire/pkg/shield/sendingchain"
)

type Config struct {
	crypto               Crypto
	receivingChainConfig *receivingchain.Config
	rootChainConfig      *rootchain.Config
	sendingChainConfig   *sendingchain.Config
}

func NewConfig(options ...Option) *Config {
	config := &Config{
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

type Option func(config *Config)

func WithCrypto(crypto Crypto) Option {
	return func(config *Config) {
		config.crypto = crypto
	}
}

func WithReceivingChainConfig(receivingChainConfig *receivingchain.Config) Option {
	return func(config *Config) {
		config.receivingChainConfig = receivingChainConfig
	}
}

func WithRootChainConfig(rootChainConfig *rootchain.Config) Option {
	return func(config *Config) {
		config.rootChainConfig = rootChainConfig
	}
}

func WithSendingChainConfig(sendingChainConfig *sendingchain.Config) Option {
	return func(config *Config) {
		config.sendingChainConfig = sendingChainConfig
	}
}
