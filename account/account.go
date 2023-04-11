package account

import (
	"crypto/ecdsa"
	"log"
	"os"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcutil/hdkeychain"
	util "github.com/dsxg666/snakecoin/common"
	"github.com/dsxg666/snakecoin/db"
	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/tyler-smith/go-bip39"
)

func NewAccount() common.Address {
	acc := NewAddress()
	priKey, pubKey := NewKeys()
	path := db.KeystorePath + "/" + acc.Hex()
	_ = os.Mkdir(path, 0777)
	util.WriteFile(path+"/private", string(EncodePriKey(priKey)))
	util.WriteFile(path+"/public", string(EncodePubKey(pubKey)))
	return acc
}

func NewAddress() common.Address {
	mnemonic := createMnemonic()
	priKey := newPriKeyFromMnemonic(mnemonic)
	pubKey := derivePubKey(priKey)
	return crypto.PubkeyToAddress(*pubKey)
}

func createMnemonic() string {
	entropy, _ := bip39.NewEntropy(256)
	mnemonic, _ := bip39.NewMnemonic(entropy)
	return mnemonic
}

func newPriKeyFromMnemonic(mnemonic string) *ecdsa.PrivateKey {
	path, err := accounts.ParseDerivationPath("m/44'/60'/0'/0/1")
	if err != nil {
		log.Panic("Failed to ParseDerivationPath:", err)
	}

	seed := bip39.NewSeed(mnemonic, "Secret Passphrase")
	masterKey, err := hdkeychain.NewMaster(seed, &chaincfg.MainNetParams)
	if err != nil {
		log.Panic("Failed to NewMaster:", err)
	}

	return derivePriKey(path, masterKey)
}

func derivePriKey(path accounts.DerivationPath, masterKey *hdkeychain.ExtendedKey) *ecdsa.PrivateKey {
	var err error
	key := masterKey

	for _, n := range path {
		key, err = key.Child(n)
		if err != nil {
			log.Panic("Failed to Child:", err)
		}
	}

	priKey, err := key.ECPrivKey()
	if err != nil {
		log.Panic("Failed to ECPrivKey:", err)
	}

	return priKey.ToECDSA()
}

func derivePubKey(priKey *ecdsa.PrivateKey) *ecdsa.PublicKey {
	pubKey := priKey.Public()
	pubKeyECDSA, ok := pubKey.(*ecdsa.PublicKey)
	if !ok {
		log.Panic("Failed to get pub key")
	}
	return pubKeyECDSA
}
