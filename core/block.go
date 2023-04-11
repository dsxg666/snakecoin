package core

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"log"
	"time"

	"github.com/cockroachdb/pebble"
	"github.com/dsxg666/snakecoin/common"
	"github.com/dsxg666/snakecoin/consensus/pow"
	"github.com/dsxg666/snakecoin/db"
	"github.com/ethereum/go-ethereum/crypto"
)

type Block struct {
	Header *BlockHeader
	Body   *BlockBody
}

type BlockHeader struct {
	BlockHash      common.Hash
	PrevBlockHash  common.Hash
	StateTreeRoot  common.Hash
	MerkleTreeRoot common.Hash
	Coinbase       common.Address
	Difficulty     uint64
	Number         uint64
	GasLimit       uint64
	GasUsed        uint64
	Nonce          uint64
	Time           uint64
}

type BlockBody struct {
	Txs []*Transaction
}

func NewBlock(coinbase common.Address, chainDB, mptDB *pebble.DB, txs []*Transaction) {
	lastBlockHash := db.Get([]byte{byte(99)}, chainDB)
	lastBlockBytes := db.Get(lastBlockHash, chainDB)
	lastBlock := DeserializeBlock(lastBlockBytes)
	merkleTree := NewMerkleTree(txs)
	number := lastBlock.Header.Number + 1
	block := &Block{
		&BlockHeader{
			Number:   number,
			Time:     uint64(time.Now().Unix()),
			Coinbase: coinbase,
		},
		&BlockBody{
			Txs: txs,
		},
	}
	mptBytes := db.Get([]byte{byte(99)}, mptDB)
	db.Set([]byte{byte(number)}, mptBytes, mptDB)
	block.Header.StateTreeRoot.SetBytes(crypto.Keccak256(mptBytes))
	block.Header.MerkleTreeRoot.SetBytes(merkleTree.RootNode.Hash)
	block.Header.PrevBlockHash.SetBytes(lastBlock.Header.BlockHash.Bytes())
	fmt.Println("Mining is underway now, please wait patiently.")
	nonce, diff := pow.Pow(lastBlock.Header.Difficulty, pow.CombinedData(
		pow.ToBytes(block.Header.Number),
		pow.ToBytes(block.Header.Time),
		block.Header.Coinbase.Bytes(),
		block.Header.PrevBlockHash.Bytes(),
		block.Header.MerkleTreeRoot.Bytes(),
		block.Header.StateTreeRoot.Bytes(),
	))
	block.Header.Difficulty = diff
	block.Header.Nonce = nonce
	blockHash := block.Hash()
	block.Header.BlockHash.SetBytes(blockHash)
	db.Set(blockHash, block.Serialize(), chainDB)
	db.Set([]byte{byte(99)}, blockHash, chainDB)
}

func NewGenesisBlock(chainDB *pebble.DB) {
	block := &Block{
		&BlockHeader{
			Time:       uint64(time.Now().Unix()),
			Difficulty: 2195456,
		},
		&BlockBody{},
	}
	blockHash := block.Hash()
	block.Header.BlockHash.SetBytes(blockHash)
	db.Set(blockHash, block.Serialize(), chainDB)
	db.Set([]byte{byte(99)}, blockHash, chainDB)
}

func (b *Block) Hash() []byte {
	return crypto.Keccak256(b.Serialize())
}

func (b *Block) Serialize() []byte {
	var buf bytes.Buffer
	encoder := gob.NewEncoder(&buf)
	err := encoder.Encode(b)
	if err != nil {
		log.Panic("Failed to Encode:", err)
	}
	return buf.Bytes()
}

func DeserializeBlock(b []byte) *Block {
	var block Block
	decoder := gob.NewDecoder(bytes.NewReader(b))
	err := decoder.Decode(&block)
	if err != nil {
		log.Panic("Failed to Decode:", err)
	}
	return &block
}
