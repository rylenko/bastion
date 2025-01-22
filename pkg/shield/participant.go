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
}

// NewRecipient creates a receiving participant in the conversation.
func NewRecipient(localPrivateKey *keys.Private, rootKey *keys.Root) *Participant {
	return newParticipant(localPrivateKey, nil, chains.NewRoot(rootKey))
}

// NewSender creates a sending participant in the conversation.
func NewSender(provider Provider, remotePublicKey *keys.Public, rootKey *keys.Root) (*Participant, error) {
	if provider == nil {
		return nil, fmt.Errorf("%w: provider is nil", ErrInvalidValue)
	}

	localPrivateKey, err := provider.GeneratePrivateKey()
	if err != nil {
		return nil, fmt.Errorf("GeneratePrivateKey(): %w", err)
	}

	rootChain := chains.NewRoot(rootKey)
	// sendingKey, sendingNextHeaderKey := rootChain.Forward(
	// provider.ComputeSharedSecretKey(localPrivateKey, remotePublicKey))

	participant := newParticipant(localPrivateKey, remotePublicKey, rootChain)

	return participant, nil
}

func newParticipant(localPrivateKey *keys.Private, remotePublicKey *keys.Public, rootChain *chains.Root) *Participant {
	return &Participant{
		localPrivateKey: localPrivateKey,
		remotePublicKey: remotePublicKey,
		rootChain:       rootChain,
	}
}
