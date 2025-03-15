package utils

import "testing"

func TestIsNil(t *testing.T) {
	t.Parallel()

	intValue := 123

	tests := []struct {
		arg    any
		result bool
	}{
		{nil, true},
		{[]int(nil), true},
		{map[string]string(nil), true},
		{(*int)(nil), true},
		{5, false},
		{"string", false},
		{[]int{1, 2, 3}, false},
		{&intValue, false},
	}

	for _, test := range tests {
		if result := IsNil(test.arg); result != test.result {
			t.Fatalf("IsNil(%v) expected %t but got %t", test.arg, test.result, result)
		}
	}
}
