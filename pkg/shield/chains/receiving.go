package chains

import (
	"fmt"

	"github.com/rylenko/sapphire/pkg/shield/keys"
)

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

// Forward moves the receiving chain forward. In other words, creating a new message master key and a new message key.
//
// This method is a wrapper around the provider that sets a new message master key into the current chain.
func (chain *Receiving) Forward(provider Provider) (*keys.Message, error) {
	if provider == nil {
		return nil, fmt.Errorf("%w: provider is nil", ErrInvalidValue)
	}

	newMasterKey, messageKey, err := provider.ForwardMessageChain(chain.masterKey)
	if err != nil {
		return nil, fmt.Errorf("%w: forward message chain: %w", ErrProvider, err)
	}

	chain.masterKey = newMasterKey
	chain.nextMessageNumber++

	return messageKey, nil
}
