package utils

import (
	"fmt"

	"github.com/rylenko/bastion/pkg/ratchet/errors"
)

type TxFunc[T any] func(dirty *T) error

// UpdateWithTx passes the dirty value to the transaction function and, if the transaction occurred without errors,
// replaces the target with the dirty version. This is useful for use in public methods to avoid corrupting the state
// with errors.
func UpdateWithTx[T any](target *T, dirty T, tx TxFunc[T]) error {
	if err := tx(&dirty); err != nil {
		return err
	}

	if target == nil {
		return fmt.Errorf("%w: target is nil", errors.ErrInvalidValue)
	}

	*target = dirty

	return nil
}
