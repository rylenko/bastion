package keys

func cloneBytes(bytes []byte) []byte {
	if len(bytes) == 0 {
		return nil
	}

	clone := make([]byte, len(bytes))
	copy(clone, bytes)

	return clone
}
