package logic

import (
	"fmt"
	"log"
	"path/filepath"
	"strings"

	"github.com/cockroachdb/pebble"
	"github.com/dsxg666/snakecoin/common"
	"github.com/dsxg666/snakecoin/core"
	"github.com/dsxg666/snakecoin/db"
	"github.com/dsxg666/snakecoin/mpt"
	"github.com/dsxg666/snakecoin/rlp"
)

func AccountIsExist(address string) bool {
	files, _ := filepath.Glob(db.KeystorePath + "/*")
	for i := 0; i < len(files); i++ {
		if strings.Compare(files[i][14:], address) == 0 {
			return true
		}
	}
	return false
}

func ShowAccountBalance(account string, mptDB *pebble.DB) {
	mptBytes := db.Get([]byte("latest"), mptDB)
	var e []interface{}
	err := rlp.DecodeBytes(mptBytes, &e)
	if err != nil {
		log.Panic("Failed to DecodeBytes:", err)
	}
	trie := mpt.NewTrieWithDecodeData(e)
	stateB, _ := trie.Get(common.Hex2Bytes(account[2:]))
	state := core.DeserializeState(stateB)
	fmt.Println(state.Balance.String(), "skc")
}
