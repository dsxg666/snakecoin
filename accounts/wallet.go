package accounts

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"log"
)

// Wallet 钱包结构体
type Wallet struct {
	PrivateKey ecdsa.PrivateKey // 私钥
	PublicKey  ecdsa.PublicKey  // 公钥
}

// NewWallet 创建钱包
func NewWallet() *Wallet {
	//随机生成秘钥对
	private, public := newKeyPair()
	wallet := &Wallet{private, public}
	return wallet
}

// 生成公私钥函数
func newKeyPair() (ecdsa.PrivateKey, ecdsa.PublicKey) {
	// 获得椭圆曲线
	curve := elliptic.P256()
	// 生成私钥
	privateKey, err := ecdsa.GenerateKey(curve, rand.Reader)
	if err != nil {
		log.Panic(err)
	}
	publicKey := privateKey.PublicKey
	return *privateKey, publicKey
}
