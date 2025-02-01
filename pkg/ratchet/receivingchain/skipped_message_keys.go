package receivingchain

type SkippedMessageKeys interface{}

type skippedMessageKeys struct{}

func newSkippedMessageKeys() SkippedMessageKeys {
	return &skippedMessageKeys{}
}
