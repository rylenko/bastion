package config

// TODO: find optimal value.
const defaultMessageKeysSkipLimit = 1024

// Config is a collection of parameters for the shield algorithm.
//
// MessageKeysSkipLimit represents the allowed limit for skipping messages that have not yet delivered.
type Config struct {
	MessageKeysSkipLimit uint32
}

// NewConfig creates a new default config and applies each config changer passed in. This way you can change only the
// parameters you need.
func New(changers ...Changer) *Config {
	config := &Config{
		MessageKeysSkipLimit: defaultMessageKeysSkipLimit,
	}

	for _, changer := range changers {
		changer(config)
	}

	return config
}
