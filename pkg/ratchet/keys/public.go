package keys

import "github.com/rylenko/bastion/pkg/ratchet/utils"

type Public struct {
	Bytes []byte
}

func (pk Public) Clone() Public {
	return Public{Bytes: utils.CloneByteSlice(pk.Bytes)}
}

func (pk *Public) ClonePtr() *Public {
	if pk == nil {
		return nil
	}

	return &Public{Bytes: utils.CloneByteSlice(pk.Bytes)}
}
