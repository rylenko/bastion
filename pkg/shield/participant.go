package shield

import (
	"fmt"

	"github.com/rylenko/sapphire/pkg/shield/chains"
	"github.com/rylenko/sapphire/pkg/shield/keys"
)

// Participant is a participant in the coversation.
type Participant struct {
	localPrivateKey *keys.Private
	remotePublicKey *keys.Public
	rootChain       *chains.Root
	sendingChain    *chains.Sending
	receivingChain  *chains.Receiving
}

// NewRecipient creates a receiving participant in the conversation.
func NewRecipient(localPrivateKey *keys.Private, rootKey *keys.Root) *Participant {
	return newParticipant(localPrivateKey, nil, chains.NewRoot(rootKey), nil, nil)
}

// NewSender creates a sending participant in the conversation.
func NewSender(
	provider Provider,
	remotePublicKey *keys.Public,
	rootKey *keys.Root,
	sendingChainHeaderKey *keys.Header,
	receivingChainNextHeaderKey *keys.Header,
) (*Participant, error) {
	if provider == nil {
		return nil, fmt.Errorf("%w: provider is nil", ErrInvalidValue)
	}

	localPrivateKey, err := provider.GeneratePrivateKey()
	if err != nil {
		return nil, fmt.Errorf("%w: generate private key: %w", ErrProvider, err)
	}

	sharedSecretKey, err := provider.ComputeSharedSecretKey(localPrivateKey, remotePublicKey)
	if err != nil {
		return nil, fmt.Errorf("%w: compute shared secret key: %w", ErrProvider, err)
	}

	rootChain := chains.NewRoot(rootKey)

	sendingChainKey, sendingChainNextHeaderKey, err := rootChain.Forward(provider, sharedSecretKey)
	if err != nil {
		return nil, fmt.Errorf("%w: root: %w", ErrForwardChain, err)
	}

	participant := newParticipant(
		localPrivateKey,
		remotePublicKey,
		rootChain,
		chains.NewSending(sendingChainKey, sendingChainHeaderKey, sendingChainNextHeaderKey),
		chains.NewReceiving(receivingChainNextHeaderKey),
	)

	return participant, nil
}

func newParticipant(
	localPrivateKey *keys.Private,
	remotePublicKey *keys.Public,
	rootChain *chains.Root,
	sendingChain *chains.Sending,
	receivingChain *chains.Receiving,
) *Participant {
	return &Participant{
		localPrivateKey: localPrivateKey,
		remotePublicKey: remotePublicKey,
		rootChain:       rootChain,
		sendingChain:    sendingChain,
		receivingChain:  receivingChain,
	}
}
