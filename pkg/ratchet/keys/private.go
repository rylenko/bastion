package keys

import "github.com/rylenko/bastion/pkg/ratchet/utils"

type Private struct {
	Bytes []byte
}

func (pk Private) Clone() Private {
	pk.Bytes = utils.CloneByteSlice(pk.Bytes)
	return pk
}
