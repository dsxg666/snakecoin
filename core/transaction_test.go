package core

import (
	"github.com/dsxg666/snakecoin/accounts"
	"github.com/shopspring/decimal"
	"testing"
	"time"
)

func TestTransaction(t *testing.T) {
	wallet := accounts.NewWallet()
	_, pub := accounts.Encode(&wallet.PrivateKey, &wallet.PublicKey)
	tx := NewTransaction(decimal.NewFromFloat(5.1), []byte("abcd"), []byte("efgh"), pub, time.Now().Unix())
	tx2 := NewTransaction(decimal.NewFromFloat(5.8), []byte("abcd"), []byte("efgh"), pub, time.Now().Unix())
	tx.Sign(&wallet.PrivateKey)
	tx2.Sign(&wallet.PrivateKey)
	if !tx.Verify(&wallet.PublicKey) {
		t.Error("unexpected occur")
	}
	if !tx2.Verify(&wallet.PublicKey) {
		t.Error("unexpected occur")
	}
	if tx.ID == tx2.ID {
		t.Error("unexpected occur")
	}
}
