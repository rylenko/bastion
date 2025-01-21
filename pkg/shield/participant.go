package shield

import (
	"fmt"

	"github.com/rylenko/sapphire/pkg/shield/keys"
)

// Participant is a participant in the coversation.
type Participant struct {
	localPrivateKey keys.Private
}

// NewRecipient creates a receiving participant in the conversation.
func NewRecipient(localPrivateKey keys.Private) *Participant {
	return &Participant{
		localPrivateKey: localPrivateKey,
	}
}

// NewSender creates a sending participant in the conversation.
func NewSender(provider Provider) (*Participant, error) {
	localPrivateKey, err := provider.GeneratePrivateKey()
	if err != nil {
		return nil, fmt.Errorf("GeneratePrivateKey(): %w", err)
	}

	participant := &Participant{
		localPrivateKey: localPrivateKey,
	}

	return participant, nil
}
