package chains

import "github.com/rylenko/sapphire/pkg/shield/keys"

// Receiving is the receiving chain, which is responsible for decrypting the messages being received. The receiving
// chain of the recipient is equal to the sending chain of the sender.
type Receiving struct {
	masterKey          *keys.MessageMaster
	headerKey          *keys.Header
	nextHeaderKey      *keys.Header
	nextMessageNumber  uint32
	skippedMessageKeys any // TODO
}

// NewReceiving creates a new receiving chain.
func NewReceiving(nextHeaderKey *keys.Header) *Receiving {
	return &Receiving{
		masterKey:          nil,
		headerKey:          nil,
		nextHeaderKey:      nextHeaderKey,
		nextMessageNumber:  0,
		skippedMessageKeys: nil,
	}
}
