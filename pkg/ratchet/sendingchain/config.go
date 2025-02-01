package sendingchain

type Config struct {
	crypto Crypto
}

func NewConfig(options ...ConfigOption) *Config {
	config := &Config{crypto: newCrypto()}
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
