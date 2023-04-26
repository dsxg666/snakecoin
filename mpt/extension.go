package mpt

import (
	"github.com/ethereum/go-ethereum/crypto"
)

type ExtensionNode struct {
	Shared []Nibble
	Next   Node
}

func NewExtensionNode(nibbles []Nibble, next Node) *ExtensionNode {
	return &ExtensionNode{
		Shared: nibbles,
		Next:   next,
	}
}

func NewExtensionNodeWithDecodeData(ns []Nibble, is interface{}) *ExtensionNode {
	return &ExtensionNode{
		Shared: ns,
		Next:   NewBranchNodeWithDecodeData(is),
	}
}

func (e *ExtensionNode) Hash() []byte {
	return crypto.Keccak256(e.Serialize())
}

func (e *ExtensionNode) Serialize() []byte {
	return Serialize(e)
}

func (e *ExtensionNode) Raw() []interface{} {
	hashes := make([]interface{}, 2)
	hashes[0] = NibblesToBytes(AddPrefixedByIsLeafNode(e.Shared, false))
	hashes[1] = e.Next.Raw()

	return hashes
}
