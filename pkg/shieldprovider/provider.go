package shieldprovider

// Provider is a crypto provider for shield module.
type Provider struct{}

// New creates a new instance of provider.
func New() Provider {
	return Provider{}
}
