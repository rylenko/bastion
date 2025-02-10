package keys

type Message struct {
	Bytes []byte
}

func NewMessage(bytes []byte) *Message {
	return &Message{Bytes: bytes}
}
