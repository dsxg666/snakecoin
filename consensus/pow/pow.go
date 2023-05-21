package pow

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"log"
	"math/big"
	"time"

	"github.com/dsxg666/snakecoin/common"
)

func Pow(diff *big.Int, data []byte) (*big.Int, *big.Int) {
	nonce := new(big.Int)
	begin := time.Now().UnixNano()
	for !Mine(diff, nonce, data) {
		nonce.Add(nonce, common.Big1)
	}
	end := time.Now().UnixNano()
	consumedTime := (end - begin) / 1e6
	if consumedTime < 10000 {
		diff.Add(diff, big.NewInt(16384))
	} else {
		diff.Sub(diff, big.NewInt(16384))
	}
	return nonce, diff
}

func Mine(diff, nonce *big.Int, data []byte) bool {
	t := append(data, nonce.Bytes()...)
	res := common.Bytes2BigInt(sha256.Sum256(t))
	temp := new(big.Int).Exp(common.Big2, common.Big256, nil) // 2**256
	target := new(big.Int).Div(temp, diff)
	return target.Cmp(res) > 0
}

func CombinedData(bss ...[]byte) []byte {
	data := bytes.Join(
		bss,
		[]byte{},
	)
	return data
}

func ToBytes(i interface{}) []byte {
	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.BigEndian, i)
	if err != nil {
		log.Panic("Failed to Write:", err)
	}
	return buf.Bytes()
}
