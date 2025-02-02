package ratchet

import (
	"fmt"

	"github.com/rylenko/bastion/pkg/ratchet/keys"
	"github.com/rylenko/bastion/pkg/ratchet/receivingchain"
	"github.com/rylenko/bastion/pkg/ratchet/rootchain"
	"github.com/rylenko/bastion/pkg/ratchet/sendingchain"
)

// Participant is a participant in the coversation.
//
// TODO: Make the state error-resistant, i.e. save any changes only after successful operations. For example, make a
// wrapper around the state. Then outside the wrapper make a method update(func(state *state) error), which would clone
// the current state and send it to the passed callback. After the callback successfully completes, it would set the new
// state. Maybe we should to add a lock for parallel updating.
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

//nolint:unused // TODO: use
func (participant *Participant) ratchet(remotePublicKey *keys.Public) error {
	participant.remotePublicKey = remotePublicKey

	if participant.config.crypto == nil {
		return fmt.Errorf("%w: config crypto is nil", ErrInvalidValue)
	}

	sharedSecretKey, err := participant.config.crypto.ComputeSharedSecretKey(participant.localPrivateKey, remotePublicKey)
	if err != nil {
		return fmt.Errorf("%w: compute shared secret key for receiving chain upgrade: %w", ErrCrypto, err)
	}

	newMasterKey, newNextHeaderKey, err := participant.rootChain.Advance(sharedSecretKey)
	if err != nil {
		return fmt.Errorf("%w: root chain for receiving chain upgrade: %w", ErrAdvanceChain, err)
	}

	participant.receivingChain.Upgrade(newMasterKey, newNextHeaderKey)

	participant.localPrivateKey, err = participant.config.crypto.GeneratePrivateKey()
	if err != nil {
		return fmt.Errorf("%w: generate new private key: %w", ErrCrypto, err)
	}

	sharedSecretKey, err = participant.config.crypto.ComputeSharedSecretKey(participant.localPrivateKey, remotePublicKey)
	if err != nil {
		return fmt.Errorf("%w: compute shared secret key for sending chain upgrade: %w", ErrCrypto, err)
	}

	newMasterKey, newNextHeaderKey, err = participant.rootChain.Advance(sharedSecretKey)
	if err != nil {
		return fmt.Errorf("%w: root chain for sending chain upgrade: %w", ErrAdvanceChain, err)
	}

	participant.sendingChain.Upgrade(newMasterKey, newNextHeaderKey)

	return nil
}
