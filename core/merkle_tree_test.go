package core

import (
	"crypto/sha256"
	"testing"
)

func TestNewMerkleTree(t *testing.T) {
	tests := []*Transaction{
		&Transaction{ID: sha256.Sum256([]byte("test1"))},
		&Transaction{ID: sha256.Sum256([]byte("test2"))},
	}
	tree := NewMerkleTree(tests)
	if !Check(tree.RootNode) {
		t.Error("unexpected occur")
	}
}
