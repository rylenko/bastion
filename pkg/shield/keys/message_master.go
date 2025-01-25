package keys

// MessageMaster is the master key of participant's message chains: sending and receiving chains. The master key should
// be used to generate message keys. The key must not be shared anywhere.
type MessageMaster struct {
	bytes []byte
}

// NewMessageMaster creates a new instance of master key.
func NewMessageMaster(bytes []byte) *MessageMaster {
	return &MessageMaster{bytes: bytes}
}

// Bytes returns message master key bytes.
func (key *MessageMaster) Bytes() []byte {
	return key.bytes
}
