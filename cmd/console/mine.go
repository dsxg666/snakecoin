package console

import (
	"fmt"
	"github.com/cockroachdb/pebble"
	"github.com/dsxg666/snakecoin/accounts"
	"github.com/dsxg666/snakecoin/core"
	"github.com/dsxg666/snakecoin/db"
	"github.com/dsxg666/snakecoin/util"
	"github.com/fatih/color"
	"github.com/shopspring/decimal"
	"go.uber.org/zap"
	"strconv"
	"strings"
)

func Mine(num, txNum int, account string, txDB, accountDB, chainDB *pebble.DB, logger *zap.SugaredLogger) {
	var txs []*core.Transaction
	t := txNum
	for i := 0; i < num; i++ {
		t--
		txB := db.Get(util.StringToBytes(strconv.FormatInt(int64(t), 10)), txDB)
		txT := core.DeserializeTransaction(txB)
		pubKey := accounts.PubDecode(txT.PublicKey)
		if !txT.Verify(pubKey) {
			color.Red("Transaction verification failed!")
			fmt.Println()
			return
		}
		txs = append(txs, txT)
		txT.State = "Writed"
		db.Set(util.StringToBytes(strconv.FormatInt(int64(t), 10)), txT.Serialize(), txDB)
	}

	var bc core.Blockchain
	fmt.Println("Please wait a moment, now is digging blocks for you.")
	bc.AddBlock(account, chainDB, txs)
	accB := db.Get(util.StringToBytes(account), accountDB)
	acc := accounts.DeserializeAccount(accB)
	acc.Balance = acc.Balance.Add(decimal.NewFromFloat(10))
	db.Set(acc.Address, acc.Serialize(), accountDB)
	s := strings.Split(util.CurrentTimeFormant(), " ")
	color.Green("INFO [%s|%s] Successfully digging into a block, you will receive 10 skc.", s[0], s[1])
	logger.Infof("INFO [%s|%s] Successfully digging into a block, you will receive 10 skc.", s[0], s[1])
	fmt.Println()
	color.Blue("Welcome back to the SnakeCoin Blockchain console!")
	fmt.Println()
	fmt.Printf("CurrentAccountAddress: %s\n", account)
	fmt.Println("You can enter the following instruction to use blockchain:")
	fmt.Println("- [ transaction ] Conduct a transfer transaction")
	fmt.Println("- [ txpool ] You can view the situation in the txpool")
	fmt.Println("- [ mining ] Enter the mining program")
	fmt.Println("- [ blockchain ] See block information")
	fmt.Println("- [ balance ] Check your account balance")
	fmt.Println()
	fmt.Println("To exit, input quit")
	fmt.Println()
}
