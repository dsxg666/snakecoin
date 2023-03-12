package main

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

// Mine 挖矿处理逻辑
func Mine(account string, txDB, accountDB, chainDB *pebble.DB, logger *zap.SugaredLogger) {
	txNum := core.NumOfTx(txDB)
	if txNum == 0 {
		color.Yellow("Mining requires at least one transaction!")
		fmt.Println()
		return
	}
	var txs []*core.Transaction
	// 从交易池获取交易
	for i := 0; i < txNum; i++ {
		txB := db.Get(util.StringToBytes(strconv.FormatInt(int64(i), 10)), txDB)
		txT := core.DeserializeTransaction(txB)
		pubKey := accounts.PubDecode(txT.PublicKey)
		// 验证签名
		if !txT.Verify(pubKey) {
			color.Red("Transaction verification failed!")
			fmt.Println()
			return
		}
		txs = append(txs, txT)
		txT.State = "Writed"
		db.Set(util.StringToBytes(strconv.FormatInt(int64(i), 10)), txT.Serialize(), txDB)
	}
	accB := db.Get(util.StringToBytes(account), accountDB)
	acc := accounts.DeserializeAccount(accB)
	acc.Balance = acc.Balance.Add(decimal.NewFromFloat(10))
	db.Set(acc.Address, acc.Serialize(), accountDB)
	var bc core.Blockchain
	fmt.Println("Please wait a moment, now is digging blocks for you.")
	bc.AddBlock(account, chainDB, txs)
	fmt.Println()
	s := strings.Split(util.CurrentTimeFormant(), " ")
	color.Green("INFO [%s|%s] Successfully digging into a block, you will receive 10 skc.", s[0], s[1])
	logger.Infof("INFO [%s|%s] Successfully digging into a block, you will receive 10 skc.", s[0], s[1])
	fmt.Println()
	meetAgain(account)
	fmt.Println()
}
