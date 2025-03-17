package utils

import (
	"errors"
	"testing"
)

func TestUpdateWithTxTargetIsNil(t *testing.T) {
	t.Parallel()

	err := UpdateWithTx(nil, 5, func(_ *int) error { return nil })
	if err == nil || err.Error() != "target is nil" {
		t.Fatalf("UpdateWithTx(nil, 5) expected \"target is nil\" error but got %v", err)
	}
}

func TestUpdateWithTxSuccess(t *testing.T) {
	t.Parallel()

	value := 123
	err := UpdateWithTx(&value, value, func(x *int) error {
		*x *= 2
		return nil
	})

	if err != nil || value != 246 {
		t.Fatalf("UpdateWithTx returned %d, %v", value, err)
	}
}

func TestUpdateWithTxError(t *testing.T) {
	t.Parallel()

	value := 123
	err := UpdateWithTx(&value, value, func(x *int) error {
		*x *= 123
		return errors.New("runtime error")
	})

	if err == nil || err.Error() != "runtime error" || value != 123 {
		t.Fatalf("UpdateWithTx returned %d, %v", value, err)
	}
}
