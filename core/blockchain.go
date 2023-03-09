package core

import (
	"github.com/cockroachdb/pebble"
	"github.com/dsxg666/snakecoin/db"
	"github.com/dsxg666/snakecoin/util"
)

// Blockchain 区块链结构体
type Blockchain struct{}

// AddBlock 增加一个区块
func (bc *Blockchain) AddBlock(account string, chainDB *pebble.DB, txs []*Transaction) {
	previousHash := db.Get([]byte("last"), chainDB)
	previousBlockB := db.Get(previousHash, chainDB)
	previousBlock := DeserializeBlock(previousBlockB)
	newBlock := NewBlock(previousBlock.Header.Number+1, previousHash, util.StringToBytes(account), txs)
	db.Set(newBlock.Header.Hash, newBlock.Serialize(), chainDB)
	db.Set([]byte("last"), newBlock.Header.Hash, chainDB)
}

// NewBlockchain 初始化区块链，首个区块为创世区块
func NewBlockchain(chainDB *pebble.DB) {
	genesis := NewGenesisBlock()
	db.Set(genesis.Header.Hash, genesis.Serialize(), chainDB)
	db.Set([]byte("last"), genesis.Header.Hash, chainDB)
}
