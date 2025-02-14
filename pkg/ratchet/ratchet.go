package ratchet

import (
	"fmt"

	"github.com/rylenko/bastion/pkg/ratchet/errors"
	"github.com/rylenko/bastion/pkg/ratchet/header"
	"github.com/rylenko/bastion/pkg/ratchet/keys"
	"github.com/rylenko/bastion/pkg/ratchet/receivingchain"
	"github.com/rylenko/bastion/pkg/ratchet/rootchain"
	"github.com/rylenko/bastion/pkg/ratchet/sendingchain"
	"github.com/rylenko/bastion/pkg/ratchet/utils"
)

// Ratchet is the participant of the conversation.
//
// Please note that the structure is not safe for concurrent programs.
type Ratchet struct {
	localPrivateKey keys.Private
	localPublicKey  keys.Public
	remotePublicKey *keys.Public
	rootChain       rootchain.Chain
	sendingChain    sendingchain.Chain
	receivingChain  receivingchain.Chain
	cfg             config
}

// TODO: try to reduce arguments count.
func NewRecipient(
	localPrivateKey keys.Private,
	localPublicKey keys.Public,
	rootKey keys.Root,
	sendingChainNextHeaderKey keys.Header,
	receivingChainNextHeaderKey keys.Header,
	options ...Option,
) (Ratchet, error) {
	cfg, err := newConfig(options)
	if err != nil {
		return Ratchet{}, fmt.Errorf("new config: %w", err)
	}

	rootChain, err := rootchain.New(rootKey, cfg.rootOptions...)
	if err != nil {
		return Ratchet{}, fmt.Errorf("new root chain: %w", err)
	}

	sendingChain, err := sendingchain.New(nil, nil, sendingChainNextHeaderKey, 0, 0, cfg.sendingOptions...)
	if err != nil {
		return Ratchet{}, fmt.Errorf("new sending chain: %w", err)
	}

	receivingChain, err := receivingchain.New(nil, nil, receivingChainNextHeaderKey, 0, cfg.receivingOptions...)
	if err != nil {
		return Ratchet{}, fmt.Errorf("new receiving chain: %w", err)
	}

	ratchet := Ratchet{
		localPrivateKey: localPrivateKey,
		localPublicKey:  localPublicKey,
		rootChain:       rootChain,
		sendingChain:    sendingChain,
		receivingChain:  receivingChain,
		cfg:             cfg,
	}

	return ratchet, nil
}

// TODO: try to reduce arguments count.
func NewSender(
	remotePublicKey keys.Public,
	rootKey keys.Root,
	sendingChainHeaderKey keys.Header,
	receivingChainNextHeaderKey keys.Header,
	options ...Option,
) (Ratchet, error) {
	cfg, err := newConfig(options)
	if err != nil {
		return Ratchet{}, fmt.Errorf("new config: %w", err)
	}

	localPrivateKey, localPublicKey, err := cfg.crypto.GenerateKeyPair()
	if err != nil {
		return Ratchet{}, fmt.Errorf("%w: generate key pair: %w", errors.ErrCrypto, err)
	}

	sharedKey, err := cfg.crypto.ComputeSharedKey(localPrivateKey, remotePublicKey)
	if err != nil {
		return Ratchet{}, fmt.Errorf("%w: compute shared key: %w", errors.ErrCrypto, err)
	}

	rootChain, err := rootchain.New(rootKey, cfg.rootOptions...)
	if err != nil {
		return Ratchet{}, fmt.Errorf("new root chain: %w", err)
	}

	sendingChainKey, sendingChainNextHeaderKey, err := rootChain.Advance(sharedKey)
	if err != nil {
		return Ratchet{}, fmt.Errorf("advance root chain: %w", err)
	}

	sendingChain, err := sendingchain.New(
		&sendingChainKey, &sendingChainHeaderKey, sendingChainNextHeaderKey, 0, 0, cfg.sendingOptions...)
	if err != nil {
		return Ratchet{}, fmt.Errorf("new sending chain: %w", err)
	}

	receivingChain, err := receivingchain.New(nil, nil, receivingChainNextHeaderKey, 0, cfg.receivingOptions...)
	if err != nil {
		return Ratchet{}, fmt.Errorf("new receiving chain: %w", err)
	}

	ratchet := Ratchet{
		localPrivateKey: localPrivateKey,
		localPublicKey:  localPublicKey,
		remotePublicKey: &remotePublicKey,
		rootChain:       rootChain,
		sendingChain:    sendingChain,
		receivingChain:  receivingChain,
		cfg:             cfg,
	}

	return ratchet, nil
}

