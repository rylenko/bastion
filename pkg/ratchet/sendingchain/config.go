package sendingchain

import "github.com/rylenko/bastion/pkg/ratchet/messagechaincommon"

type Config struct {
	crypto Crypto
}

func NewConfig(options ...ConfigOption) *Config {
	config := &Config{crypto: messagechaincommon.NewCrypto()}
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
