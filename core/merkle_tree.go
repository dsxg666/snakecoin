package core

import (
	"bytes"
	"crypto/sha256"
	"fmt"
)

// MerkleTree 默克尔树结构体
type MerkleTree struct {
	RootNode *MerkleNode
}

// MerkleNode 默克尔树节点结构体
type MerkleNode struct {
	Left  *MerkleNode
	Right *MerkleNode
	Hash  []byte
}

// NewMerkleNode 创建节点
func NewMerkleNode(left, right *MerkleNode, tx *Transaction) *MerkleNode {
	node := &MerkleNode{}
	// 如果 left 和 right 都为空，则表示当前创建的节点是叶子节点
	if left == nil && right == nil {
		node.Hash = tx.Serialize()
	} else {
		// 表示当前创建的节点是中间节点
		// 将左右子树的 hash 组合在一起
		combinedHash := append(left.Hash, right.Hash...)
		hash := sha256.Sum256(combinedHash)
		node.Hash = hash[:]
	}
	// 给左右子树赋值
	node.Right = right
	node.Left = left
	return node
}

// NewMerkleTree 将节点组建成树
func NewMerkleTree(txs []*Transaction) *MerkleTree {
	var nodes []*MerkleNode
	for _, tx := range txs {
		node := NewMerkleNode(nil, nil, tx)
		nodes = append(nodes, node)
	}
	// 用双层for循环完成节点构造
	for i := 0; i < len(txs)/2; i++ {
		var newLevel []*MerkleNode
		// 将左右子节点进行合并
		for j := 0; j < len(nodes); j += 2 {
			node := NewMerkleNode(nodes[j], nodes[j+1], nil)
			newLevel = append(newLevel, node)
		}
		// 将 nodes 更新为合并后的新节点
		nodes = newLevel
	}
	return &MerkleTree{nodes[0]}
}

// ShowMerkleTree 先序遍历
func ShowMerkleTree(root *MerkleNode) {
	if root == nil {
		return
	} else {
		fmt.Println(root)
	}
	ShowMerkleTree(root.Left)
	ShowMerkleTree(root.Right)
}

func Check(node *MerkleNode) bool {
	if node.Left == nil {
		return true
	}
	prevHashes := append(node.Left.Hash, node.Right.Hash...)
	hash32 := sha256.Sum256(prevHashes)
	hash := hash32[:]
	return bytes.Compare(hash, node.Hash) == 0
}
