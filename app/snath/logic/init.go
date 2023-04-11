package logic

import (
	"fmt"
	"os"
	"strings"

	"github.com/dsxg666/snakecoin/common"
	"github.com/dsxg666/snakecoin/core"
	"github.com/dsxg666/snakecoin/db"
	mptrie "github.com/dsxg666/snakecoin/trie"
	"github.com/fatih/color"
)

func Init() {
	b := InitDir()
	if b {
		// Get and close db
		chainDB := db.GetDB(db.ChainPath)
		defer db.CloseDB(chainDB)
		mptDB := db.GetDB(db.MPTirePath)
		defer db.CloseDB(mptDB)
		// Blockchain initialization
		core.NewGenesisBlock(chainDB)
		// MPTire initialization
		trie := mptrie.NewTrie()
		trie.Put([]byte("hello"), []byte("world"))
		db.Set([]byte{byte(99)}, mptrie.Serialize(trie.Root), mptDB)
		// Prompt
		time := strings.Split(common.GetCurrentTime(), " ")
		color.Green("INFO [%s|%s] Initialization is successful!", time[0], time[1])
		fmt.Println("The data directory is generated for you in the current directory.")
		fmt.Println()
	} else {
		color.Red("The initialization is done! please init other command.")
		fmt.Println()
	}
}

func InitDir() bool {
	_, err := os.Stat("./data/chain")
	if os.IsNotExist(err) {
		_, err := os.Stat("./data")
		if os.IsNotExist(err) {
			_ = os.Mkdir("./data", 0777)
		}
		_ = os.Mkdir(db.MPTirePath, 0777)
		_ = os.Mkdir(db.ChainPath, 0777)
		_ = os.Mkdir(db.KeystorePath, 0777)
		_ = os.Mkdir(db.LogPath, 0777)
		_ = os.Mkdir(db.TxPath, 0777)
		return true
	} else {
		return false
	}
}
