package core

import (
	"github.com/cockroachdb/pebble"
	"github.com/dsxg666/snakecoin/db"
	"github.com/dsxg666/snakecoin/util"
	"strconv"
	"strings"
)

// PoolSize 交易池最大可以容纳8笔交易
const PoolSize = 8

func Posh(loc int, tx *Transaction, txDB *pebble.DB) {
	db.Set(util.StringToBytes(strconv.FormatInt(int64(loc), 10)), tx.Serialize(), txDB)
}

// IsFull 判断交易池是否满了，并返回空余的位置
func IsFull(txDB *pebble.DB) (bool, int) {
	for i := 0; i < PoolSize; i++ {
		txB := db.Get(util.StringToBytes(strconv.FormatInt(int64(i), 10)), txDB)
		if txB == nil {
			return false, i
		} else {
			txT := DeserializeTransaction(txB)
			if strings.Compare(txT.State, "Writed") == 0 {
				return false, i
			}
		}
	}
	return true, -1
}

func NumOfTx(txDB *pebble.DB) int {
	for i := 0; i < PoolSize; i++ {
		txB := db.Get(util.StringToBytes(strconv.FormatInt(int64(i), 10)), txDB)
		if txB == nil {
			return i
		} else {
			txT := DeserializeTransaction(txB)
			if strings.Compare(txT.State, "Writed") == 0 {
				return i
			}
		}
	}
	return PoolSize
}
