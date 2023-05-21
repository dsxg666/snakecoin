package wallet

import (
	"fmt"
	"log"
	"testing"
)

func Example() {
	strEncrypted, err := DesEncrypt([]byte("hello world"), []byte("12345678"))
	if err != nil {
		log.Fatal(err)
	}
	strDecrypted, err := DesDecrypt(strEncrypted, []byte("12345678"))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Encrypted:", strEncrypted)
	fmt.Println("Decrypted:", string(strDecrypted))
	//Output:
	//Encrypted: 28dba02eb5f6dd476042daebfa59687a
	//Decrypted: hello world
}

func TestTrimKey(t *testing.T) {
	tests := []struct {
		b   []byte
		ans int
	}{
		{[]byte("1234"), 8},
		{[]byte(";]anba"), 8},
		{[]byte("12345678"), 8},
		{[]byte("123456789lx,s;"), 8},
	}
	for _, tt := range tests {
		fmt.Println(TrimKey(tt.b))
		if actual := TrimKey(tt.b); len(actual) != tt.ans {
			t.Errorf("TrimKey(%s) expected %d, but got %d", string(tt.b), tt.ans, len(actual))
		}
	}
}
