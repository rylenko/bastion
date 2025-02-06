package header

import (
	"encoding/binary"
	"fmt"
	"unsafe"

	"github.com/rylenko/bastion/pkg/ratchet/errors"
	"github.com/rylenko/bastion/pkg/ratchet/keys"
)

var uint64Size = int(unsafe.Sizeof(uint64(0)))

type Header struct {
	publicKey                         *keys.Public
	previousSendingChainMessagesCount uint64
	messageNumber                     uint64
}

func Decode(bytes []byte) (*Header, error) {
	if len(bytes) <= 2*uint64Size {
		return nil, fmt.Errorf("%w: not enough bytes", errors.ErrInvalidValue)
	}

	key := keys.NewPublic(make([]byte, len(bytes)-2*uint64Size))
	copy(key.Bytes(), bytes[2*uint64Size:])

	previousMessagesCount := binary.LittleEndian.Uint64(bytes[uint64Size : 2*uint64Size])
	messageNumber := binary.LittleEndian.Uint64(bytes[:uint64Size])

	return New(key, previousMessagesCount, messageNumber), nil
}

func New(publicKey *keys.Public, previousSendingChainMessagesCount, messageNumber uint64) *Header {
	return &Header{
		publicKey:                         publicKey,
		previousSendingChainMessagesCount: previousSendingChainMessagesCount,
		messageNumber:                     messageNumber,
	}
}

func (h *Header) Encode() ([]byte, error) {
	if h.publicKey == nil {
		return nil, fmt.Errorf("%w: public key is nil", errors.ErrInvalidValue)
	}

	buf := make([]byte, 2*uint64Size+len(h.publicKey.Bytes()))

	binary.LittleEndian.PutUint64(buf[:8], h.messageNumber)
	binary.LittleEndian.PutUint64(buf[8:16], h.previousSendingChainMessagesCount)
	copy(buf[8:], h.publicKey.Bytes())

	return buf, nil
}

func (h *Header) MessageNumber() uint64 {
	return h.messageNumber
}
