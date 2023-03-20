package main

import (
	"fmt"
	"strings"

	"github.com/cockroachdb/pebble"
	"github.com/dsxg666/snakecoin/console"
	"github.com/dsxg666/snakecoin/core"
	"github.com/dsxg666/snakecoin/db"
	"github.com/dsxg666/snakecoin/util"
	"github.com/fatih/color"
)

func TxPoolDeal(account string, txDB *pebble.DB) {
	num := core.NumOfTx(txDB)
	if num == 0 {
		color.Yellow("There is no transaction in the pool.")
		fmt.Println()
	} else {
		color.Blue("Welcome to the Snath txpool mode!")
		fmt.Println()
		color.Yellow("Note that the pool has a capacity of only 8, meaning it can only hold 8 transactions.")
		fmt.Printf("There are now %d transactions in the pool.\n", num)
		fmt.Printf("You can enter %d numbers from 0-%d to see the status of a given transaction.\n", num, num-1)
		fmt.Println()
		fmt.Println("To exit, input leave")
		b := true
		for b {
			fmt.Print("> ")
			var input string
			fmt.Scanf("%s\n", &input)
			if util.StringIs0ToN(input, num) {
				txB := db.Get(util.StringToBytes(input), txDB)
				tx := core.DeserializeTransaction(txB)
				fmt.Printf("ID: %x\n", tx.ID)
				fmt.Printf("From: %s\n", util.BytesToString(tx.From))
				fmt.Printf("To: %s\n", util.BytesToString(tx.To))
				fmt.Println("Amount:", tx.Amount, "skc")
				fmt.Printf("State: %s\n", tx.State)
				fmt.Printf("Time: %s\n", util.TimestampFormat(tx.Timestamp))
				fmt.Println()
			} else if strings.Compare(input, "leave") == 0 {
				b = false
				fmt.Println()
				console.MeetAgain(account)
				fmt.Println()
			} else {
				color.Red("Your input is not valid!")
				fmt.Println()
			}
		}
	}
}
