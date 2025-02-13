package keys

import "github.com/rylenko/bastion/pkg/ratchet/utils"

type Header struct {
	Bytes []byte
}

func (hk Header) Clone() Header {
	return Header{Bytes: utils.CloneByteSlice(hk.Bytes)}
}

func (hk *Header) ClonePtr() *Header {
	if hk == nil {
		return nil
	}

	return &Header{Bytes: utils.CloneByteSlice(hk.Bytes)}
}
