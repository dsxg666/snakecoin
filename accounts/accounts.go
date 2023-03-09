package accounts

import (
	"bytes"
	"encoding/gob"
	"github.com/shopspring/decimal"
	"log"
)

// Account 账户的结构体
type Account struct {
	Address []byte          // 地址
	Balance decimal.Decimal // 余额
}

// NewAccount 创建一个账户
func NewAccount() (*Account, *Wallet) {
	wallet := NewWallet()
	address := GetAddress(wallet)
	return &Account{Address: address, Balance: decimal.NewFromFloat(10)}, wallet
}

func (a *Account) Serialize() []byte {
	var buf bytes.Buffer
	// 得到编码器
	encoder := gob.NewEncoder(&buf)
	// 进行编码
	err := encoder.Encode(a)
	if err != nil {
		log.Panic(err)
	}
	return buf.Bytes()
}

func DeserializeAccount(buf []byte) *Account {
	var account Account
	// 得到解码器
	decoder := gob.NewDecoder(bytes.NewReader(buf))
	// 进行解码
	err := decoder.Decode(&account)
	if err != nil {
		log.Panic(err)
	}
	return &account
}
