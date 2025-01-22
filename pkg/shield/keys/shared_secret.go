package keys

// SharedSecret is the shared secret key computed using the private key of one participant and the public key of the
// another. The key must be stored locally and not shared anywhere.
type SharedSecret struct {
	bytes []byte
}

// NewPrivate creates a new instance of private key.
func NewSharedSecret(bytes []byte) *SharedSecret {
	return &SharedSecret{bytes: bytes}
}
