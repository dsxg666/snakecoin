package account

import (
	"fmt"
	"testing"
)

func TestNewAddress(t *testing.T) {
	fmt.Println(NewAddress().Hex())
}
