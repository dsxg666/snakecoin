package console

import (
	"fmt"
	"github.com/cockroachdb/pebble"
	"github.com/dsxg666/snakecoin/accounts"
	"github.com/dsxg666/snakecoin/core"
	"github.com/dsxg666/snakecoin/db"
	"github.com/dsxg666/snakecoin/util"
	"github.com/fatih/color"
	"go.uber.org/zap"
)

func Start(logger *zap.SugaredLogger, account string, txDB, accountDB, chainDB *pebble.DB) {
	fmt.Println()
	color.Blue("Welcome to the SnakeCoin console!")
	fmt.Println()
	fmt.Printf("CurrentAccountAddress: %s\n", account)
	fmt.Println("You can enter the following instruction to use blockchain:")
	fmt.Println("- [ transaction ] Conduct a transfer transaction")
	fmt.Println("- [ txpool ] Look at the transactions in the pool")
	fmt.Println("- [ mining ] Enter the mining program")
	fmt.Println("- [ blockchain ] See block information")
	fmt.Println("- [ balance ] Check your account balance")
	fmt.Println()
	fmt.Println("To exit, input quit")
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
			num := core.NumOfTx(txDB)
			if num < 2 {
				color.Yellow("A minimum of two transactions in the pool are required to run the mining process")
				fmt.Println()
			} else if num >= 2 && num < 4 {
				Mine(2, num, account, txDB, accountDB, chainDB, logger)
			} else if num >= 4 && num < 8 {
				Mine(4, num, account, txDB, accountDB, chainDB, logger)
			} else if num == 8 {
				Mine(8, num, account, txDB, accountDB, chainDB, logger)
			}
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
