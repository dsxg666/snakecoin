package accounts

import (
	"crypto/sha256"
	"golang.org/x/crypto/ripemd160"
)

const version = byte(0x00)

const ChecksumLen = 4

// HashPublicKey 计算公钥hash
func HashPublicKey(pubKey []byte) []byte {
	// 先hash一次
	publicSHA256 := sha256.Sum256(pubKey)
	// 计算ripemd160
	RIPEMD160Hasher := ripemd160.New()
	RIPEMD160Hasher.Write(publicSHA256[:])

	publicRIPEMD160 := RIPEMD160Hasher.Sum(nil)

	return publicRIPEMD160
}

// Checksum 计算校验和，输入为0x00+公钥hash
func Checksum(payload []byte) []byte {

	firstSHA := sha256.Sum256(payload)
	secondSHA := sha256.Sum256(firstSHA[:])

	return secondSHA[:ChecksumLen]
}

// GetAddress 生成地址
func GetAddress(wallet *Wallet) []byte {
	private := &wallet.PrivateKey
	// 利用私钥推导出公钥
	pubKey := append(private.PublicKey.X.Bytes(), private.PublicKey.Y.Bytes()...)
	// 计算公钥hash
	pubKeyHash := HashPublicKey(pubKey)
	// 计算校验和
	versionedPayload := append([]byte{version}, pubKeyHash...)
	checksum := Checksum(versionedPayload)
	// 计算base58编码
	fullPayload := append(versionedPayload, checksum...)
	address := Base58Encode(fullPayload)

	return address
}
