package accounts

import (
	"github.com/dsxg666/snakecoin/core"
	"github.com/shopspring/decimal"
	"testing"
	"time"
)

func TestDecode(t *testing.T) {
	wallet := NewWallet()
	priStr, pubStr := Encode(&wallet.PrivateKey, &wallet.PublicKey)
	pubKey := PubDecode(pubStr)
	priKey := PriDecode(priStr)
	tx := core.NewTransaction(decimal.NewFromFloat(5), []byte("abcd"), []byte("efgh"), pubStr, time.Now().Unix())
	tx.Sign(priKey)
	if !tx.Verify(pubKey) {
		t.Error("unexpected occur")
	}
}
