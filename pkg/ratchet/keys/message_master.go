package keys

import "github.com/rylenko/bastion/pkg/ratchet/utils"

type MessageMaster struct {
	Bytes []byte
}

func (mk *MessageMaster) ClonePtr() *MessageMaster {
	return &MessageMaster{Bytes: utils.CloneByteSlice(mk.Bytes)}
}
