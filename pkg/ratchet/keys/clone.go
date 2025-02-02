package keys

func cloneBytes(bytes []byte) []byte {
	clone := make([]byte, 0, len(bytes))
	return append(clone, bytes...)
}
