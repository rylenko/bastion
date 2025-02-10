package keys

type Root struct {
	Bytes []byte
}

func NewRoot(bytes []byte) *Root {
	return &Root{Bytes: bytes}
}

func (rk *Root) Clone() *Root {
	if rk == nil {
		return nil
	}

	return NewRoot(cloneBytes(rk.Bytes))
}
