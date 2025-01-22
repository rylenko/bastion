package chains

import "github.com/rylenko/sapphire/pkg/shield/keys"

// Root is the root chain of the participant of conversation. The main task of the root chain is to generate new keys
// of the sending and receiving chains and their next header encryption keys.
type Root struct {
	key *keys.Root
}

// NewRoot creates a new instance of root chain.
func NewRoot(key *keys.Root) *Root {
	return &Root{key: key}
}
