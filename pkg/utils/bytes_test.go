package utils

import (
	"slices"
	"testing"
)

func TestCloneByteSlice(t *testing.T) {
	t.Parallel()

	testSlices := [][]byte{
		nil,
		[]byte{},
		make([]byte, 0, 10),
		make([]byte, 4, 10),
		[]byte{1, 2, 3, 4},
	}

	for _, slice := range testSlices {
		clone := CloneByteSlice(slice)
		if !slices.Equal(clone, slice) {
			t.Fatalf("CloneByteSlice(%v) returned %v", slice, clone)
		}

		if cap(clone) != len(slice) {
			t.Fatalf("CloneByteSlice(%v) expected capacity %d but got %d", slice, len(slice), cap(clone))
		}

		if clone == nil {
			continue
		}

		if len(clone) == 0 {
			if slices.Equal(append(clone, 123), slice) {
				t.Fatalf("CloneByteSlice(%v) returned the same slice (append)", slice)
			}
		}

		if len(clone) > 0 {
			clone[0]++
			if slices.Equal(clone, slice) {
				t.Fatalf("CloneByteSlice(%v) returned the same slice (change)", slice)
			}
		}
	}
}

func TestConcatByteSlices(t *testing.T) {
	t.Parallel()

	tests := []struct {
		slices [][]byte
		result []byte
	}{
		{[][]byte{[]byte{1, 2, 3}, []byte{4, 5, 6}, []byte{7, 8, 9}}, []byte{1, 2, 3, 4, 5, 6, 7, 8, 9}},
		{nil, nil},
		{[][]byte{}, nil},
		{[][]byte{[]byte{}, []byte{}, []byte{}}, nil},
		{[][]byte{[]byte{}, []byte{1}, []byte{2, 3}, []byte{}}, []byte{1, 2, 3}},
		{[][]byte{nil, []byte{1, 2}, nil, []byte{3}}, []byte{1, 2, 3}},
	}

	for _, test := range tests {
		result := ConcatByteSlices(test.slices...)
		if !slices.Equal(result, test.result) {
			t.Fatalf("ConcatByteSlices(%v): expected %v but got %v", test.slices, test.result, result)
		}

		if cap(result) != len(test.result) {
			t.Fatalf("ConcatByteSlices(%v): expected capacity %d but got %d", test.slices, cap(test.result), cap(result))
		}
	}
}
