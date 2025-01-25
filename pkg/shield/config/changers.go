package config

// Changer is a callback that sets some parameters for the config.
type Changer func(config *Config)

// WithMessageKeysSkipLimit returns the config changer, which sets the config to the passed limit value.
//
// The changer can be passed to the constructor.
func WithMessageKeysSkipLimit(limit uint32) Changer {
	return func(config *Config) {
		config.MessageKeysSkipLimit = limit
	}
}
