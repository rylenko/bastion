package receivingchain

type SkippedMessageKeysStorage interface{}

type skippedMessageKeysStorage struct{}

func newSkippedMessageKeysStorage() SkippedMessageKeysStorage {
	return &skippedMessageKeysStorage{}
}
