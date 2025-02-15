package utils

import (
	"fmt"

	"github.com/rylenko/bastion/pkg/ratchet/errors"
)

type TxFunc[T any] func(dirty *T) error

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
