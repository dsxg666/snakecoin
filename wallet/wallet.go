package wallet

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/sha256"
	"log"

	"github.com/dsxg666/snakecoin/common"
)

type Wallet struct {
	PriKey  *ecdsa.PrivateKey
	PubKey  *ecdsa.PublicKey
	Address common.Address
}

func NewWallet() *Wallet {
	priKey, pubKey := NewKeys()
	w := &Wallet{PriKey: priKey, PubKey: pubKey}
	h := sha256.New()
	h.Write(w.PubKey.X.Bytes())
	h.Write(w.PubKey.Y.Bytes())
	digest := h.Sum(nil)
	w.Address.SetBytes(digest[:20])
	return w
}

func (w *Wallet) Sign(data []byte) []byte {
	signature, _ := ecdsa.SignASN1(rand.Reader, w.PriKey, data)
	return signature
}

func Verity(data, sign []byte, pubKey *ecdsa.PublicKey) bool {
	return ecdsa.VerifyASN1(pubKey, data, sign)
}

func (w *Wallet) StoreKey(filename string, pass []byte) {
	text, err := DesEncrypt(EncodePriKey(w.PriKey), TrimKey(pass))
	if err != nil {
		log.Panic("Failed to DesEncrypt:", err)
	}
	common.WriteFile(filename, text)
}

func LoadWallet(filename string, pass []byte, acc string) *Wallet {
	text := common.ReadFile(filename)
	priKeyBytes, err := DesDecrypt(string(text), TrimKey(pass))
	if err != nil {
		return nil
	}
	priKey := DecodePriKey(priKeyBytes)
	if priKey == nil {
		return nil
	}
	pubKey := &priKey.PublicKey
	h := sha256.New()
	h.Write(pubKey.X.Bytes())
	h.Write(pubKey.Y.Bytes())
	digest := h.Sum(nil)
	accByte := common.Hex2Bytes(acc[2:])
	if !bytes.Equal(accByte, digest[:20]) {
		return nil
	}
	w := &Wallet{PriKey: priKey, PubKey: pubKey}
	w.Address.SetBytes(accByte)
	return w
}
