package keys

type MessageMaster struct {
	bytes []byte
}

func NewMessageMaster(bytes []byte) *MessageMaster {
	return &MessageMaster{bytes: bytes}
}

func (mk *MessageMaster) Bytes() []byte {
	return mk.bytes
}

func (mk *MessageMaster) Clone() *MessageMaster {
	if mk == nil {
		return nil
	}

	return NewMessageMaster(cloneBytes(mk.bytes))
}
