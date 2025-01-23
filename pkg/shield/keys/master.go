package keys

// Master is the main key of participant's sending or receiving chain. The key must not be shared anywhere.
type Master struct {
	bytes []byte
}

// NewMaster creates a new instance of master key.
func NewMaster(bytes []byte) *Master {
	return &Master{bytes: bytes}
}
