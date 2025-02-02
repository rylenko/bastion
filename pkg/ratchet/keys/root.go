package keys

type Root struct {
	bytes []byte
}

func NewRoot(bytes []byte) *Root {
	return &Root{bytes: bytes}
}

func (r *Root) Bytes() []byte {
	return r.bytes
}

func (r *Root) Clone() *Root {
	return NewRoot(cloneBytes(r.bytes))
}
