package keys

// Header is the header encryption key of participant's sending or receiving chain. The key must not be shared anywhere.
type Header struct {
	bytes []byte
}

// NewHeader creates a new instance of master key.
func NewHeader(bytes []byte) *Header {
	return &Header{bytes: bytes}
}
