package keys

type Header struct {
	bytes []byte
}

func NewHeader(bytes []byte) *Header {
	return &Header{bytes: bytes}
}

func (hk *Header) Clone() *Header {
	return NewHeader(cloneBytes(hk.bytes))
}
