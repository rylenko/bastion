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
	PublicKey                         *keys.Public
	PreviousSendingChainMessagesCount uint64
	MessageNumber                     uint64
}

func Decode(bytes []byte) (*Header, error) {
	if len(bytes) <= 2*uint64Size {
		return nil, fmt.Errorf("%w: not enough bytes", errors.ErrInvalidValue)
	}

	messageNumber := binary.LittleEndian.Uint64(bytes[:uint64Size])
	previousMessagesCount := binary.LittleEndian.Uint64(bytes[uint64Size : 2*uint64Size])

	key := keys.NewPublic(make([]byte, len(bytes)-2*uint64Size))
	copy(key.Bytes, bytes[2*uint64Size:])

	return New(key, previousMessagesCount, messageNumber), nil
}

func New(publicKey *keys.Public, previousSendingChainMessagesCount, messageNumber uint64) *Header {
	return &Header{
		PublicKey:                         publicKey,
		PreviousSendingChainMessagesCount: previousSendingChainMessagesCount,
		MessageNumber:                     messageNumber,
	}
}

func (h *Header) Encode() ([]byte, error) {
	if h.PublicKey == nil {
		return nil, fmt.Errorf("%w: public key is nil", errors.ErrInvalidValue)
	}

	buf := make([]byte, 2*uint64Size+len(h.PublicKey.Bytes))

	binary.LittleEndian.PutUint64(buf[:uint64Size], h.MessageNumber)
	binary.LittleEndian.PutUint64(buf[uint64Size:2*uint64Size], h.PreviousSendingChainMessagesCount)
	copy(buf[2*uint64Size:], h.PublicKey.Bytes)

	return buf, nil
}
