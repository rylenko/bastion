package keys

import "github.com/rylenko/bastion/pkg/ratchet/utils"

type Public struct {
	Bytes []byte
}

func (pk Public) Clone() Public {
	pk.Bytes = utils.CloneByteSlice(pk.Bytes)
	return pk
}

func (pk *Public) ClonePtr() *Public {
	if pk == nil {
		return nil
	}

	return &Public{Bytes: utils.CloneByteSlice(pk.Bytes)}
}
