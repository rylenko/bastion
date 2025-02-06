package receivingchain

import "github.com/rylenko/bastion/pkg/ratchet/messagechaincommon"

const messageKeysSkipLimit = 1024

type Config struct {
	crypto                    Crypto
	messageKeysSkipLimit      uint64
	skippedMessageKeysStorage SkippedMessageKeysStorage
}

func NewConfig(options ...ConfigOption) *Config {
	cfg := &Config{
		crypto:                    messagechaincommon.NewCrypto(),
		messageKeysSkipLimit:      messageKeysSkipLimit,
		skippedMessageKeysStorage: newSkippedMessageKeysStorage(),
	}
	cfg.ApplyOptions(options...)

	return cfg
}

func (cfg *Config) ApplyOptions(options ...ConfigOption) {
	for _, option := range options {
		option(cfg)
	}
}

type ConfigOption func(cfg *Config)

func WithCrypto(crypto Crypto) ConfigOption {
	return func(c *Config) {
		c.crypto = crypto
	}
}

func WithMessageKeysSkipLimit(limit uint64) ConfigOption {
	return func(cfg *Config) {
		cfg.messageKeysSkipLimit = limit
	}
}

func WithSkippedMessageKeysStorage(storage SkippedMessageKeysStorage) ConfigOption {
	return func(cfg *Config) {
		cfg.skippedMessageKeysStorage = storage
	}
}
