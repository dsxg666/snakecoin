package account

import (
	"crypto/ecdsa"
	"crypto/rand"
	"reflect"
	"testing"
)

func TestEcdsa(t *testing.T) {
	priKey, pubKey := NewKeys()
	message := []byte("hello world")
	r, s, _ := ecdsa.Sign(rand.Reader, priKey, message)
	if !ecdsa.Verify(pubKey, message, r, s) {
		t.Errorf("Some unexpected happened.")
	}
	if ecdsa.Verify(pubKey, []byte("hello dsxg"), r, s) {
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
