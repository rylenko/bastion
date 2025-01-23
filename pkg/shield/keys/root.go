package keys

// Root is the key of conversation participant's root chain. The key must not be shared anywhere.
type Root struct {
	bytes []byte
}

// NewRoot creates a new instance of root key.
func NewRoot(bytes []byte) *Root {
	return &Root{bytes: bytes}
}

// Bytes returns root key bytes.
func (key *Root) Bytes() []byte {
	return key.bytes
}
