package receivingchain

import (
	stderrors "errors"
	"fmt"

	"github.com/rylenko/bastion/pkg/ratchet/errors"
	"github.com/rylenko/bastion/pkg/ratchet/header"
	"github.com/rylenko/bastion/pkg/ratchet/keys"
	"github.com/rylenko/bastion/pkg/ratchet/utils"
)

// Ratchet receiving chain.
//
// Please note that this structure may corrupt its state in case of errors. Therefore, clone the data at the top level
// and replace the current data with it if there are no errors.
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

func (ch Chain) Clone() Chain {
	ch.masterKey = ch.masterKey.ClonePtr()
	ch.headerKey = ch.headerKey.ClonePtr()
	ch.nextHeaderKey = ch.nextHeaderKey.Clone()
	ch.cfg = ch.cfg.clone()

	return ch
}

func (ch *Chain) Decrypt(
	encryptedHeader []byte,
	encryptedData []byte,
	auth []byte,
	ratchet RatchetCallback,
) (header.Header, []byte, error) {
	decryptedHeader, decryptedData, err := ch.decryptWithSkippedKeys(encryptedHeader, encryptedData, auth)
	if err == nil {
		return decryptedHeader, decryptedData, nil
	}

	err = fmt.Errorf("decrypt with skipped keys: %w", err)

	decryptedHeader, handleErr := ch.handleEncryptedHeader(encryptedHeader, ratchet)
	if handleErr != nil {
		err = stderrors.Join(err, fmt.Errorf("handle encrypted header: %w", handleErr))
		return header.Header{}, nil, err
	}

	messageKey, advanceErr := ch.advance()
	if advanceErr != nil {
		err = stderrors.Join(err, fmt.Errorf("advance chain: %w", err))
		return header.Header{}, nil, err
	}

	auth = utils.ConcatByteSlices(encryptedHeader, auth)

	decryptedData, decryptErr := ch.cfg.crypto.DecryptMessage(messageKey, encryptedData, auth)
	if decryptErr != nil {
		err = stderrors.Join(err, fmt.Errorf("%w: decrypt message: %w", errors.ErrCrypto, decryptErr))
		return header.Header{}, nil, err
	}

	// Note that here it is ok to ignore an error when decrypting with skipped keys if decryption with the next message key
	// succeeds.
	return decryptedHeader, decryptedData, nil
}

func (ch *Chain) Upgrade(masterKey keys.MessageMaster, nextHeaderKey keys.Header) {
	ch.masterKey = &masterKey
	ch.headerKey = &ch.nextHeaderKey
	ch.nextHeaderKey = nextHeaderKey
	ch.nextMessageNumber = 0
}

func (ch *Chain) advance() (keys.Message, error) {
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

// decryptHeaderWithCurrentOrNextKeys must decrypt passed encrypted header with current or next header key.
//
// Note that ratchet is needed if header decrypted with next header key.
func (ch *Chain) decryptHeaderWithCurrentOrNextKey(
	encryptedHeader []byte,
) (decryptedHeader header.Header, needRatchet bool, err error) {
	if ch.headerKey != nil {
		header, decryptErr := ch.cfg.crypto.DecryptHeader(*ch.headerKey, encryptedHeader, ch.nextMessageNumber)
		if decryptErr == nil {
			return header, false, nil
		}

		err = stderrors.Join(err, fmt.Errorf("%w: decrypt header with current key: %w", errors.ErrCrypto, decryptErr))
	}

	decryptedHeader, decryptErr := ch.cfg.crypto.DecryptHeader(ch.nextHeaderKey, encryptedHeader, ch.nextMessageNumber)
	if decryptErr != nil {
		err = stderrors.Join(err, fmt.Errorf("%w: decrypt header with next header key: %w", errors.ErrCrypto, decryptErr))
		return header.Header{}, false, err
	}

	// Note that here it is ok to ignore an error when decrypting with the current key if decryption with the next key
	// succeeds.
	return decryptedHeader, true, nil
}

func (ch *Chain) decryptWithSkippedKeys(encryptedHeader, encryptedData, auth []byte) (header.Header, []byte, error) {
	iter, err := ch.cfg.skippedKeysStorage.GetIter()
	if err != nil {
		return header.Header{}, nil, fmt.Errorf("%w: get iter: %w", errors.ErrSkippedKeysStorage, err)
	}

	for headerKey, messageNumberKeys := range iter {
		decryptedHeader, err := ch.cfg.crypto.DecryptHeader(headerKey, encryptedHeader, ch.nextMessageNumber)
		if err != nil {
			continue
		}

		for messageNumber, messageKey := range messageNumberKeys {
			if messageNumber != decryptedHeader.MessageNumber {
				continue
			}

			decryptedData, err := ch.cfg.crypto.DecryptMessage(messageKey, encryptedData, auth)
			if err != nil {
				return header.Header{}, nil, fmt.Errorf("%w: decrypt message: %w", errors.ErrCrypto, err)
			}

			if err := ch.cfg.skippedKeysStorage.Delete(headerKey, messageNumber); err != nil {
				return header.Header{}, nil, fmt.Errorf("%w: delete: %w", errors.ErrSkippedKeysStorage, err)
			}

			return decryptedHeader, decryptedData, nil
		}
	}

	return header.Header{}, nil, stderrors.New("no keys to decrypt header and data")
}

func (ch *Chain) handleEncryptedHeader(encryptedHeader []byte, ratchet RatchetCallback) (header.Header, error) {
	decryptedHeader, needRatchet, err := ch.decryptHeaderWithCurrentOrNextKey(encryptedHeader)
	if err != nil {
		return header.Header{}, fmt.Errorf("decrypt header: %w", err)
	}

	if needRatchet {
		if err := ch.skipKeys(decryptedHeader.PreviousSendingChainMessagesCount); err != nil {
			return header.Header{}, fmt.Errorf("skip %d keys: %w", decryptedHeader.PreviousSendingChainMessagesCount, err)
		}

		if err := ratchet(decryptedHeader.PublicKey); err != nil {
			return header.Header{}, fmt.Errorf("ratchet: %w", err)
		}
	}

	if err := ch.skipKeys(decryptedHeader.MessageNumber); err != nil {
		return header.Header{}, fmt.Errorf("skip %d message keys in upgraded chain: %w", decryptedHeader.MessageNumber, err)
	}

	return decryptedHeader, nil
}

func (ch *Chain) skipKeys(untilMessageNumber uint64) error {
	if untilMessageNumber < ch.nextMessageNumber {
		return fmt.Errorf("message number is small for the current chain, next message number is %d", ch.nextMessageNumber)
	}

	keysCount := untilMessageNumber - ch.nextMessageNumber

	for messageNumber := ch.nextMessageNumber; messageNumber < untilMessageNumber; messageNumber++ {
		messageKey, err := ch.advance()
		if err != nil {
			return fmt.Errorf("advance chain: %w", err)
		}

		if ch.headerKey == nil {
			return fmt.Errorf("%w: header key is nil", errors.ErrInvalidValue)
		}

		if err := ch.cfg.skippedKeysStorage.Add(keysCount, *ch.headerKey, messageNumber, messageKey); err != nil {
			return fmt.Errorf("%w: add: %w", errors.ErrSkippedKeysStorage, err)
		}
	}

	return nil
}

// RatchetCallback must perform ratchet and upgrade receiving chain.
type RatchetCallback func(remotePublicKey keys.Public) error
