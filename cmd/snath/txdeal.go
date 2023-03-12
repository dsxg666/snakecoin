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
	"time"
)

func Transaction(account string, logger *zap.SugaredLogger, txDB, accountDB *pebble.DB) {
	b, loc := core.IsFull(txDB)
	if b {
		color.Yellow("Honey, the current transaction pool is full and cannot be traded.")
	} else {
		color.Blue("Welcome to the Snath transaction mode!")
		fmt.Println()
		accB := db.Get(util.StringToBytes(account), accountDB)
		acc := accounts.DeserializeAccount(accB)
		balance := acc.Balance
		b1 := true
		b2 := true
		for b1 {
			fmt.Println("Which account do you want to transfer money to?")
			var to string
			fmt.Print("> ")
			fmt.Scanf("%s\n", &to)
			res := db.AccountIsInDB(account, to)
			if res == 1 {
				fmt.Println()
				for b2 {
					fmt.Println("How many skc do you want to transfer?")
					var v string
					fmt.Print("> ")
					fmt.Scanf("%s\n", &v)
					if util.StringIsDigit(v) {
						amount, _ := strconv.ParseFloat(v, 64)
						if balance.Sub(decimal.NewFromFloat(amount)).LessThan(decimal.NewFromFloat(0)) {
							color.Red("You don't have enough balance!")
							fmt.Println()
						} else {
							acc.Balance = balance.Sub(decimal.NewFromFloat(amount))
							db.Set(acc.Address, acc.Serialize(), accountDB)
							toB := db.Get(util.StringToBytes(to), accountDB)
							acc2 := accounts.DeserializeAccount(toB)
							acc2.Balance = acc2.Balance.Add(decimal.NewFromFloat(amount))
							db.Set(acc2.Address, acc2.Serialize(), accountDB)
							pubPath := db.KeystoreDataPath + "/" + account + "/publickey.txt"
							priPath := db.KeystoreDataPath + "/" + account + "/privatekey.txt"
							pubStr := util.FileOutput(pubPath)
							priStr := util.FileOutput(priPath)
							tx := core.NewTransaction(decimal.NewFromFloat(amount), acc.Address, acc2.Address, pubStr, time.Now().Unix())
							priKey := accounts.PriDecode(priStr)
							// 对交易签名
							tx.Sign(priKey)
							// 将交易放入交易池
							core.Posh(loc, tx, txDB)
							b2 = false
							b1 = false
							s := strings.Split(util.CurrentTimeFormant(), " ")
							color.Green("INFO [%s|%s] A transaction was made successfully.", s[0], s[1])
							logger.Infof("INFO [%s|%s] A transaction was made successfully.", s[0], s[1])
							fmt.Println()
							meetAgain(account)
							fmt.Println()
						}
					} else {
						color.Red("Please do not enter amounts with letters!")
						fmt.Println()
					}
				}
			} else if res == 0 {
				color.Red("Do not transfer money to yourself!")
				fmt.Println()
			} else if res == -1 {
				color.Red("Target account does not exist!")
				fmt.Println()
			}
		}
	}
}
