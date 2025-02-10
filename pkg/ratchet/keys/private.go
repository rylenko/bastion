package keys

type Private struct {
	Bytes []byte
}

func NewPrivate(bytes []byte) *Private {
	return &Private{Bytes: bytes}
}

func (pk *Private) Clone() *Private {
	if pk == nil {
		return nil
	}

	return NewPrivate(cloneBytes(pk.Bytes))
}
