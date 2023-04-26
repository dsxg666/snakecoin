package wallet

import (
	"crypto/ecdsa"
	"crypto/rand"
	"reflect"
	"testing"
)

func TestEcdsa(t *testing.T) {
	priKey, pubKey := NewKeys()
	message := []byte("hello world")
	sign, _ := ecdsa.SignASN1(rand.Reader, priKey, message)
	if !ecdsa.VerifyASN1(pubKey, message, sign) {
		t.Errorf("Some unexpected happened.")
	}
	if ecdsa.VerifyASN1(pubKey, []byte("hello dsxg"), sign) {
		t.Errorf("Some unexpected happened.")
	}
}

func TestEncodeAndDecode(t *testing.T) {
	priKey, pubKey := NewKeys()

	pubKeyBytes := EncodePubKey(pubKey)
	priKeyBytes := EncodePriKey(priKey)

	pubKey2 := DecodePubKey(pubKeyBytes)
	priKey2 := DecodePriKey(priKeyBytes)
	if !reflect.DeepEqual(priKey, priKey2) {
		t.Errorf("Private keys do not match.")
	}
	if !reflect.DeepEqual(pubKey, pubKey2) {
		t.Errorf("Public keys do not match.")
	}
}
