package receivingchain

import "github.com/rylenko/bastion/pkg/ratchet/messagechaincommon"

const messageKeysSkipLimit = 1024

type Config struct {
	crypto                    Crypto
	messageKeysSkipLimit      uint32
	skippedMessageKeysStorage SkippedMessageKeysStorage
}

func NewConfig(options ...ConfigOption) *Config {
	config := &Config{
		crypto:                    messagechaincommon.NewCrypto(),
		messageKeysSkipLimit:      messageKeysSkipLimit,
		skippedMessageKeysStorage: newSkippedMessageKeysStorage(),
	}
	config.ApplyOptions(options...)

	return config
}

func (config *Config) ApplyOptions(options ...ConfigOption) {
	for _, option := range options {
		option(config)
	}
}

type ConfigOption func(config *Config)

func WithCrypto(crypto Crypto) ConfigOption {
	return func(config *Config) {
		config.crypto = crypto
	}
}

func WithMessageKeysSkipLimit(limit uint32) ConfigOption {
	return func(config *Config) {
		config.messageKeysSkipLimit = limit
	}
}

func WithSkippedMessageKeysStorage(storage SkippedMessageKeysStorage) ConfigOption {
	return func(config *Config) {
		config.skippedMessageKeysStorage = storage
	}
}
