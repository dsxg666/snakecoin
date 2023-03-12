package core

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"github.com/dsxg666/snakecoin/consensus"
	"log"
	"time"
)

// Block 区块结构体
type Block struct {
	Header *BlockHeader // 块头
	Body   *BlockBody   // 块身
}

// BlockHeader 块头结构体
type BlockHeader struct {
	Number                  int64                // 区块编号
	Timestamp               int64                // 挖矿当前时间的时间戳
	Nonce                   int64                // 随机数
	MiningTimestamp         int64                // 挖矿所用时间的时间戳
	Miner                   []byte               // 矿工地址
	Hash                    []byte               // 自己的哈希
	MerkleTreeRootHash      []byte               // 默克尔树的根哈希
	PreviousBlockHeaderHash []byte               // 前一区块的块头哈希
	Difficulty              consensus.Difficulty // 挖矿难度
}

// BlockBody 块身体结构体
type BlockBody struct {
	Txs []*Transaction // 交易列表
}

func NewBlock(number, bit int64, previousBlockHeaderHash, miner []byte, txs []*Transaction) *Block {
	merkleTree := NewMerkleTree(txs)
	blockHeader := &BlockHeader{
		Number:                  number,
		Timestamp:               time.Now().Unix(),
		Miner:                   miner,
		MerkleTreeRootHash:      merkleTree.RootNode.Hash,
		PreviousBlockHeaderHash: previousBlockHeaderHash,
		Difficulty:              consensus.Difficulty{Bits: bit},
	}
	blockBody := &BlockBody{
		Txs: txs,
	}
	block := &Block{blockHeader, blockBody}
	pow := consensus.NewProofOfWork(block.Header.Difficulty.Bits)
	miningTimestamp, nonce, hash := pow.Mine(number, block.Header.Timestamp, miner, merkleTree.RootNode.Hash, previousBlockHeaderHash)
	block.Header.MiningTimestamp = miningTimestamp
	block.Header.Nonce = nonce
	block.Header.Hash = hash
	return block
}

func NewGenesisBlock() *Block {
	sum := sha256.Sum256([]byte("Hello, I'm SnakeCoin!"))
	return &Block{
		&BlockHeader{
			Number:                  0,
			Timestamp:               time.Now().Unix(),
			Nonce:                   0,
			MiningTimestamp:         0,
			Miner:                   []byte("0000000000000000000000000000000000"),
			Hash:                    sum[:],
			MerkleTreeRootHash:      nil,
			PreviousBlockHeaderHash: nil,
			Difficulty:              consensus.Difficulty{Bits: 24}, // 设定的初始难度
		},
		&BlockBody{
			Txs: nil,
		},
	}
}

// Serialize 将 Block 结构体序列化
func (b *Block) Serialize() []byte {
	var buf bytes.Buffer
	// 得到编码器
	encoder := gob.NewEncoder(&buf)
	// 进行编码
	err := encoder.Encode(b)
	if err != nil {
		log.Panic(err)
	}
	return buf.Bytes()
}

// DeserializeBlock 将 []byte 反序列化为 Block 结构体
func DeserializeBlock(buf []byte) *Block {
	var block Block
	// 得到解码器
	decoder := gob.NewDecoder(bytes.NewReader(buf))
	// 进行解码
	err := decoder.Decode(&block)
	if err != nil {
		log.Panic(err)
	}
	return &block
}
