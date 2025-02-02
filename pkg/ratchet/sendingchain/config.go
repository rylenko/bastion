package sendingchain

import "github.com/rylenko/bastion/pkg/ratchet/messagechaincommon"

type Config struct {
	crypto Crypto
}

func NewConfig(options ...ConfigOption) *Config {
	cfg := &Config{crypto: messagechaincommon.NewCrypto()}
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
	return func(cfg *Config) {
		cfg.crypto = crypto
	}
}
