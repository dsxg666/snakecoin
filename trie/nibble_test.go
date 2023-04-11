package mptrie

import (
	"encoding/binary"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIsNibble(t *testing.T) {
	bs := []byte{'a', 'b', 0xff, 3, 44}
	ns := BytesToNibbles(bs)

	for _, n := range ns {
		require.True(t, IsNibble(byte(n)))
	}
}

func TestByteToNibbles(t *testing.T) {
	b := byte(0xff) // represents 11111111 (255) when do b >> 4 should return 00001111 1 + 2 + 4 + 8 (15)
	expected := []Nibble{15, 15}

	ns := ByteToNibbles(b)

	require.Equal(t, expected, ns)
	require.Len(t, ns, 2)
}

func TestBytesToNibbles(t *testing.T) {
	bs := []byte{0xff, 0xff, 0xff, 0xff}
	ns := BytesToNibbles(bs)
	expected := []Nibble{15, 15, 15, 15, 15, 15, 15, 15}

	require.Equal(t, expected, ns)
	require.Len(t, ns, len(bs)*2)
}

func TestNibblesToBytes(t *testing.T) {
	b := byte(0xff)        // 11111111
	ns := ByteToNibbles(b) // []Nibble{15, 15}

	bs := NibblesToBytes(ns) // 15 (00001111) << 4 = 11110000 (240) + 15 = 255 (0xff)

	require.Len(t, bs, 1)
	require.Equal(t, b, bs[0])
}

func TestAddPrefixedByIsLeafNode(t *testing.T) {
	b := byte(0xff)
	ns := ByteToNibbles(b)

	prefixedNibbles := AddPrefixedByIsLeafNode(ns, true)
	prefixedNibbles2 := AddPrefixedByIsLeafNode(ns, false)
	expected := []Nibble{2, 0, 15, 15}
	expected2 := []Nibble{0, 0, 15, 15}

	require.Len(t, prefixedNibbles, 4)
	require.Equal(t, expected, prefixedNibbles)
	require.Len(t, prefixedNibbles2, 4)
	require.Equal(t, expected2, prefixedNibbles2)
}

func TestPrefixMatchedLen(t *testing.T) {
	key1 := []byte{'a', 'a', 'c'}
	key2 := []byte{'a', 'a'}

	expectMatchedLen := 4 // [a, b, c] [a, b] == 2 * 2 (each byte = 2 nibbles)
	value := PrefixMatchedLen(BytesToNibbles(key1), BytesToNibbles(key2))

	require.Equal(t, expectMatchedLen, value)

	putKey := BytesToNibbles([]byte("account.Balance"))
	putValue := make([]byte, 4)
	binary.LittleEndian.PutUint32(putValue, 100)

	leafn := NewLeafNode(putKey, putValue)

	equalPrefixed := PrefixMatchedLen(leafn.Suffix, putKey)
	require.Equal(t, len(putKey), equalPrefixed)
	require.Equal(t, len(leafn.Suffix), equalPrefixed)
}
