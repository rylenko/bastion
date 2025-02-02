package keys

type MessageMaster struct {
	bytes []byte
}

func NewMessageMaster(bytes []byte) *MessageMaster {
	return &MessageMaster{bytes: bytes}
}

func (m *MessageMaster) Bytes() []byte {
	return m.bytes
}

func (m *MessageMaster) Clone() *MessageMaster {
	return NewMessageMaster(cloneBytes(m.bytes))
}
