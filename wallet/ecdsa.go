package wallet

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"encoding/pem"
)

func NewKeys() (*ecdsa.PrivateKey, *ecdsa.PublicKey) {
	priKey, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	return priKey, &priKey.PublicKey
}

func EncodePubKey(pubKey *ecdsa.PublicKey) []byte {
	x509EncodedPub, _ := x509.MarshalPKIXPublicKey(pubKey)
	return pem.EncodeToMemory(&pem.Block{Type: "PUBLIC KEY", Bytes: x509EncodedPub})
}

func DecodePubKey(b []byte) *ecdsa.PublicKey {
	blockPub, _ := pem.Decode(b)
	x509EncodedPub := blockPub.Bytes
	genericPublicKey, _ := x509.ParsePKIXPublicKey(x509EncodedPub)
	return genericPublicKey.(*ecdsa.PublicKey)
}

func EncodePriKey(priKey *ecdsa.PrivateKey) []byte {
	x509Encoded, _ := x509.MarshalECPrivateKey(priKey)
	return pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: x509Encoded})
}

func DecodePriKey(b []byte) *ecdsa.PrivateKey {
	block, _ := pem.Decode(b)
	if block == nil {
		return nil
	}
	x509Encoded := block.Bytes
	privateKey, _ := x509.ParseECPrivateKey(x509Encoded)
	return privateKey
}
