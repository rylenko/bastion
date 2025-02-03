package keys

type SharedSecret struct {
	bytes []byte
}

func NewSharedSecret(bytes []byte) *SharedSecret {
	return &SharedSecret{bytes: bytes}
}

func (sk *SharedSecret) Bytes() []byte {
	return sk.bytes
}
