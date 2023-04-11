package mptrie

import (
	"errors"
)

type Trie struct {
	Root Node
}

func NewTrie() *Trie {
	return new(Trie)
}

func NewTrieWithDecodeData(is []interface{}) *Trie {
	trie := NewTrie()
	trie.Root = IsToNode(is)
	return trie
}

func IsToNode(is []interface{}) Node {
	if len(is) == 2 {
		bs := is[0].([]byte)
		ns := BytesToNibbles(bs)
		if ns[0] == 0 {
			// Even Nibble ExtensionNode
			return NewExtensionNodeWithDecodeData(ns[2:], is[1])
		}
		if ns[0] == 1 {
			// Odd Nibble ExtensionNode
			return NewExtensionNodeWithDecodeData(ns[1:], is[1])
		}
		if ns[0] == 2 {
			// Even Nibble LeafNode
			return NewLeafNodeWithDecodeData(ns[2:], is[1])
		}
		if ns[0] == 3 {
			// Odd Nibble LeafNode
			return NewLeafNodeWithDecodeData(ns[1:], is[1])
		}
	}
	if len(is) == 17 {
		return NewBranchNodeWithDecodeData(is)
	}
	return nil
}

func (t *Trie) Hash() []byte {
	if t.Root == nil {
		return EmptyNodeHash
	}

	return t.Root.Hash()
}

func (t *Trie) Get(key []byte) ([]byte, bool) {
	node := t.Root
	nibbles := BytesToNibbles(key)

	for {
		if node == nil {
			return nil, false
		}

		if leaf, ok := node.(*LeafNode); ok {
			matched := PrefixMatchedLen(nibbles, leaf.Suffix)
			if matched != len(nibbles) || matched != len(leaf.Suffix) {
				return nil, false
			}

			return leaf.Value, true
		}

		if branch, ok := node.(*BranchNode); ok {
			if len(nibbles) == 0 {
				return branch.Value, branch.HasValue()
			}

			b, remaining := nibbles[0], nibbles[1:]
			nibbles = remaining

			node = branch.Branches[b]
			continue
		}

		if ext, ok := node.(*ExtensionNode); ok {
			matched := PrefixMatchedLen(ext.Shared, nibbles)

			if matched < len(ext.Shared) {
				return nil, false
			}

			nibbles = nibbles[matched:]
			node = ext.Next
			continue
		}

		return nil, false
	}
}

func (t *Trie) Update(key, value []byte) bool {
	node := t.Root
	nibbles := BytesToNibbles(key)

	for {
		if node == nil {
			return false
		}

		if leaf, ok := node.(*LeafNode); ok {
			matched := PrefixMatchedLen(nibbles, leaf.Suffix)
			if matched != len(nibbles) || matched != len(leaf.Suffix) {
				return false
			}
			leaf.Value = value
			return true
		}

		if branch, ok := node.(*BranchNode); ok {
			if len(nibbles) == 0 {
				branch.Value = value
				return branch.HasValue()
			}

			b, remaining := nibbles[0], nibbles[1:]
			nibbles = remaining

			node = branch.Branches[b]
			continue
		}

		if ext, ok := node.(*ExtensionNode); ok {
			matched := PrefixMatchedLen(ext.Shared, nibbles)

			if matched < len(ext.Shared) {
				return false
			}

			nibbles = nibbles[matched:]
			node = ext.Next
			continue
		}

		return false
	}
}

// Put inserts a key -> value in the merkle tree
// EmptyNode     -> replace with a leaf node with the path
// LeafNode      -> transform into a Extension Node add a new branch node and a new leaf node
// ExtensionNode -> convert to a Extension Node with a shorter path, create a branch node that points to a new Extension Node
func (t *Trie) Put(key, value []byte) error {
	node := &t.Root
	nibbles := BytesToNibbles(key)

	if len(nibbles) <= 0 {
		return errors.New("cannot insert empty keys")
	}

	for {
		if *node == nil {
			leaf := NewLeafNode(nibbles, value)
			*node = leaf
			return nil
		}

		if leaf, ok := (*node).(*LeafNode); ok {
			matched := PrefixMatchedLen(leaf.Suffix, nibbles)

			// all the leaf.Shared matches with nibbles then update the value
			if matched == len(leaf.Suffix) && matched == len(nibbles) {
				newleaf := NewLeafNode(leaf.Suffix, value)
				*node = newleaf
				return nil
			}

			branch := NewBranchNode()

			if matched == len(leaf.Suffix) {
				branch.SetValue(leaf.Value)
			}

			if matched == len(nibbles) {
				branch.SetValue(value)
			}

			// if there is matched nibbles, an extension node will be created
			if matched > 0 {
				ext := NewExtensionNode(leaf.Suffix[:matched], branch)
				*node = ext
			} else {
				*node = branch
			}

			if matched < len(leaf.Suffix) {
				branchNibble, leafNibbles := leaf.Suffix[matched], leaf.Suffix[matched+1:]
				newLeaf := NewLeafNode(leafNibbles, leaf.Value)

				branch.SetBranch(branchNibble, newLeaf)
			}

			if matched < len(nibbles) {
				branchNibble, leafNode := nibbles[matched], nibbles[matched+1:]
				newLeaf := NewLeafNode(leafNode, value)

				branch.SetBranch(branchNibble, newLeaf)
			}

			return nil
		}

		if branch, ok := (*node).(*BranchNode); ok {
			branchNibble, remaining := nibbles[0], nibbles[1:]
			nibbles = remaining
			node = &branch.Branches[int(branchNibble)]

			continue
		}

		if ext, ok := (*node).(*ExtensionNode); ok {
			matched := PrefixMatchedLen(ext.Shared, nibbles)

			if matched < len(ext.Shared) {
				extNibbles, branchNibble, extRemaining := ext.Shared[:matched], ext.Shared[matched], ext.Shared[matched+1:]
				newBranchNibble, newLeafNibbles := nibbles[matched], nibbles[matched+1:]

				branch := NewBranchNode()
				if len(extRemaining) == 0 {
					branch.SetBranch(branchNibble, ext.Next)
				} else {
					newExt := NewExtensionNode(extRemaining, ext.Next)
					branch.SetBranch(branchNibble, newExt)
				}

				newleaf := NewLeafNode(newLeafNibbles, value)
				branch.SetBranch(newBranchNibble, newleaf)

				if len(extNibbles) == 0 {
					*node = branch
				} else {
					*node = NewExtensionNode(extNibbles, branch)
				}

				return nil
			}
			nibbles = nibbles[matched:]
			node = &ext.Next
			continue
		}
	}
}
