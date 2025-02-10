package keys

type SharedSecret struct {
	Bytes []byte
}

func NewSharedSecret(bytes []byte) *SharedSecret {
	return &SharedSecret{Bytes: bytes}
}
