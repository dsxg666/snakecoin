package logic

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/dsxg666/snakecoin/common"
	"github.com/dsxg666/snakecoin/core"
	"github.com/dsxg666/snakecoin/db"
	"github.com/dsxg666/snakecoin/mpt"
	"github.com/dsxg666/snakecoin/rlp"
	"github.com/dsxg666/snakecoin/wallet"
	"github.com/fatih/color"
	"github.com/howeyc/gopass"
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
	// New and store Account
	w := wallet.NewWallet()
	fmt.Println("Your new account is locked with a password. Please give a password. Do not forget this password.")
	fmt.Print("Password: ")
	pass, _ := gopass.GetPasswd()
	fmt.Print("Repeat password: ")
	pass2, _ := gopass.GetPasswd()
	if !bytes.Equal(pass, pass2) {
		color.Yellow("Passwords do not match!")
		fmt.Println()
		return
	}
	path := db.KeystorePath + "/" + w.Address.Hex()
	w.StoreKey(path, pass)

	// Save state to mptrie
	mptBytes := db.Get([]byte("latest"), mptDB)
	var e []interface{}
	err = rlp.DecodeBytes(mptBytes, &e)
	if err != nil {
		log.Panic("Failed to DecodeBytes:", err)
	}
	trie := mpt.NewTrieWithDecodeData(e)
	s := core.NewState()
	err = trie.Put(w.Address.Bytes(), s.Serialize())
	if err != nil {
		log.Panic("Failed to Put:", err)
	}
	db.Set([]byte("latest"), mpt.Serialize(trie.Root), mptDB)

	// Prompt
	time := strings.Split(common.GetCurrentTime(), " ")
	color.Green("INFO [%s|%s] Account creation succeeded! address: %s", time[0], time[1], w.Address.Hex())
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
