package keys

import "github.com/rylenko/bastion/pkg/utils"

type Root struct {
	Bytes []byte
}

func (rk Root) Clone() Root {
	rk.Bytes = utils.CloneByteSlice(rk.Bytes)
	return rk
}
