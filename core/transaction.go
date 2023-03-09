package core

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/sha256"
	"encoding/gob"
	"github.com/shopspring/decimal"
	"log"
	"math/big"
)

// Transaction 交易结构体
type Transaction struct {
	ID        [32]byte        // 交易ID
	State     string          // 当前状态：Wait 表示等待被加入区块，Writed 表示已经被写入区块
	Amount    decimal.Decimal // 转账数额
	Timestamp int64           // 交易时间
	From      []byte          // 交易输入项
	To        []byte          // 交易输出项
	PublicKey string          // 发送者的公钥
	R, S      big.Int         // 加密参数
}

func NewTransaction(amount decimal.Decimal, from, to []byte, publicKey string, timestamp int64) *Transaction {
	return &Transaction{
		State:     "Wait",
		Amount:    amount,
		Timestamp: timestamp,
		From:      from,
		To:        to,
		PublicKey: publicKey,
	}
}

func (tx *Transaction) SetID() {
	tx.ID = sha256.Sum256(tx.Serialize())
}

func (tx *Transaction) Sign(privateKey *ecdsa.PrivateKey) {
	tx.SetID()
	r, s, _ := ecdsa.Sign(rand.Reader, privateKey, tx.ID[:])
	tx.R = *r
	tx.S = *s
}

func (tx *Transaction) Verify(publicKey *ecdsa.PublicKey) bool {
	if ecdsa.Verify(publicKey, tx.ID[:], &tx.R, &tx.S) {
		return true
	} else {
		return false
	}
}

// Serialize 将 Transaction 结构体序列化
func (tx *Transaction) Serialize() []byte {
	var b bytes.Buffer
	// 得到编码器
	encoder := gob.NewEncoder(&b)
	// 进行编码
	err := encoder.Encode(tx)
	if err != nil {
		log.Panic(err)
	}
	return b.Bytes()
}

// DeserializeTransaction 将 []byte 反序列化为 Transaction 结构体
func DeserializeTransaction(b []byte) *Transaction {
	var tx Transaction
	// 得到解码器
	decoder := gob.NewDecoder(bytes.NewReader(b))
	// 进行解码
	err := decoder.Decode(&tx)
	if err != nil {
		log.Panic(err)
	}
	return &tx
}
