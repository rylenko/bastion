package keys

type SharedSecret struct {
	bytes []byte
}

func NewSharedSecret(bytes []byte) *SharedSecret {
	return &SharedSecret{bytes: bytes}
}

func (s *SharedSecret) Bytes() []byte {
	return s.bytes
}
