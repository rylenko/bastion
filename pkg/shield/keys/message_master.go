package keys

type MessageMaster struct {
	bytes []byte
}

func NewMessageMaster(bytes []byte) *MessageMaster {
	return &MessageMaster{bytes: bytes}
}

func (key *MessageMaster) Bytes() []byte {
	return key.bytes
}
