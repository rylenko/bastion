package messagechaincommon

import "errors"

var (
	ErrInvalidValue = errors.New("invalid value")
	ErrMAC          = errors.New("MAC")
	ErrNewHash      = errors.New("new hash")
)
