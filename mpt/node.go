package mpt

import (
	"encoding/hex"
	"github.com/dsxg666/snakecoin/rlp"
)

var (
	EmptyNodeRaw     = []byte{}
	EmptyNodeHash, _ = hex.DecodeString("56e81f171bcc55a6ff8345e692c0f86e5b48e01b996cadc001622fb5e363b421")
)

type Node interface {
	Hash() []byte
	Raw() []interface{}
}

func Hash(n Node) []byte {
	if n == nil {
		return EmptyNodeHash
	}

	return n.Hash()
}

func Serialize(n Node) []byte {
	var raw interface{}

	if n == nil {
		raw = EmptyNodeRaw
	} else {
		raw = n.Raw()
	}

	b, err := rlp.EncodeToBytes(raw)
	if err != nil {
		panic(err)
	}

	return b
}
