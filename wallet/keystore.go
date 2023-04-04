package wallet

import (
	"crypto/ecdsa"
	"crypto/rand"
	"github.com/dsxg666/snakecoin/db"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"io"
	"log"
	"os"
	"path/filepath"
)

type Keystore struct {
	KeystorePath string
	scryptN      int
	scryptP      int
	Key          keystore.Key
}

func NewKeystore(priKey *ecdsa.PrivateKey, address common.Address) *Keystore {
	uuid := []byte(newUUID())
	key := keystore.Key{
		Id:         uuid,
		Address:    address,
		PrivateKey: priKey,
	}
	return &Keystore{
		KeystorePath: db.KeystorePath,
		scryptN:      keystore.LightScryptN,
		scryptP:      keystore.LightScryptP,
		Key:          key,
	}
}

type UUID []byte

func newUUID() UUID {
	uuid := make([]byte, 16)
	_, _ = io.ReadFull(rand.Reader, uuid)
	uuid[6] = (uuid[6] & 0x0f) | 0x40
	uuid[8] = (uuid[8] & 0x3f) | 0x80
	return uuid
}

func (k Keystore) storeKey(filename string, key *keystore.Key, pass string) {
	keyjson, err := keystore.EncryptKey(key, pass, k.scryptN, k.scryptP)
	if err != nil {
		log.Panic("Failed to EncryptKey:", err)
	}

	writeKeyFile(filename, keyjson)
}

func writeKeyFile(filename string, content []byte) {
	const perm = 0700
	if err := os.MkdirAll(filepath.Dir(filename), perm); err != nil {
		log.Panic("Failed to MkdirAll:", err)
	}

	f, err := os.CreateTemp(filepath.Dir(filename), "."+filepath.Base(filename)+".tmp")
	if err != nil {
		log.Panic("Failed to CreateTemp:", err)
	}

	if _, err := f.Write(content); err != nil {
		_ = f.Close()
		_ = os.Remove(f.Name())
		log.Panic("Failed to Write:", err)
	}

	_ = f.Close()
	_ = os.Rename(f.Name(), filename)
}

func (k Keystore) joinPath(filename string) string {
	if filepath.IsAbs(filename) {
		return filename
	}
	return filepath.Join(k.KeystorePath, filename)
}