func (r Ratchet) Clone() Ratchet {
	r.localPrivateKey = r.localPrivateKey.Clone()
	r.localPublicKey = r.localPublicKey.Clone()
	r.remotePublicKey = r.remotePublicKey.ClonePtr()
	r.rootChain = r.rootChain.Clone()
	r.sendingChain = r.sendingChain.Clone()
	r.receivingChain = r.receivingChain.Clone()

	return r
}

func (r *Ratchet) Encrypt(data, auth []byte) ([]byte, []byte, error) {
	var encryptedHeader, encryptedData []byte

	tx := func(newR *Ratchet) error {
		header := header.Header{
			PublicKey:                         newR.localPublicKey,
			PreviousSendingChainMessagesCount: newR.sendingChain.PreviousChainMessagesCount(),
			MessageNumber:                     newR.sendingChain.NextMessageNumber(),
		}

		headerKey := newR.sendingChain.HeaderKey()
		if headerKey == nil {
			return fmt.Errorf("%w: sending chain header key is nil", errors.ErrInvalidValue)
		}

		encryptedHeader, err := newR.cfg.crypto.EncryptHeader(*headerKey, header)
		if err != nil {
			return fmt.Errorf("%w: encrypt header: %w", errors.ErrCrypto, err)
		}

		messageKey, err := newR.sendingChain.Advance()
		if err != nil {
			return fmt.Errorf("advance sending chain: %w", err)
		}

		encryptedData, err = newR.cfg.crypto.Encrypt(messageKey, data, utils.ConcatByteSlices(encryptedHeader, auth))
		if err != nil {
			return fmt.Errorf("%w: encrypt data: %w", errors.ErrCrypto, err)
		}

		return nil
	}

	if err := r.updateWithTx(tx); err != nil {
		return nil, nil, err
	}

	return encryptedHeader, encryptedData, nil
}

//nolint:unused // TODO: use
func (r *Ratchet) ratchet(remotePublicKey keys.Public) error {
	r.remotePublicKey = &remotePublicKey

	sharedKey, err := r.cfg.crypto.ComputeSharedKey(r.localPrivateKey, remotePublicKey)
	if err != nil {
		return fmt.Errorf("%w: compute shared secret key for receiving chain upgrade: %w", errors.ErrCrypto, err)
	}

	newMasterKey, newNextHeaderKey, err := r.rootChain.Advance(sharedKey)
	if err != nil {
		return fmt.Errorf("advance root chain for receiving chain upgrade: %w", err)
	}

	r.receivingChain.Upgrade(newMasterKey, newNextHeaderKey)

	r.localPrivateKey, r.localPublicKey, err = r.cfg.crypto.GenerateKeyPair()
	if err != nil {
		return fmt.Errorf("%w: generate new key pair: %w", errors.ErrCrypto, err)
	}

	sharedKey, err = r.cfg.crypto.ComputeSharedKey(r.localPrivateKey, remotePublicKey)
	if err != nil {
		return fmt.Errorf("%w: compute shared secret key for sending chain upgrade: %w", errors.ErrCrypto, err)
	}

	newMasterKey, newNextHeaderKey, err = r.rootChain.Advance(sharedKey)
	if err != nil {
		return fmt.Errorf("advance root chain for sending chain upgrade: %w", err)
	}

	r.sendingChain.Upgrade(newMasterKey, newNextHeaderKey)

	return nil
}

// TODO: docs.
func (r *Ratchet) updateWithTx(txFunc func(r *Ratchet) error) error {
	newR := r.Clone()
	if err := txFunc(&newR); err != nil {
		return err
	}

	*r = newR

	return nil
}
