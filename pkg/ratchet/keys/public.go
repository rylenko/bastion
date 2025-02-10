package keys

type Public struct {
	Bytes []byte
}

func NewPublic(bytes []byte) *Public {
	return &Public{Bytes: bytes}
}

func (pk *Public) Clone() *Public {
	if pk == nil {
		return nil
	}

	return NewPublic(cloneBytes(pk.Bytes))
}
