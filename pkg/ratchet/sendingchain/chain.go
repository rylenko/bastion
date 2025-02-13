package sendingchain

import (
	"fmt"

	"github.com/rylenko/bastion/pkg/ratchet/errors"
	"github.com/rylenko/bastion/pkg/ratchet/keys"
)

type Chain struct {
	masterKey                  *keys.MessageMaster
	headerKey                  *keys.Header
	nextHeaderKey              keys.Header
	nextMessageNumber          uint64
	previousChainMessagesCount uint64
	cfg                        config
}

func New(
	masterKey *keys.MessageMaster,
	headerKey *keys.Header,
	nextHeaderKey keys.Header,
	nextMessageNumber uint64,
	previousChainMessagesCount uint64,
	options ...Option,
) (Chain, error) {
	cfg, err := newConfig(options)
	if err != nil {
		return Chain{}, fmt.Errorf("new config: %w", err)
	}

	chain := Chain{
		masterKey:                  masterKey,
		headerKey:                  headerKey,
		nextHeaderKey:              nextHeaderKey,
		nextMessageNumber:          nextMessageNumber,
		previousChainMessagesCount: previousChainMessagesCount,
		cfg:                        cfg,
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

func (ch *Chain) Clone() Chain {
	return Chain{
		masterKey:                  ch.masterKey.ClonePtr(),
		headerKey:                  ch.headerKey.ClonePtr(),
		nextHeaderKey:              ch.nextHeaderKey.Clone(),
		nextMessageNumber:          ch.nextMessageNumber,
		previousChainMessagesCount: ch.previousChainMessagesCount,
		cfg:                        ch.cfg,
	}
}

func (ch *Chain) HeaderKey() *keys.Header {
	return ch.headerKey
}

func (ch *Chain) NextMessageNumber() uint64 {
	return ch.nextMessageNumber
}

func (ch *Chain) PreviousChainMessagesCount() uint64 {
	return ch.previousChainMessagesCount
}

func (ch *Chain) Upgrade(masterKey keys.MessageMaster, nextHeaderKey keys.Header) {
	ch.masterKey = &masterKey
	ch.headerKey = &ch.nextHeaderKey
	ch.nextHeaderKey = nextHeaderKey
	ch.previousChainMessagesCount = ch.nextMessageNumber
	ch.nextMessageNumber = 0
}
