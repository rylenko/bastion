package chains

import "github.com/rylenko/sapphire/pkg/shield/keys"

// Sending is the sending chain, which is responsible for encrypting the messages being sent. The sending chain of the
// sender is equal to the receiving chain of the recipient.
type Sending struct {
	masterKey                         *keys.MessageMaster
	headerKey                         *keys.Header
	nextHeaderKey                     *keys.Header
	nextMessageNumber                 uint32
	previousSendingChainMessagesCount uint32
}

// NewSending creates a new sending chain.
func NewSending(masterKey *keys.MessageMaster, headerKey, nextHeaderKey *keys.Header) *Sending {
	return &Sending{
		masterKey:                         masterKey,
		headerKey:                         headerKey,
		nextHeaderKey:                     nextHeaderKey,
		nextMessageNumber:                 0,
		previousSendingChainMessagesCount: 0,
	}
}
