package main

import (
	"fmt"
	"github.com/cockroachdb/pebble"
	"github.com/dsxg666/snakecoin/core"
	"github.com/dsxg666/snakecoin/db"
	"github.com/dsxg666/snakecoin/util"
	"github.com/fatih/color"
	"strconv"
	"strings"
)

func Blockchain(account string, chainDB *pebble.DB) {
	prevHash := db.Get([]byte("last"), chainDB)
	blockB := db.Get(prevHash, chainDB)
	block := core.DeserializeBlock(blockB)
	num := block.Header.Number
	if num == 0 {
		color.Yellow("There is no block in the current blockchain.")
		fmt.Println()
	} else {
		color.Blue("Welcome to the Snath blockchain mode!")
		fmt.Println()
		fmt.Printf("The height of the current block is %d.\n", num)
		fmt.Printf("You can specify the height of the block to find the block information.\n")
		fmt.Println()
		fmt.Println("To exit, input leave")
		b := true
		for b {
			fmt.Print("> ")
			var input string
			fmt.Scanf("%s\n", &input)
			if util.StringIs1ToN(input, int(num)) {
				var temp *core.Block
				temp = block
				in, _ := strconv.ParseInt(input, 10, 64)
				for i := 0; i < int(num-in); i++ {
					prevBlockHash := db.Get(temp.Header.PreviousBlockHeaderHash, chainDB)
					prevBLock := core.DeserializeBlock(prevBlockHash)
					temp = prevBLock
				}
				color.Green("BlockHeaderInfo")
				fmt.Printf("Number: %d\n", temp.Header.Number)
				fmt.Printf("Time: %s\n", util.TimestampFormat(temp.Header.Timestamp))
				fmt.Printf("Nonce: %d\n", temp.Header.Nonce)
				fmt.Printf("MiningTimestamp: %d\n", temp.Header.MiningTimestamp)
				fmt.Printf("Difficulty: %d\n", temp.Header.Difficulty.Bits)
				fmt.Printf("Miner: %s\n", util.BytesToString(temp.Header.Miner))
				fmt.Printf("Hash: %x\n", temp.Header.Hash)
				fmt.Printf("PreviousBlockHash: %x\n", temp.Header.PreviousBlockHeaderHash)
				fmt.Printf("MerkleTreeRootHash: %x\n", temp.Header.MerkleTreeRootHash)
				fmt.Println()
				color.Green("BlockBodyInfo")
				txs := temp.Body.Txs
				for i := 0; i < len(txs); i++ {
					fmt.Printf("Transaction%d:\n", i+1)
					fmt.Printf("ID: %x\n", txs[i].ID)
					fmt.Printf("From: %s\n", util.BytesToString(txs[i].From))
					fmt.Printf("To: %s\n", util.BytesToString(txs[i].To))
					fmt.Println("Amount:", txs[i].Amount, "skc")
					fmt.Printf("State: %s\n", txs[i].State)
					fmt.Printf("Time: %s\n", util.TimestampFormat(txs[i].Timestamp))
					fmt.Println()
				}
			} else if strings.Compare(input, "leave") == 0 {
				b = false
				fmt.Println()
				meetAgain(account)
				fmt.Println()
			} else {
				color.Red("Your input is not valid!")
				fmt.Println()
			}
		}
	}
}
