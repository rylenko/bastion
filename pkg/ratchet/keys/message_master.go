package keys

import "github.com/rylenko/bastion/pkg/ratchet/utils"

type MessageMaster struct {
	Bytes []byte
}

func (mk *MessageMaster) ClonePtr() *MessageMaster {
	if mk == nil {
		return nil
	}

	return &MessageMaster{Bytes: utils.CloneByteSlice(mk.Bytes)}
}
