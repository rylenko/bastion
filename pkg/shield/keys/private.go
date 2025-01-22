package keys

// Private is the private key of the conversation participant. The private key must be stored locally and not shared
// anywhere.
type Private struct {
	bytes []byte
}

// NewPrivate creates a new instance of private key.
func NewPrivate(bytes []byte) *Private {
	return &Private{bytes: bytes}
}
