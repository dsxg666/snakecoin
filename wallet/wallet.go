package wallet

import (
	"crypto/ecdsa"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcutil/hdkeychain"
	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/tyler-smith/go-bip39"
	"log"
)

const defaultHDPath = "m/44'/60'/0'/0/1"

type Wallet struct {
	Address  common.Address
	Keystore *Keystore
}

func NewWallet() *Wallet {
	mnemonic := createMnemonic()
	priKey := newPriKeyFromMnemonic(mnemonic)
	pubKey := derivePubKey(priKey)
	address := crypto.PubkeyToAddress(*pubKey)
	k := NewKeystore(priKey, address)
	return &Wallet{Address: address, Keystore: k}
}

func (w Wallet) StoreKey(pass string) {
	filename := w.Keystore.joinPath(w.Address.Hex())
	w.Keystore.storeKey(filename, &w.Keystore.Key, pass)
}

func createMnemonic() string {
	entropy, _ := bip39.NewEntropy(256)
	mnemonic, _ := bip39.NewMnemonic(entropy)
	return mnemonic
}

func newPriKeyFromMnemonic(mnemonic string) *ecdsa.PrivateKey {
	path, err := accounts.ParseDerivationPath(defaultHDPath)
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
		log.Panic("Failed to ECPrivKey")
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
