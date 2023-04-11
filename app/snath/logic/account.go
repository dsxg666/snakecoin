package logic

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/dsxg666/snakecoin/account"
	"github.com/dsxg666/snakecoin/common"
	"github.com/dsxg666/snakecoin/core"
	"github.com/dsxg666/snakecoin/db"
	"github.com/dsxg666/snakecoin/rlp"
	mptrie "github.com/dsxg666/snakecoin/trie"
	"github.com/fatih/color"
)

func New() {
	_, err := os.Stat("./data/chain")
	if os.IsNotExist(err) {
		color.Red("The blockchain hasn't been initialized yet, please enter the following command to initialize:")
		fmt.Println(" - [snath init]")
		fmt.Println()
		return
	}
	// Get and close db
	mptDB := db.GetDB(db.MPTirePath)
	defer db.CloseDB(mptDB)
	acc := account.NewAccount()
	// Save account state to mptrie
	mptBytes := db.Get([]byte{byte(99)}, mptDB)
	var e []interface{}
	err = rlp.DecodeBytes(mptBytes, &e)
	if err != nil {
		log.Panic("Failed to DecodeBytes:", err)
	}
	trie := mptrie.NewTrieWithDecodeData(e)
	s := core.NewState()
	err = trie.Put(acc.Bytes(), s.Serialize())
	if err != nil {
		log.Panic("Failed to Put:", err)
	}
	db.Set([]byte{byte(99)}, mptrie.Serialize(trie.Root), mptDB)
	// Prompt
	time := strings.Split(common.GetCurrentTime(), " ")
	color.Green("INFO [%s|%s] Account creation succeeded! address: %s", time[0], time[1], acc.Hex())
	fmt.Println()
}

func List() {
	files, _ := filepath.Glob(db.KeystorePath + "/*")
	if len(files) == 0 {
		color.Red("You don't have an account yet, please enter the following command to new an account:")
		fmt.Println(" - [snath account new]")
		fmt.Println()
		return
	}
	for i := 0; i < len(files); i++ {
		fmt.Printf("Address%d: %s\n", i+1, files[i][14:])
	}
	fmt.Println()
}
