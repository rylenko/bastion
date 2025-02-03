package keys

type Root struct {
	bytes []byte
}

func NewRoot(bytes []byte) *Root {
	return &Root{bytes: bytes}
}

func (rk *Root) Bytes() []byte {
	return rk.bytes
}

func (rk *Root) Clone() *Root {
	return NewRoot(cloneBytes(rk.bytes))
}
