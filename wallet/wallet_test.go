package wallet

import (
	"fmt"
	"testing"
)

func TestCreateMnemonic(t *testing.T) {
	fmt.Println(createMnemonic())
}

func TestNewWallet(t *testing.T) {
	w := NewWallet()
	fmt.Println(w.Address.Hex())
}
