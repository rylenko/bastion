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
	if rk == nil {
		return nil
	}

	return NewRoot(cloneBytes(rk.bytes))
}
