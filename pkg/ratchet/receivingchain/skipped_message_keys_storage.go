package receivingchain

type SkippedMessageKeysStorage interface {
	Clone() SkippedMessageKeysStorage
}

type skippedMessageKeysStorage struct{}

func newSkippedMessageKeysStorage() SkippedMessageKeysStorage {
	return skippedMessageKeysStorage{}
}

func (st skippedMessageKeysStorage) Clone() SkippedMessageKeysStorage {
	// TODO
	return st
}
