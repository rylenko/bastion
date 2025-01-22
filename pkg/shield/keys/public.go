package keys

// Public is the public key of the conversation participant. The public key may be shared.
type Public struct {
	bytes []byte
}

// NewPublic creates a new instance of public key.
func NewPublic(bytes []byte) *Public {
	return &Public{bytes: bytes}
}
