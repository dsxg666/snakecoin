package consensus

import (
	"bytes"
	"crypto/sha256"
	"github.com/dsxg666/snakecoin/util"
	"math/big"
	"strconv"
	"time"
)

// ProofOfWork pow结构体
type ProofOfWork struct {
	Target *big.Int // 目标域值
}

// Difficulty 挖矿难度结构体
type Difficulty struct {
	Bits int64 // 位移数
}

// NewProofOfWork 创建pow结构体
func NewProofOfWork(bits int64) *ProofOfWork {
	target := big.NewInt(1)
	// 相当于：target = target << (256-targetBits)
	target.Lsh(target, uint(256-bits))
	return &ProofOfWork{target}
}

// Mine 挖矿
func (pow *ProofOfWork) Mine(number, timestamp int64, miner, merkleTreeRootHash, previousBlockHeaderHash []byte) (int64, int64, []byte) {
	var hashInt big.Int
	var hash [32]byte
	var nonce int64
	t1 := time.Now().Unix()
	for {
		data := pow.PrepareData(nonce, number, timestamp, miner, merkleTreeRootHash, previousBlockHeaderHash)
		hash = sha256.Sum256(data)
		hashInt.SetBytes(hash[:])
		if hashInt.Cmp(pow.Target) == -1 {
			break
		} else {
			nonce++
		}
	}
	mixHash := append(hash[:], []byte(strconv.FormatInt(nonce, 10))...)
	finalHash := sha256.Sum256(mixHash)
	t2 := time.Now().Unix()
	return t2 - t1, nonce, finalHash[:]
}

func (pow *ProofOfWork) PrepareData(nonce, number, timestamp int64, miner, merkleTreeRootHash, previousBlockHeaderHash []byte) []byte {
	data := bytes.Join(
		[][]byte{
			merkleTreeRootHash,
			util.Int64ToBytes(timestamp),
			util.Int64ToBytes(number),
			miner,
			util.Int64ToBytes(nonce),
			previousBlockHeaderHash,
		},
		[]byte{},
	)
	return data
}
