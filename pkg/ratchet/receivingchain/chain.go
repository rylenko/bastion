package receivingchain

import (
	stderrors "errors"
	"fmt"

	"github.com/rylenko/bastion/pkg/ratchet/errors"
	"github.com/rylenko/bastion/pkg/ratchet/header"
	"github.com/rylenko/bastion/pkg/ratchet/keys"
)

type RatchetCallback func(remotePublicKey keys.Public) error

type Chain struct {
	masterKey         *keys.MessageMaster
	headerKey         *keys.Header
	nextHeaderKey     keys.Header
	nextMessageNumber uint64
	cfg               config
}

func New(
	masterKey *keys.MessageMaster,
	headerKey *keys.Header,
	nextHeaderKey keys.Header,
	nextMessageNumber uint64,
	options ...Option,
) (Chain, error) {
	cfg, err := newConfig(options)
	if err != nil {
		return Chain{}, fmt.Errorf("new config: %w", err)
	}

	chain := Chain{
		masterKey:         masterKey,
		headerKey:         headerKey,
		nextHeaderKey:     nextHeaderKey,
		nextMessageNumber: nextMessageNumber,
		cfg:               cfg,
	}

	return chain, nil
}

func (ch *Chain) Advance() (keys.Message, error) {
	if ch.masterKey == nil {
		return keys.Message{}, fmt.Errorf("%w: master key is nil", errors.ErrInvalidValue)
	}

	newMasterKey, messageKey, err := ch.cfg.crypto.AdvanceChain(*ch.masterKey)
	if err != nil {
		return keys.Message{}, fmt.Errorf("%w: advance via crypto: %w", errors.ErrCrypto, err)
	}

	ch.masterKey = &newMasterKey
	ch.nextMessageNumber++

	return messageKey, nil
}

func (ch Chain) Clone() Chain {
	ch.masterKey = ch.masterKey.ClonePtr()
	ch.headerKey = ch.headerKey.ClonePtr()
	ch.nextHeaderKey = ch.nextHeaderKey.Clone()
	ch.cfg = ch.cfg.clone()

	return ch
}

func (ch *Chain) HandleEncryptedHeader(encryptedHeader []byte, ratchet RatchetCallback) (header.Header, error) {
	decryptedHeader, needRatchet, err := ch.decryptHeader(encryptedHeader)
	if err != nil {
		return header.Header{}, fmt.Errorf("decrypt header: %w", err)
	}

	if needRatchet {
		if err := ch.skipMessageKeys(decryptedHeader.PreviousSendingChainMessagesCount); err != nil {
			return header.Header{}, fmt.Errorf("skip %d keys: %w", decryptedHeader.PreviousSendingChainMessagesCount, err)
		}

		if err := ratchet(decryptedHeader.PublicKey); err != nil {
			return header.Header{}, fmt.Errorf("ratchet: %w", err)
		}
	}

	if err := ch.skipMessageKeys(decryptedHeader.MessageNumber); err != nil {
		return header.Header{}, fmt.Errorf("skip %d message keys in upgraded chain: %w", decryptedHeader.MessageNumber, err)
	}

	return decryptedHeader, nil
}

func (ch *Chain) Upgrade(masterKey keys.MessageMaster, nextHeaderKey keys.Header) {
	ch.masterKey = &masterKey
	ch.headerKey = &ch.nextHeaderKey
	ch.nextHeaderKey = nextHeaderKey
	ch.nextMessageNumber = 0
}

func (ch *Chain) decryptHeader(encryptedHeader []byte) (decryptedHeader header.Header, needRatchet bool, err error) {
	if ch.headerKey != nil {
		header, decryptErr := ch.cfg.crypto.DecryptHeader(*ch.headerKey, encryptedHeader, ch.nextMessageNumber)
		if decryptErr == nil {
			return header, false, nil
		}

		err = stderrors.Join(err, decryptErr)
	}

	decryptedHeader, decryptErr := ch.cfg.crypto.DecryptHeader(ch.nextHeaderKey, encryptedHeader, ch.nextMessageNumber)
	if decryptErr != nil {
		err = stderrors.Join(err, decryptErr)
		return header.Header{}, false, fmt.Errorf("decrypt header with current and next header keys: %w", err)
	}

	// Note that here it is ok to ignore an error when decrypting with the current key if decryption with the next key
	// succeeds.
	return decryptedHeader, true, nil
}

func (ch *Chain) skipMessageKeys(untilMessageNumber uint64) error {
	// TODO
	return fmt.Errorf("%d", untilMessageNumber)
}
