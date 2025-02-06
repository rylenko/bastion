package keys

func cloneBytes(bytes []byte) []byte {
	clone := make([]byte, len(bytes))
	copy(clone, bytes)

	return clone
}
