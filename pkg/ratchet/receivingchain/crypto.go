package receivingchain

import "github.com/rylenko/bastion/pkg/ratchet/keys"

type Crypto interface {
	AdvanceChain(masterKey keys.MessageMaster) (keys.MessageMaster, keys.Message, error)
}
