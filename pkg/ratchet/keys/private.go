package keys

import "github.com/rylenko/bastion/pkg/ratchet/utils"

type Private struct {
	Bytes []byte
}

func (pk Private) Clone() Private {
	return Private{Bytes: utils.CloneByteSlice(pk.Bytes)}
}
