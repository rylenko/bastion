package keys

type Header struct {
	Bytes []byte
}

func NewHeader(bytes []byte) *Header {
	return &Header{Bytes: bytes}
}

func (hk *Header) Clone() *Header {
	if hk == nil {
		return nil
	}

	return NewHeader(cloneBytes(hk.Bytes))
}
