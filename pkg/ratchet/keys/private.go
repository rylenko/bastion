package keys

type Private struct {
	bytes []byte
}

func NewPrivate(bytes []byte) *Private {
	return &Private{bytes: bytes}
}

func (p *Private) Bytes() []byte {
	return p.bytes
}

func (p *Private) Clone() *Private {
	return NewPrivate(cloneBytes(p.bytes))
}
