package mptrie

import (
	"fmt"
	"log"
	"testing"

	"github.com/dsxg666/snakecoin/rlp"
	"github.com/stretchr/testify/require"
)

func ExampleSerialize() {
	key := []byte{0xa7, 0x11, 0x35, 0x51}
	value := []byte("45")
	key2 := []byte{0xa7, 0x7d, 0x33, 0x71}
	value2 := []byte("1")
	key3 := []byte{0xa7, 0xf9, 0x36, 0x51}
	value3 := []byte("2")
	key4 := []byte{0xa7, 0x7d, 0x39, 0x71}
	value4 := []byte("3")
	trie := NewTrie()
	trie.Put(key, value)
	trie.Put(key2, value2)
	trie.Put(key3, value3)
	trie.Put(key4, value4)
	i := Serialize(trie.Root)
	var e []interface{}
	err := rlp.DecodeBytes(i, &e)
	if err != nil {
		log.Panic("Failed to DecodeBytes:", err)
	}
	trie2 := NewTrieWithDecodeData(e)
	v, _ := trie2.Get(key)
	fmt.Println(v)
	// Output:
	// [52 53]
}

func TestPut_ShouldReturnLeafWhenTrieIsEmpty(t *testing.T) {
	key := []byte("account.address")
	value := []byte("XYZABCDEF")

	trie := NewTrie()
	err := trie.Put(key, value)
	require.NoError(t, err)

	require.IsType(t, &LeafNode{}, trie.Root)

	leaf := trie.Root.(*LeafNode)
	expectedNibbles := BytesToNibbles(key)
	require.Equal(t, value, leaf.Value)
	require.Equal(t, expectedNibbles, leaf.Suffix)
}

func TestPut_ShouldReturnExtensionNode(t *testing.T) {
	trie := NewTrie()

	err := trie.Put([]byte("accounts.address"), []byte("some_fake_addresss"))
	require.NoError(t, err)
	require.IsType(t, &LeafNode{}, trie.Root)

	err = trie.Put([]byte("accounts.value"), []byte("9000"))
	require.NoError(t, err)
	require.IsType(t, &ExtensionNode{}, trie.Root)
}

func TestPut_ShouldReturnLeafWhenUpdateValue(t *testing.T) {
	trie := NewTrie()

	err := trie.Put([]byte("accounts.address"), []byte("some_fake_addresss"))
	require.NoError(t, err)
	require.IsType(t, &LeafNode{}, trie.Root)
	require.Equal(t, trie.Root.(*LeafNode).Value, []byte("some_fake_addresss"))

	err = trie.Put([]byte("accounts.address"), []byte("another_address"))
	require.NoError(t, err)
	require.IsType(t, &LeafNode{}, trie.Root)
	require.Equal(t, trie.Root.(*LeafNode).Value, []byte("another_address"))
}

func TestPut_ShouldReturnBranchNodeWhenThereIsNoMatch(t *testing.T) {
	trie := NewTrie()

	err := trie.Put([]byte("accounts.balance"), []byte("10000"))
	require.NoError(t, err)

	err = trie.Put([]byte("system.version"), []byte("1.0.0.0"))
	require.NoError(t, err)

	require.IsType(t, &BranchNode{}, trie.Root)

	firstNibbles, secondNibbles := BytesToNibbles([]byte("accounts.balance")), BytesToNibbles([]byte("system.version"))
	branch := trie.Root.(*BranchNode)

	require.NotNil(t, branch.Branches[int(firstNibbles[0])])
	require.NotNil(t, branch.Branches[int(secondNibbles[0])])

	firstLeaf := branch.Branches[int(firstNibbles[0])].(*LeafNode)
	secondLeaf := branch.Branches[int(secondNibbles[0])].(*LeafNode)

	require.Equal(t, firstNibbles[1:], firstLeaf.Suffix)
	require.Equal(t, secondNibbles[1:], secondLeaf.Suffix)

	require.Equal(t, firstLeaf.Value, []byte("10000"))
	require.Equal(t, secondLeaf.Value, []byte("1.0.0.0"))
}

func TestPut_ShouldStoreLeafValueAtBranchNode(t *testing.T) {
	trie := NewTrie()

	err := trie.Put([]byte("transfer.input"), []byte("my-address"))
	require.NoError(t, err)

	err = trie.Put([]byte("transfer.input.value"), []byte("50"))
	require.NoError(t, err)

	require.IsType(t, &ExtensionNode{}, trie.Root)

	expectedExtNibbles := BytesToNibbles([]byte("transfer.input"))
	extnode := trie.Root.(*ExtensionNode)

	require.Equal(t, expectedExtNibbles, extnode.Shared)
	require.NotNil(t, extnode.Next)

	require.IsType(t, &BranchNode{}, extnode.Next)
	branch := extnode.Next.(*BranchNode)

	require.True(t, branch.HasValue())
	require.Equal(t, []byte("my-address"), branch.Value)
}

func TestPut_ShouldStoreNewValueAtBranchNode(t *testing.T) {
	trie := NewTrie()

	err := trie.Put([]byte("transfer.input.value"), []byte("50"))
	require.NoError(t, err)

	err = trie.Put([]byte("transfer.input"), []byte("my-address"))
	require.NoError(t, err)

	require.IsType(t, &ExtensionNode{}, trie.Root)

	expectedExtNibbles := BytesToNibbles([]byte("transfer.input"))
	extnode := trie.Root.(*ExtensionNode)

	require.Equal(t, expectedExtNibbles, extnode.Shared)
	require.NotNil(t, extnode.Next)

	require.IsType(t, &BranchNode{}, extnode.Next)
	branch := extnode.Next.(*BranchNode)

	require.True(t, branch.HasValue())
	require.Equal(t, []byte("my-address"), branch.Value)
}

