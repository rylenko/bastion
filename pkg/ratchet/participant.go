package ratchet

import (
	"fmt"

	"github.com/rylenko/bastion/pkg/ratchet/errors"
	"github.com/rylenko/bastion/pkg/ratchet/header"
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
	localPublicKey  *keys.Public
	remotePublicKey *keys.Public
	rootChain       *rootchain.RootChain
	sendingChain    *sendingchain.SendingChain
	receivingChain  *receivingchain.ReceivingChain
	config          *config
}

// TODO: try to reduce arguments count.
func NewRecipient(
	localPrivateKey *keys.Private,
	localPublicKey *keys.Public,
	rootKey *keys.Root,
	sendingChainNextHeaderKey *keys.Header,
	receivingChainNextHeaderKey *keys.Header,
	configOptions ...ConfigOption,
) *Participant {
	config := newConfig(configOptions...)

	return newParticipant(
		localPrivateKey,
		localPublicKey,
		nil,
		rootchain.New(rootKey, config.rootChainConfig),
		sendingchain.New(nil, nil, sendingChainNextHeaderKey, 0, 0, config.sendingChainConfig),
		receivingchain.New(nil, nil, receivingChainNextHeaderKey, 0, config.receivingChainConfig),
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
		return nil, fmt.Errorf("%w: config crypto is nil", errors.ErrInvalidValue)
	}

	localPrivateKey, localPublicKey, err := config.crypto.GenerateKeyPair()
	if err != nil {
		return nil, fmt.Errorf("%w: generate key pair: %w", errors.ErrCrypto, err)
	}

	sharedSecretKey, err := config.crypto.ComputeSharedSecretKey(localPrivateKey, remotePublicKey)
	if err != nil {
		return nil, fmt.Errorf("%w: compute shared secret key: %w", errors.ErrCrypto, err)
	}

	rootChain := rootchain.New(rootKey, config.rootChainConfig)

	sendingChainKey, sendingChainNextHeaderKey, err := rootChain.Advance(sharedSecretKey)
	if err != nil {
		return nil, fmt.Errorf("advance root chain: %w", err)
	}

	participant := newParticipant(
		localPrivateKey,
		localPublicKey,
		remotePublicKey,
		rootChain,
		sendingchain.New(sendingChainKey, sendingChainHeaderKey, sendingChainNextHeaderKey, 0, 0, config.sendingChainConfig),
		receivingchain.New(nil, nil, receivingChainNextHeaderKey, 0, config.receivingChainConfig),
		config,
	)

	return participant, nil
}

func newParticipant(
	localPrivateKey *keys.Private,
	localPublicKey *keys.Public,
	remotePublicKey *keys.Public,
	rootChain *rootchain.RootChain,
	sendingChain *sendingchain.SendingChain,
	receivingChain *receivingchain.ReceivingChain,
	config *config,
) *Participant {
	return &Participant{
		localPrivateKey: localPrivateKey,
		localPublicKey:  localPublicKey,
		remotePublicKey: remotePublicKey,
		rootChain:       rootChain,
		sendingChain:    sendingChain,
		receivingChain:  receivingChain,
		config:          config,
	}
}

func (p *Participant) Clone() *Participant {
	if p == nil {
		return nil
	}

	return newParticipant(
		p.localPrivateKey.Clone(),
		p.localPublicKey.Clone(),
		p.remotePublicKey.Clone(),
		p.rootChain.Clone(),
		p.sendingChain.Clone(),
		p.receivingChain.Clone(),
		p.config,
	)
}

func (p *Participant) Encrypt(data, auth []byte) ([]byte, []byte, error) {
	var encryptedHeader, encryptedData []byte

	tx := func(newP *Participant) error {
		messageKey, err := newP.sendingChain.Advance()
		if err != nil {
			return fmt.Errorf("advance sending chain: %w", err)
		}

		header := header.New(
			newP.localPublicKey,
			newP.sendingChain.PreviousSendingChainMessagesCount(),
			newP.sendingChain.NextMessageNumber(),
		)

		if p.config.crypto == nil {
			return fmt.Errorf("%w: config crypto is nil", errors.ErrInvalidValue)
		}

		encryptedHeader, err = newP.config.crypto.EncryptHeader(newP.sendingChain.HeaderKey(), header)
		if err != nil {
			return fmt.Errorf("%w: encrypt header: %w", errors.ErrCrypto, err)
		}

		cryptoAuth := make([]byte, len(encryptedHeader)+len(auth))
		copy(cryptoAuth[:len(encryptedHeader)], encryptedHeader)
		copy(cryptoAuth[len(encryptedHeader):], auth)

		encryptedData, err = newP.config.crypto.Encrypt(messageKey, data, cryptoAuth)
		if err != nil {
			return fmt.Errorf("%w: encrypt data: %w", errors.ErrCrypto, err)
		}

		return nil
	}

	if err := p.updateWithTx(tx); err != nil {
		return nil, nil, err
	}

	return encryptedHeader, encryptedData, nil
}

//nolint:unused // TODO: use
func (p *Participant) ratchet(remotePublicKey *keys.Public) error {
	p.remotePublicKey = remotePublicKey

	if p.config.crypto == nil {
		return fmt.Errorf("%w: config crypto is nil", errors.ErrInvalidValue)
	}

	sharedSecretKey, err := p.config.crypto.ComputeSharedSecretKey(p.localPrivateKey, remotePublicKey)
	if err != nil {
		return fmt.Errorf("%w: compute shared secret key for receiving chain upgrade: %w", errors.ErrCrypto, err)
	}

	newMasterKey, newNextHeaderKey, err := p.rootChain.Advance(sharedSecretKey)
	if err != nil {
		return fmt.Errorf("advance root chain for receiving chain upgrade: %w", err)
	}

	p.receivingChain.Upgrade(newMasterKey, newNextHeaderKey)

	p.localPrivateKey, p.localPublicKey, err = p.config.crypto.GenerateKeyPair()
	if err != nil {
		return fmt.Errorf("%w: generate new key pair: %w", errors.ErrCrypto, err)
	}

	sharedSecretKey, err = p.config.crypto.ComputeSharedSecretKey(p.localPrivateKey, remotePublicKey)
	if err != nil {
		return fmt.Errorf("%w: compute shared secret key for sending chain upgrade: %w", errors.ErrCrypto, err)
	}

	newMasterKey, newNextHeaderKey, err = p.rootChain.Advance(sharedSecretKey)
	if err != nil {
		return fmt.Errorf("advance root chain for sending chain upgrade: %w", err)
	}

	p.sendingChain.Upgrade(newMasterKey, newNextHeaderKey)

	return nil
}

// TODO: docs.
func (p *Participant) updateWithTx(tx func(p *Participant) error) error {
	newP := p.Clone()
	if err := tx(newP); err != nil {
		return err
	}

	*p = *newP

	return nil
}
