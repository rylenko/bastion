package keys

type MessageMaster struct {
	Bytes []byte
}

func NewMessageMaster(bytes []byte) *MessageMaster {
	return &MessageMaster{Bytes: bytes}
}

func (mk *MessageMaster) Clone() *MessageMaster {
	if mk == nil {
		return nil
	}

	return NewMessageMaster(cloneBytes(mk.Bytes))
}