func TestPut_ShouldReturnErrWhenKeyEmpty(t *testing.T) {
	trie := NewTrie()
	err := trie.Put([]byte(""), []byte(""))

	require.Error(t, err)
}

func TestPut_WhenRootIsBranchNodeWithEmptySlot_ShouldAddLeafNode(t *testing.T) {
	trie := NewTrie()

	err := trie.Put([]byte("zirst_info"), []byte("address"))
	require.NoError(t, err)

	err = trie.Put([]byte("other_info"), []byte("8000"))
	require.NoError(t, err)

	require.IsType(t, &BranchNode{}, trie.Root)

	err = trie.Put([]byte("10hird_info"), []byte("some-hash"))
	require.NoError(t, err)

	require.IsType(t, &BranchNode{}, trie.Root)
	branch := trie.Root.(*BranchNode)

	firstNibbles := BytesToNibbles([]byte("prevote"))
	secondNibbles := BytesToNibbles([]byte("transfer.input"))
	thirdNibbles := BytesToNibbles([]byte("block.header"))

	require.NotNil(t, branch.Branches[int(firstNibbles[0])])
	require.NotNil(t, branch.Branches[int(secondNibbles[0])])
	require.NotNil(t, branch.Branches[int(thirdNibbles[0])])

	require.IsType(t, &LeafNode{}, branch.Branches[int(firstNibbles[0])])
	require.IsType(t, &LeafNode{}, branch.Branches[int(secondNibbles[0])])
	require.IsType(t, &LeafNode{}, branch.Branches[int(thirdNibbles[0])])
}

func TestPut_WhenRootIsExtensionNodeShouldAddNewLeafNode(t *testing.T) {
	trie := NewTrie()

	err := trie.Put([]byte("block.header"), []byte("some_hash"))
	require.NoError(t, err)

	err = trie.Put([]byte("block.number"), []byte("1"))
	require.NoError(t, err)

	require.IsType(t, &ExtensionNode{}, trie.Root)

	newNibbles := BytesToNibbles([]byte("block.number"))
	leafNibbles := BytesToNibbles([]byte("block.header"))

	matched := PrefixMatchedLen(newNibbles, leafNibbles)

	extnode := trie.Root.(*ExtensionNode)
	require.Equal(t, leafNibbles[:matched], extnode.Shared)
	require.IsType(t, &BranchNode{}, extnode.Next)

	err = trie.Put([]byte("transfer.input"), []byte("1000"))
	require.NoError(t, err)

	require.IsType(t, &BranchNode{}, trie.Root)
	txNibbles := BytesToNibbles([]byte("transfer.input"))

	branch := trie.Root.(*BranchNode)
	txLeafNode := branch.Branches[int(txNibbles[0])]

	require.IsType(t, &LeafNode{}, txLeafNode)

	extNode := branch.Branches[int(extnode.Shared[0])]
	require.IsType(t, &ExtensionNode{}, extNode)
}

func TestPut_WhenExtensionDoesntHaveRemaining(t *testing.T) {
	firstKey, firstValue := []byte("transfer.to"), []byte("some-addr")
	secondKey, secondValue := []byte("transfer.input"), []byte("some-value")
	thirdKey, thirdValue := []byte("transfer.gas"), []byte("some-fee")

	trie := NewTrie()

	err := trie.Put(firstKey, firstValue)

	require.NoError(t, err)
	require.IsType(t, &LeafNode{}, trie.Root)

	err = trie.Put(secondKey, secondValue)

	require.NoError(t, err)
	require.IsType(t, &ExtensionNode{}, trie.Root)

	err = trie.Put(thirdKey, thirdValue)

	require.NoError(t, err)
	require.IsType(t, &ExtensionNode{}, trie.Root)
}

func TestGet_ShouldReturnFalse_WhenPrefixDoesntMathc(t *testing.T) {
	trie := NewTrie()

	// get when trie.Root is nil
	v, ok := trie.Get([]byte("doesnt exits"))
	require.False(t, ok)
	require.Nil(t, v)

	err := trie.Put([]byte("some-key"), []byte("some-value"))
	require.NoError(t, err)

	err = trie.Put([]byte("some-other-key"), []byte("other-value"))
	require.NoError(t, err)

	// get when extension node path doesn match
	v, ok = trie.Get([]byte("doesnt exits"))
	require.False(t, ok)
	require.Nil(t, v)

	err = trie.Put([]byte("another-key"), []byte("another-value"))
	require.NoError(t, err)

	// get when branch is Root
	v, ok = trie.Get([]byte("doesnt exits"))
	require.False(t, ok)
	require.Nil(t, v)
}

func TestGet_WhenValueExists(t *testing.T) {
	trie := NewTrie()

	err := trie.Put([]byte("some-key"), []byte("some-value"))
	require.NoError(t, err)

	v, ok := trie.Get([]byte("some-key"))
	require.True(t, ok)
	require.NotNil(t, v)
	require.Equal(t, v, []byte("some-value"))

	err = trie.Put([]byte("some-other-key"), []byte("other-value"))
	require.NoError(t, err)

	v, ok = trie.Get([]byte("some-other-key"))
	require.True(t, ok)
	require.NotNil(t, v)
	require.Equal(t, v, []byte("other-value"))

	err = trie.Put([]byte("another-key"), []byte("another-value"))
	require.NoError(t, err)

	v, ok = trie.Get([]byte("another-key"))
	require.True(t, ok)
	require.NotNil(t, v)
	require.Equal(t, v, []byte("another-value"))

	v, ok = trie.Get([]byte("some-key"))
	require.True(t, ok)
	require.NotNil(t, v)
	require.Equal(t, v, []byte("some-value"))
}
