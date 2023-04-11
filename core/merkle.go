package core

import (
	"bytes"
	"crypto/sha256"
	"fmt"
)

type MerkleTree struct {
	RootNode *MerkleNode
}

type MerkleNode struct {
	Left  *MerkleNode
	Right *MerkleNode
	Hash  []byte
}

func NewMerkleNode(left, right *MerkleNode, tx *Transaction) *MerkleNode {
	node := new(MerkleNode)
	// If left and right both equal nil, then indicates that the current node is a leaf node.
	if left == nil && right == nil {
		node.Hash = tx.Hash()
	} else {
		combiendHash := append(left.Hash, right.Hash...)
		hash := sha256.Sum256(combiendHash)
		node.Hash = hash[:]
	}
	node.Right = right
	node.Left = left
	return node
}

func NewMerkleTree(txs []*Transaction) *MerkleTree {
	var nodes []*MerkleNode
	// Ensure that the number of nodes is an integer multiple of 2.
	if len(txs)%2 != 0 {
		txs = append(txs, txs[len(txs)-1])
	}
	for _, tx := range txs {
		node := NewMerkleNode(nil, nil, tx)
		nodes = append(nodes, node)
	}
	// A two-layer for loop is used to complete the tree construction of nodes
	for i := 0; i < len(txs)/2; i++ {
		if len(nodes)%2 != 0 {
			nodes = append(nodes, nodes[len(nodes)-1])
		}
		var newLevel []*MerkleNode
		for j := 0; j < len(nodes); j += 2 {
			node := NewMerkleNode(nodes[j], nodes[j+1], nil)
			newLevel = append(newLevel, node)
		}
		nodes = newLevel
	}
	return &MerkleTree{nodes[0]}
}

func PreOrderTraversal(root *MerkleNode) {
	if root == nil {
		return
	} else {
		PrintNode(root)
	}
	PreOrderTraversal(root.Left)
	PreOrderTraversal(root.Right)
}

func PrintNode(node *MerkleNode) {
	fmt.Printf("%p\n", node)
	if node != nil {
		fmt.Printf("left[%p], right[%p], data(%x)\n", node.Left, node.Right, node.Hash)
		fmt.Printf("check:%t\n\n", check(node))
	}
}

func check(node *MerkleNode) bool {
	if node.Left == nil {
		return true
	}
	prevHashes := append(node.Left.Hash, node.Right.Hash...)
	hash32 := sha256.Sum256(prevHashes)
	hash := hash32[:]
	return bytes.Compare(hash, node.Hash) == 0
}
