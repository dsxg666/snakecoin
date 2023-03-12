package main

import (
	"fmt"
	"github.com/cockroachdb/pebble"
	"github.com/dsxg666/snakecoin/accounts"
	"github.com/dsxg666/snakecoin/db"
	"github.com/dsxg666/snakecoin/util"
	"github.com/fatih/color"
	"go.uber.org/zap"
)

// Begin 为 console 命令提供一个DOS窗口用户界面
func Begin(logger *zap.SugaredLogger, account string, txDB, accountDB, chainDB *pebble.DB) {
	fmt.Println()
	firstMeet(account)
	b := true
	for b {
		fmt.Print("> ")
		var input string
		fmt.Scanf("%s\n", &input)
		switch input {
		case "transaction":
			Transaction(account, logger, txDB, accountDB)
		case "txpool":
			TxPool(account, txDB)
		case "mining":
			Mine(account, txDB, accountDB, chainDB, logger)
		case "blockchain":
			Blockchain(account, chainDB)
		case "balance":
			b := db.Get(util.StringToBytes(account), accountDB)
			acc := accounts.DeserializeAccount(b)
			fmt.Println(acc.Balance, "skc")
			fmt.Println()
		case "quit":
			b = false
			fmt.Println()
			color.Blue("bye")
		default:
			color.Yellow("Honey, very sorry, we don't support your instruction yet.")
		}
	}
}
