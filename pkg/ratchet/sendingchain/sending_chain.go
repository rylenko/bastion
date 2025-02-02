package sendingchain

import (
	"fmt"

	"github.com/rylenko/bastion/pkg/ratchet/keys"
)

type SendingChain struct {
	masterKey                         *keys.MessageMaster
	headerKey                         *keys.Header
	nextHeaderKey                     *keys.Header
	nextMessageNumber                 uint32
	previousSendingChainMessagesCount uint32
	config                            *Config
}

func New(masterKey *keys.MessageMaster, headerKey, nextHeaderKey *keys.Header, config *Config) *SendingChain {
	return &SendingChain{
		masterKey:                         masterKey,
		headerKey:                         headerKey,
		nextHeaderKey:                     nextHeaderKey,
		nextMessageNumber:                 0,
		previousSendingChainMessagesCount: 0,
		config:                            config,
	}
}

func (chain *SendingChain) Advance() (*keys.Message, error) {
	if chain.config == nil {
		return nil, fmt.Errorf("%w: config is nil", ErrInvalidValue)
	}

	if chain.config.crypto == nil {
		return nil, fmt.Errorf("%w: config crypto is nil", ErrInvalidValue)
	}

	newMasterKey, messageKey, err := chain.config.crypto.AdvanceChain(chain.masterKey)
	if err != nil {
		return nil, fmt.Errorf("%w: advance via crypto: %w", ErrCrypto, err)
	}

	chain.masterKey = newMasterKey
	chain.nextMessageNumber++

	return messageKey, nil
}

func (chain *SendingChain) Upgrade(masterKey *keys.MessageMaster, nextHeaderKey *keys.Header) {
	chain.masterKey = masterKey
	chain.headerKey = chain.nextHeaderKey
	chain.nextHeaderKey = nextHeaderKey
	chain.previousSendingChainMessagesCount = chain.nextMessageNumber
	chain.nextMessageNumber = 0
}
