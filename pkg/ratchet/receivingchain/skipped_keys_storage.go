package receivingchain

import (
	"fmt"

	"github.com/rylenko/bastion/pkg/ratchet/keys"
)

const skippedKeysStorageAddAtOnceLimit = 32

type (
	SkippedKeysIter  func(yield SkippedKeysYield)
	SkippedKeysYield func(headerKey keys.Header, messageNumberKeysIter SkippedMessageNumberKeysIter) bool

	SkippedMessageNumberKeysIter  func(yield SkippedMessageNumberKeysYield)
	SkippedMessageNumberKeysYield func(number uint64, key keys.Message) bool
)

type SkippedKeysStorage interface {
	// Add must add new skipped keys to storage.
	//
	// Total count indicates the total count of additions that will be made at once. This is useful to only allocate memory
	// once or to limit clutter from a large number of missed keys.
	Add(totalAtOnceCount uint64, headerKey keys.Header, messageNumber uint64, messageKey keys.Message) error

	// Clone must deep clone a storage.
	Clone() SkippedKeysStorage

	// Delete must delete skipped keys by header key and message number.
	Delete(headerKey keys.Header, messageNumber uint64) error

	// GetIter must return function, which iterates over all skipped keys.
	GetIter() (SkippedKeysIter, error)
}

type skippedKeysStorage map[string]map[uint64]keys.Message

func (st skippedKeysStorage) Add(
	totalAtOnceCount uint64,
	headerKey keys.Header,
	messageNumber uint64,
	messageKey keys.Message,
) error {
	if totalAtOnceCount > skippedKeysStorageAddAtOnceLimit {
		return fmt.Errorf("total count limit: %d > %d", totalAtOnceCount, skippedKeysStorageAddAtOnceLimit)
	}

	headerKeyString := string(headerKey.Bytes)
	if _, exists := st[headerKeyString]; !exists {
		st[headerKeyString] = make(map[uint64]keys.Message, totalAtOnceCount)
	}

	st[headerKeyString][messageNumber] = messageKey

	return nil
}

func (st skippedKeysStorage) Clone() SkippedKeysStorage {
	stClone := make(skippedKeysStorage, len(st))

	for headerKeyString, messageNumberKeys := range st {
		messageNumberKeysClone := make(map[uint64]keys.Message, len(messageNumberKeys))

		for messageNumber, messageKey := range messageNumberKeys {
			messageNumberKeysClone[messageNumber] = messageKey.Clone()
		}

		stClone[headerKeyString] = messageNumberKeysClone
	}

	return stClone
}

func (st skippedKeysStorage) Delete(headerKey keys.Header, messageNumber uint64) error {
	delete(st[string(headerKey.Bytes)], messageNumber)

	return nil
}

func (st skippedKeysStorage) GetIter() (SkippedKeysIter, error) {
	iter := func(yield SkippedKeysYield) {
		for headerKeyString, messageNumberKeys := range st {
			headerKey := keys.Header{Bytes: []byte(headerKeyString)}

			messageNumberKeysIter := func(yield SkippedMessageNumberKeysYield) {
				for messageNumber, messageKey := range messageNumberKeys {
					if !yield(messageNumber, messageKey) {
						return
					}
				}
			}

			if !yield(headerKey, messageNumberKeysIter) {
				return
			}
		}
	}

	return iter, nil
}
