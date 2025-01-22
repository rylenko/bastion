package keys

// Root is the key of conversation participant's root chain. The key must be stored locally and not shared
// anywhere.
type Root struct {
	bytes []byte
}

// NewRoot creates a new instance of root key.
func NewRoot(bytes []byte) *Root {
	return &Root{bytes: bytes}
}
