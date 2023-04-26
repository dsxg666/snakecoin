package mpt

import (
	"github.com/ethereum/go-ethereum/crypto"
)

type BranchNode struct {
	Branches [16]Node
	Value    []byte
}

func NewBranchNode() *BranchNode {
	return &BranchNode{
		Branches: [16]Node{},
	}
}

func NewBranchNodeWithDecodeData(is interface{}) *BranchNode {
	branch := new(BranchNode)
	its, _ := is.([]interface{})
	for i := 0; i < 16; i++ {
		if _, ok := its[i].([]byte); ok == true {
			branch.Branches[i] = nil
		}
		if is, ok := its[i].([]interface{}); ok == true {
			branch.Branches[i] = IsToNode(is)
		}
	}
	if v, ok := its[16].([]byte); ok == true {
		branch.Value = v
	}
	return branch
}

func (b *BranchNode) SetBranch(nb Nibble, n Node) {
	b.Branches[int(nb)] = n
}

func (b *BranchNode) SetValue(v []byte) {
	b.Value = v
}

func (b *BranchNode) Raw() []interface{} {
	hashes := make([]interface{}, 17)
	for i := 0; i < 16; i++ {
		if b.Branches[i] == nil {
			hashes[i] = EmptyNodeRaw
		} else {
			node := b.Branches[i]
			hashes[i] = node.Raw()
		}
	}

	hashes[16] = b.Value
	return hashes
}

func (b *BranchNode) Hash() []byte {
	return crypto.Keccak256(b.Serialize())
}

func (b *BranchNode) Serialize() []byte {
	return Serialize(b)
}

func (b *BranchNode) HasValue() bool {
	return b.Value != nil
}
