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

func Pow(diff uint64, data []byte) (uint64, uint64) {
	var nonce uint64
	begin := time.Now().UnixNano()
	for !mine(2195456, nonce, data) {
		nonce++
	}
	end := time.Now().UnixNano()
	consumedTime := (end - begin) / 1e6
	if consumedTime < 10000 {
		diff += 16384
	} else {
		diff -= 16384
	}
	return nonce, diff
}

func mine(diff, nonce uint64, data []byte) bool {
	t := append(data, ToBytes(nonce)...)
	res := common.Bytes2BigInt(sha256.Sum256(t))

	bigDiff := big.NewInt(int64(diff))
	temp := new(big.Int).Exp(common.Big2, common.Big256, nil) // 2**256
	target := new(big.Int).Div(temp, bigDiff)
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
