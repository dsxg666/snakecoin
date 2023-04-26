package core

import (
	"bytes"
	"encoding/gob"
	"log"
	"math/big"

	"github.com/dsxg666/snakecoin/common"
	"github.com/dsxg666/snakecoin/consensus/pow"
	"github.com/ethereum/go-ethereum/crypto"
)

type Transaction struct {
	TxHash    common.Hash
	From      common.Address
	To        common.Address
	Value     *big.Int
	Time      uint64
	PubKey    []byte
	Signature []byte
	State     int
	// State meaning:
	//     0 represent not yet included in the blockchain
	//     1 represent already included in the blockchain
}

func NewTransaction(amount *big.Int, time uint64, from, to common.Address, pubKey []byte) *Transaction {
	return &Transaction{
		From:   from,
		To:     to,
		Value:  amount,
		Time:   time,
		PubKey: pubKey,
	}
}

func (tx *Transaction) Serialize() []byte {
	var buf bytes.Buffer
	encoder := gob.NewEncoder(&buf)
	err := encoder.Encode(tx)
	if err != nil {
		log.Panic("Failed to Encode:", err)
	}
	return buf.Bytes()
}

func (tx *Transaction) Hash() []byte {
	return crypto.Keccak256(tx.From.Bytes(), tx.To.Bytes(), tx.Value.Bytes(), pow.ToBytes(tx.Time), tx.PubKey)
}

func DeserializeTx(b []byte) *Transaction {
	var tx Transaction
	decoder := gob.NewDecoder(bytes.NewReader(b))
	err := decoder.Decode(&tx)
	if err != nil {
		log.Panic("Failed to Decode:", err)
	}
	return &tx
}
