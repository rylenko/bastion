package keys

import "github.com/rylenko/bastion/pkg/ratchet/utils"

type Message struct {
	Bytes []byte
}

func (mk Message) Clone() Message {
	mk.Bytes = utils.CloneByteSlice(mk.Bytes)
	return mk
}
