package shield

import (
	"fmt"

	"github.com/rylenko/sapphire/pkg/shield/keys"
	"github.com/rylenko/sapphire/pkg/shield/receivingchain"
	"github.com/rylenko/sapphire/pkg/shield/rootchain"
	"github.com/rylenko/sapphire/pkg/shield/sendingchain"
)

// Participant is a participant in the coversation.
type Participant struct {
	localPrivateKey *keys.Private
	remotePublicKey *keys.Public
	rootChain       *rootchain.RootChain
	sendingChain    *sendingchain.SendingChain
	receivingChain  *receivingchain.ReceivingChain
	config          *config
}

// NewRecipient creates a receiving participant in the conversation.
//
// TODO: try to reduce arguments count.
func NewRecipient(
	localPrivateKey *keys.Private,
	rootKey *keys.Root,
	sendingChainNextHeaderKey *keys.Header,
	receivingChainNextHeaderKey *keys.Header,
	configOptions ...ConfigOption,
) *Participant {
	config := newConfig(configOptions...)

	return newParticipant(
		localPrivateKey,
		nil,
		rootchain.New(rootKey, config.rootChainConfig),
		sendingchain.New(nil, nil, sendingChainNextHeaderKey, config.sendingChainConfig),
		receivingchain.New(receivingChainNextHeaderKey, config.receivingChainConfig),
		config,
	)
}

// NewSender creates a sending participant in the conversation.
//
// TODO: try to reduce arguments count.
func NewSender(
	remotePublicKey *keys.Public,
	rootKey *keys.Root,
	sendingChainHeaderKey *keys.Header,
	receivingChainNextHeaderKey *keys.Header,
	configOptions ...ConfigOption,
) (*Participant, error) {
	config := newConfig(configOptions...)

	if config.crypto == nil {
		return nil, fmt.Errorf("%w: config crypto is nil", ErrInvalidValue)
	}

	localPrivateKey, err := config.crypto.GeneratePrivateKey()
	if err != nil {
		return nil, fmt.Errorf("%w: generate private key: %w", ErrCrypto, err)
	}

	sharedSecretKey, err := config.crypto.ComputeSharedSecretKey(localPrivateKey, remotePublicKey)
	if err != nil {
		return nil, fmt.Errorf("%w: compute shared secret key: %w", ErrCrypto, err)
	}

	rootChain := rootchain.New(rootKey, config.rootChainConfig)

	sendingChainKey, sendingChainNextHeaderKey, err := rootChain.Advance(sharedSecretKey)
	if err != nil {
		return nil, fmt.Errorf("%w: root: %w", ErrAdvanceChain, err)
	}

	participant := newParticipant(
		localPrivateKey,
		remotePublicKey,
		rootChain,
		sendingchain.New(sendingChainKey, sendingChainHeaderKey, sendingChainNextHeaderKey, config.sendingChainConfig),
		receivingchain.New(receivingChainNextHeaderKey, config.receivingChainConfig),
		config,
	)

	return participant, nil
}

func newParticipant(
	localPrivateKey *keys.Private,
	remotePublicKey *keys.Public,
	rootChain *rootchain.RootChain,
	sendingChain *sendingchain.SendingChain,
	receivingChain *receivingchain.ReceivingChain,
	config *config,
) *Participant {
	return &Participant{
		localPrivateKey: localPrivateKey,
		remotePublicKey: remotePublicKey,
		rootChain:       rootChain,
		sendingChain:    sendingChain,
		receivingChain:  receivingChain,
		config:          config,
	}
}
