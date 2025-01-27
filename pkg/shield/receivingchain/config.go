package receivingchain

const messageKeysSkipLimit = 1024

type Config struct {
	crypto               Crypto
	messageKeysSkipLimit uint32
	skippedMessageKeys   SkippedMessageKeys
}

func NewConfig(options ...Option) *Config {
	config := &Config{
		crypto:               newCrypto(),
		messageKeysSkipLimit: messageKeysSkipLimit,
		skippedMessageKeys:   newSkippedMessageKeys(),
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

func WithMessageKeysSkipLimit(limit uint32) Option {
	return func(config *Config) {
		config.messageKeysSkipLimit = limit
	}
}

func WithSkippedMessageKeys(storage SkippedMessageKeys) Option {
	return func(config *Config) {
		config.skippedMessageKeys = storage
	}
}
