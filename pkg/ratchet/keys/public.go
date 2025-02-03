package keys

type Public struct {
	bytes []byte
}

func NewPublic(bytes []byte) *Public {
	return &Public{bytes: bytes}
}

func (pk *Public) Bytes() []byte {
	return pk.bytes
}

func (pk *Public) Clone() *Public {
	return NewPublic(cloneBytes(pk.bytes))
}
