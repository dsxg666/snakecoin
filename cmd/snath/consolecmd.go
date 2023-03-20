package main

import (
	"fmt"
	"github.com/dsxg666/snakecoin/accounts"
	"github.com/dsxg666/snakecoin/util"
	"strings"

	"github.com/cockroachdb/pebble"
	"github.com/dsxg666/snakecoin/console"
	"github.com/dsxg666/snakecoin/db"
	"github.com/dsxg666/snakecoin/logs"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

// consoleCmd 代表 console 命令
var consoleCmd = &cobra.Command{
	Use:   "console",
	Short: "Start an interactive environment",
	Long: `
The snath is an interactive shell for blockchain
runtime environment which can interactwith the
blockchain.`,
	Run: func(cmd *cobra.Command, args []string) {
		account, _ := cmd.Flags().GetString("account")
		if strings.Compare("", account) == 0 {
			color.Red("This command must specify an account!")
		} else {
			res := db.AccountIsExist(account)
			if res {
				logger := logs.InitLogger()
				accountDB := db.GetDB(db.AccountDataPath)
				txDB := db.GetDB(db.TxDataPath)
				chainDB := db.GetDB(db.ChainDataPath)
				Begin(logger, account, txDB, accountDB, chainDB)
				db.CloseDB(chainDB)
				db.CloseDB(txDB)
				db.CloseDB(accountDB)
			} else {
				color.Red("The account you entered does not exist!")
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(consoleCmd)
	consoleCmd.Flags().StringP("account", "a", "", "an account for mining and transferring money")
}

// Begin 为 console 命令提供一个DOS窗口用户界面
func Begin(logger *zap.SugaredLogger, account string, txDB, accountDB, chainDB *pebble.DB) {
	fmt.Println()
	console.FirstMeet(account)
	line := console.Line()
	b := true
	for b {
		if input, err := line.Prompt("> "); err != nil {
			fmt.Println()
			fmt.Println()
			b = false
			color.Blue("bye")
		} else {
			switch input {
			case "transaction":
				TransactionDeal(account, logger, txDB, accountDB)
			case "txpool":
				TxPoolDeal(account, txDB)
			case "mining":
				MineDeal(account, txDB, accountDB, chainDB, logger)
			case "blockchain":
				BlockchainDeal(account, chainDB)
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
}
