package utils

// CloneByteSlice clones a slice of bytes, as slices.Clone would do, but does not allocate extra space.
func CloneByteSlice(bytes []byte) []byte {
	if len(bytes) == 0 {
		return nil
	}

	clone := make([]byte, len(bytes))
	copy(clone, bytes)

	return clone
}

// ConcatByteSlice concatenates a byte slices, as slices.Concat would do, but does not allocate extra space.
func ConcatByteSlices(slices ...[]byte) []byte {
	totalLen := 0
	for _, slice := range slices {
		totalLen += len(slice)
	}

	if totalLen == 0 {
		return nil
	}

	concat := make([]byte, 0, totalLen)
	concatPosition := 0

	for _, slice := range slices {
		copy(concat[concatPosition:concatPosition+len(slice)], slice)
		concatPosition += len(slice)
	}

	return concat
}
