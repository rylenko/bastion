package receivingchain

const messageKeysSkipLimit = 1024

type Config struct {
	crypto               Crypto
	messageKeysSkipLimit uint32
	skippedMessageKeys   SkippedMessageKeys
}

func NewConfig(options ...ConfigOption) *Config {
	config := &Config{
		crypto:               newCrypto(),
		messageKeysSkipLimit: messageKeysSkipLimit,
		skippedMessageKeys:   newSkippedMessageKeys(),
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

func WithSkippedMessageKeys(storage SkippedMessageKeys) ConfigOption {
	return func(config *Config) {
		config.skippedMessageKeys = storage
	}
}
