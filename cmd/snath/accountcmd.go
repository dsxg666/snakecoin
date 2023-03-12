package main

import (
	"fmt"
	"github.com/dsxg666/snakecoin/accounts"
	"github.com/dsxg666/snakecoin/db"
	"github.com/dsxg666/snakecoin/logs"
	"github.com/dsxg666/snakecoin/util"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
	"strings"
)

// accountCmd 代表 account 命令
var accountCmd = &cobra.Command{
	Use:   "account",
	Short: "Manage accounts",
	Long: `
Account command can create a new account and list all
your existing accounts.`,
}

var newCmd = &cobra.Command{
	Use:   "new",
	Short: "New account",
	Long: `
New command can create a new account for you.`,
	Run: func(cmd *cobra.Command, args []string) {
		_, err := os.Stat("./data")
		if os.IsNotExist(err) {
			color.Yellow("Honey, the blockchain hasn't been initialized yet, you can initialize by input [snakecoin init].")
			return
		}
		account, wallet := accounts.NewAccount()
		logger := logs.InitLogger()
		accDB := db.GetDB(db.AccountDataPath)
		db.Set(account.Address, account.Serialize(), accDB)
		db.CloseDB(accDB)
		_ = os.Mkdir(db.KeystoreDataPath+"/"+util.BytesToString(account.Address), 0666)
		priFilename := db.KeystoreDataPath + "/" + util.BytesToString(account.Address) + "/privatekey.txt"
		pubFilename := db.KeystoreDataPath + "/" + util.BytesToString(account.Address) + "/publickey.txt"
		addrFilename := db.KeystoreDataPath + "/" + util.BytesToString(account.Address) + "/address.txt"
		priStr, pubStr := accounts.Encode(&wallet.PrivateKey, &wallet.PublicKey)
		util.FileInput(priFilename, priStr)
		util.FileInput(pubFilename, pubStr)
		util.FileInput(addrFilename, util.BytesToString(account.Address))
		logger.Infof("Account creation succeeded! Address: %s", util.BytesToString(account.Address))
		s := strings.Split(util.CurrentTimeFormant(), " ")
		color.Green("INFO [%s|%s] Account creation succeeded! Address: %s", s[0], s[1], util.BytesToString(account.Address))
	},
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all accounts",
	Long: `
List command can list all accounts for you.`,
	Run: func(cmd *cobra.Command, args []string) {
		files, _ := filepath.Glob(db.KeystoreDataPath + "/*")
		if len(files) == 0 {
			color.Yellow("Honey, you don't have an account yet, You can create by input [snakecoin account new].")
			return
		}
		for i := 0; i < len(files); i++ {
			fmt.Printf("Address%d: %s\n", i+1, files[i][14:])
		}
	},
}

func init() {
	rootCmd.AddCommand(accountCmd)
	accountCmd.AddCommand(newCmd)
	accountCmd.AddCommand(listCmd)
}
