package keys

import "github.com/rylenko/bastion/pkg/utils"

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

	clone := pk.Clone()

	return &clone
}
