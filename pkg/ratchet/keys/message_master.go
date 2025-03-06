package keys

import "github.com/rylenko/bastion/pkg/ratchet/utils"

type MessageMaster struct {
	Bytes []byte
}

func (mk MessageMaster) Clone() MessageMaster {
	mk.Bytes = utils.CloneByteSlice(mk.Bytes)
	return mk
}

func (mk *MessageMaster) ClonePtr() *MessageMaster {
	if mk == nil {
		return nil
	}

	clone := mk.Clone()

	return &clone
}
