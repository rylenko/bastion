package receivingchain

import (
	"fmt"

	"github.com/rylenko/bastion/pkg/ratchet/keys"
)

const skippedKeysStorageAddAtOnceLimit = 32

type SkippedKeysStorage interface {
	// Add must add new skipped keys to storage.
	//
	// Total count indicates the total count of additions that will be made at once. This is useful to only allocate memory
	// once or to limit clutter from a large number of missed keys.
	Add(totalAtOnceCount uint64, headerKey keys.Header, messageNumber uint64, messageKey keys.Message) error

	// Clone must deep clone a storage.
	Clone() SkippedKeysStorage
}

type skippedKeysStorage map[string]map[uint64]keys.Message

func newSkippedKeysStorage() SkippedKeysStorage {
	return make(skippedKeysStorage)
}

func (st skippedKeysStorage) Add(
	totalAtOnceCount uint64,
	headerKey keys.Header,
	messageNumber uint64,
	messageKey keys.Message,
) error {
	if totalAtOnceCount > skippedKeysStorageAddAtOnceLimit {
		return fmt.Errorf("total count limit: %d > %d", totalAtOnceCount, skippedKeysStorageAddAtOnceLimit)
	}

	key := string(headerKey.Bytes)
	if _, exists := st[key]; !exists {
		st[key] = make(map[uint64]keys.Message, totalAtOnceCount)
	}

	st[key][messageNumber] = messageKey

	return nil
}

func (st skippedKeysStorage) Clone() SkippedKeysStorage {
	// TODO
	return st
}
