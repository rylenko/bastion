package keys

type Private struct {
	bytes []byte
}

func NewPrivate(bytes []byte) *Private {
	return &Private{bytes: bytes}
}

func (pk *Private) Bytes() []byte {
	return pk.bytes
}

func (pk *Private) Clone() *Private {
	return NewPrivate(cloneBytes(pk.bytes))
}
