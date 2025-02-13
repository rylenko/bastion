package keys

import "github.com/rylenko/bastion/pkg/ratchet/utils"

type Root struct {
	Bytes []byte
}

func (rk Root) Clone() Root {
	return Root{Bytes: utils.CloneByteSlice(rk.Bytes)}
}
