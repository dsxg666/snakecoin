package core

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/rand"
	"encoding/gob"
	"log"
	"math/big"

	"github.com/dsxg666/snakecoin/account"
	"github.com/dsxg666/snakecoin/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/shopspring/decimal"
)

type Transaction struct {
	TxHash common.Hash
	// State meaning:
	//     0 represent not yet included in the blockchain
	//     1 represent already included in the blockchain
	State  int
	Amount decimal.Decimal
	Time   uint64
	From   common.Address
	To     common.Address
	PubKey []byte
	R, S   big.Int
}

func NewTransaction(amount decimal.Decimal, time uint64, from, to common.Address, pubKey []byte) *Transaction {
	return &Transaction{
		Amount: amount,
		Time:   time,
		From:   from,
		To:     to,
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
	if tx.TxHash != [32]byte{} {
		return tx.TxHash.Bytes()
	} else {
		tx.TxHash.SetBytes(crypto.Keccak256(tx.Serialize()))
		return tx.TxHash.Bytes()
	}
}

func (tx *Transaction) Sign(privateKey *ecdsa.PrivateKey) {
	r, s, _ := ecdsa.Sign(rand.Reader, privateKey, tx.Hash())
	tx.R = *r
	tx.S = *s
}

func (tx *Transaction) Verity() bool {
	return ecdsa.Verify(account.DecodePubKey(tx.PubKey), tx.Hash(), &tx.R, &tx.S)
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
