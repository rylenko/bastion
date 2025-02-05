package keys

type Header struct {
	bytes []byte
}

func NewHeader(bytes []byte) *Header {
	return &Header{bytes: bytes}
}

func (hk *Header) Clone() *Header {
	if hk == nil {
		return nil
	}

	return NewHeader(cloneBytes(hk.bytes))
}
