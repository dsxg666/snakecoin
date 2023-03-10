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
	newBlock := NewBlock(previousBlock.Header.Number+1, previousBlock.Header.Difficulty.Bits, previousHash, util.StringToBytes(account), txs)
	// 每隔5个区块动态的调整难度，使挖矿的平均时间等于10s
	if newBlock.Header.Number%5 == 0 {
		var temp *Block
		temp = newBlock
		actualTime := newBlock.Header.MiningTimestamp
		for i := 0; i < 4; i++ {
			prevBlockHash := db.Get(temp.Header.PreviousBlockHeaderHash, chainDB)
			prevBLock := DeserializeBlock(prevBlockHash)
			actualTime += prevBLock.Header.MiningTimestamp
			temp = prevBLock
		}
		newBlock.Header.Difficulty.Bits *= int64(actualTime / 50)
	}
	db.Set(newBlock.Header.Hash, newBlock.Serialize(), chainDB)
	db.Set([]byte("last"), newBlock.Header.Hash, chainDB)

}

// NewBlockchain 初始化区块链，首个区块为创世区块
func NewBlockchain(chainDB *pebble.DB) {
	genesis := NewGenesisBlock()
	db.Set(genesis.Header.Hash, genesis.Serialize(), chainDB)
	db.Set([]byte("last"), genesis.Header.Hash, chainDB)
}
