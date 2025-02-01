package keys

// Message is the message encryption key. The key must not be shared anywhere.
type Message struct {
	bytes []byte
}

// NewMessageMaster creates a new instance of master key.
func NewMessage(bytes []byte) *Message {
	return &Message{bytes: bytes}
}
