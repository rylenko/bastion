package sendingchain

type Config struct {
	crypto Crypto
}

func NewConfig(options ...Option) *Config {
	config := &Config{crypto: newCrypto()}

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
