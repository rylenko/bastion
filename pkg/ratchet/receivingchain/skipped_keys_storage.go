package receivingchain

import (
	"fmt"

	"github.com/rylenko/bastion/pkg/ratchet/keys"
)

const (
	skippedKeysStorageMessageKeysLenLimit  = 1024
	skippedKeysStorageHeaderKeysLenToClear = 4
)

type (
	SkippedKeysIter  func(yield SkippedKeysYield)
	SkippedKeysYield func(headerKey keys.Header, messageNumberKeysIter SkippedMessageNumberKeysIter) bool

	SkippedMessageNumberKeysIter  func(yield SkippedMessageNumberKeysYield)
	SkippedMessageNumberKeysYield func(number uint64, key keys.Message) bool
)

type SkippedKeysStorage interface {
	// Add must add new skipped keys to storage.
	Add(headerKey keys.Header, messageNumber uint64, messageKey keys.Message) error

	// Clone must deep clone a storage.
	Clone() SkippedKeysStorage

	// Delete must delete skipped keys by header key and message number.
	Delete(headerKey keys.Header, messageNumber uint64) error

	// GetIter must return function, which iterates over all skipped keys.
	GetIter() (SkippedKeysIter, error)
}

type skippedKeysStorage map[string]map[uint64]keys.Message

func (st skippedKeysStorage) Add(headerKey keys.Header, messageNumber uint64, messageKey keys.Message) error {
	stKey := string(headerKey.Bytes)
	if len(st[stKey]) >= skippedKeysStorageMessageKeysLenLimit {
		return fmt.Errorf("too many message keys: %d >= %d", len(st[stKey]), skippedKeysStorageMessageKeysLenLimit)
	}

	st.clearIfNeeded()

	if _, ok := st[stKey]; !ok {
		st[stKey] = make(map[uint64]keys.Message)
	}

	st[stKey][messageNumber] = messageKey

	return nil
}

func (st skippedKeysStorage) Clone() SkippedKeysStorage {
	stClone := make(skippedKeysStorage, len(st))

	for stKey, messageNumberKeys := range st {
		messageNumberKeysClone := make(map[uint64]keys.Message, len(messageNumberKeys))

		for messageNumber, messageKey := range messageNumberKeys {
			messageNumberKeysClone[messageNumber] = messageKey.Clone()
		}

		stClone[stKey] = messageNumberKeysClone
	}

	return stClone
}

func (st skippedKeysStorage) Delete(headerKey keys.Header, messageNumber uint64) error {
	delete(st[string(headerKey.Bytes)], messageNumber)

	return nil
}

func (st skippedKeysStorage) GetIter() (SkippedKeysIter, error) {
	iter := func(yield SkippedKeysYield) {
		for stKey, messageNumberKeys := range st {
			headerKey := keys.Header{Bytes: []byte(stKey)}

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

func (st skippedKeysStorage) clearIfNeeded() {
	if len(st) > skippedKeysStorageHeaderKeysLenToClear {
		for stKey := range st {
			delete(st, stKey)
		}
	}
}
