package mpt

import (
	"github.com/ethereum/go-ethereum/crypto"
)

type LeafNode struct {
	Suffix []Nibble
	Value  []byte
}

func NewLeafNode(nibbles []Nibble, value []byte) *LeafNode {
	return &LeafNode{
		Suffix: nibbles,
		Value:  value,
	}
}

func NewLeafNodeWithDecodeData(ns []Nibble, is interface{}) *LeafNode {
	v := is.([]byte)
	return &LeafNode{
		Suffix: ns,
		Value:  v,
	}
}

func (l *LeafNode) Hash() []byte {
	return crypto.Keccak256(l.Serialize())
}

func (l *LeafNode) Serialize() []byte {
	return Serialize(l)
}

func (l *LeafNode) Raw() []interface{} {
	path := NibblesToBytes(AddPrefixedByIsLeafNode(l.Suffix, true))
	raw := []interface{}{path, l.Value}
	return raw
}
