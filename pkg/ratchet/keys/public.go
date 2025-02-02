package keys

type Public struct {
	bytes []byte
}

func NewPublic(bytes []byte) *Public {
	return &Public{bytes: bytes}
}

func (p *Public) Bytes() []byte {
	return p.bytes
}

func (p *Public) Clone() *Public {
	return NewPublic(cloneBytes(p.bytes))
}
