package keys

import (
	"testing"
)

func TestHeaderClone(t *testing.T) {
	t.Parallel()

	headers := []Header{
		Header{Bytes: nil},
		Header{Bytes: []byte{1, 2, 3, 4, 5}},
	}

	for _, header := range headers {
		clone := header.Clone()

		testBytesClone(t, header.Bytes, clone.Bytes)
	}
}

func TestHeaderClonePtr(t *testing.T) {
	t.Parallel()

	headers := []*Header{
		nil,
		&Header{Bytes: nil},
		&Header{Bytes: []byte{1, 2, 3, 4, 5}},
	}

	for _, header := range headers {
		clone := header.ClonePtr()

		testClonePtrPointers(t, header, clone)

		if header != nil && clone != nil {
			testBytesClone(t, header.Bytes, clone.Bytes)
		}
	}
}
