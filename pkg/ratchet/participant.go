package ratchet

import (
	"fmt"

	"github.com/rylenko/bastion/pkg/ratchet/keys"
	"github.com/rylenko/bastion/pkg/ratchet/receivingchain"
	"github.com/rylenko/bastion/pkg/ratchet/rootchain"
	"github.com/rylenko/bastion/pkg/ratchet/sendingchain"
)

// Participant is the participant of the conversation.
//
// Please note that the structure is not safe for concurrent programs.
type Participant struct {
	localPrivateKey *keys.Private
	remotePublicKey *keys.Public
	rootChain       *rootchain.RootChain
	sendingChain    *sendingchain.SendingChain
	receivingChain  *receivingchain.ReceivingChain
	config          *config
}

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
		sendingchain.NewEmpty(nil, nil, sendingChainNextHeaderKey, config.sendingChainConfig),
		receivingchain.NewEmpty(receivingChainNextHeaderKey, config.receivingChainConfig),
		config,
	)
}

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
		return nil, fmt.Errorf("advance root chain: %w", err)
	}

	participant := newParticipant(
		localPrivateKey,
		remotePublicKey,
		rootChain,
		sendingchain.NewEmpty(sendingChainKey, sendingChainHeaderKey, sendingChainNextHeaderKey, config.sendingChainConfig),
		receivingchain.NewEmpty(receivingChainNextHeaderKey, config.receivingChainConfig),
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

func (p *Participant) Clone() *Participant {
	return newParticipant(
		p.localPrivateKey.Clone(),
		p.remotePublicKey.Clone(),
		p.rootChain.Clone(),
		p.sendingChain.Clone(),
		p.receivingChain.Clone(),
		p.config,
	)
}

func (p *Participant) Encrypt() error {
	tx := func(_ *Participant) error {
		return nil
	}

	if err := p.updateWithTx(tx); err != nil {
		return err
	}

	return nil
}

//nolint:unused // TODO: use
func (p *Participant) ratchet(remotePublicKey *keys.Public) error {
	p.remotePublicKey = remotePublicKey

	if p.config.crypto == nil {
		return fmt.Errorf("%w: config crypto is nil", ErrInvalidValue)
	}

	sharedSecretKey, err := p.config.crypto.ComputeSharedSecretKey(p.localPrivateKey, remotePublicKey)
	if err != nil {
		return fmt.Errorf("%w: compute shared secret key for receiving chain upgrade: %w", ErrCrypto, err)
	}

	newMasterKey, newNextHeaderKey, err := p.rootChain.Advance(sharedSecretKey)
	if err != nil {
		return fmt.Errorf("advance root chain for receiving chain upgrade: %w", err)
	}

	p.receivingChain.Upgrade(newMasterKey, newNextHeaderKey)

	p.localPrivateKey, err = p.config.crypto.GeneratePrivateKey()
	if err != nil {
		return fmt.Errorf("%w: generate new private key: %w", ErrCrypto, err)
	}

	sharedSecretKey, err = p.config.crypto.ComputeSharedSecretKey(p.localPrivateKey, remotePublicKey)
	if err != nil {
		return fmt.Errorf("%w: compute shared secret key for sending chain upgrade: %w", ErrCrypto, err)
	}

	newMasterKey, newNextHeaderKey, err = p.rootChain.Advance(sharedSecretKey)
	if err != nil {
		return fmt.Errorf("advance root chain for sending chain upgrade: %w", err)
	}

	p.sendingChain.Upgrade(newMasterKey, newNextHeaderKey)

	return nil
}

func (p *Participant) updateWithTx(tx func(p *Participant) error) error {
	newParticipant := p.Clone()
	if err := tx(newParticipant); err != nil {
		return err
	}

	*p = *newParticipant

	return nil
}
