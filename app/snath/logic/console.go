package logic

import (
	"crypto/ecdsa"
	"fmt"
	"log"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/cockroachdb/pebble"
	"github.com/dsxg666/snakecoin/account"
	"github.com/dsxg666/snakecoin/common"
	"github.com/dsxg666/snakecoin/console"
	"github.com/dsxg666/snakecoin/core"
	"github.com/dsxg666/snakecoin/db"
	"github.com/dsxg666/snakecoin/rlp"
	mptrie "github.com/dsxg666/snakecoin/trie"
	"github.com/fatih/color"
	"github.com/peterh/liner"
	"github.com/shopspring/decimal"
	"github.com/spf13/cobra"
)

func Console(cmd *cobra.Command) {
	acc, _ := cmd.Flags().GetString("account")
	if strings.Compare("", acc) == 0 {
		color.Red("Console command must specify account!")
		return
	}
	if !core.AccountIsExist(acc) {
		color.Red("The account you entered does net exist!")
		return
	}
	path := db.KeystorePath + "/" + acc
	priKeyBytes := common.ReadFile(path + "/private")
	pubKeyBytes := common.ReadFile(path + "/public")
	priKey := account.DecodePriKey(priKeyBytes)
	pubKey := account.DecodePubKey(pubKeyBytes)
	Interface(acc, priKey, pubKey)
}

func Interface(acc string, priKey *ecdsa.PrivateKey, pubKey *ecdsa.PublicKey) {
	Prompts(acc)
	line := console.GetLiner()
	defer line.Close()
	txDB := db.GetDB(db.TxPath)
	defer db.CloseDB(txDB)
	mptDB := db.GetDB(db.MPTirePath)
	defer db.CloseDB(mptDB)
	chainDB := db.GetDB(db.ChainPath)
	defer db.CloseDB(chainDB)
	for {
		if input, err := line.Prompt("> "); err != nil {
			fmt.Println()
			color.Blue("bye")
			fmt.Println()
			break
		} else {
			if strings.Compare("quit", input) == 0 {
				color.Blue("bye")
				fmt.Println()
				break
			}
			switch input {
			case "transaction":
				Transaction(acc, txDB, mptDB, priKey, pubKey, line)
			case "txpool":
				Txpool(acc, txDB, line)
			case "mine":
				Mine(acc, txDB, chainDB, mptDB)
			case "blockchain":
				Blockchain(acc, chainDB, line)
			case "mptrie":
				MPTrie(acc, chainDB, mptDB, line)
			case "balance":
				Balance(acc, mptDB)
			default:
				color.Yellow("Honey, we don't support your instruction yet.")
				fmt.Println()
			}
		}
	}
}

func Prompts(account string) {
	color.Blue("Welcome to the Snath console!")
	fmt.Printf("CurrentAccount: %s\n", account)
	fmt.Println("You can enter the following instruction to use snath:")
	fmt.Println("- [ transaction ] Conduct a transfer transaction")
	fmt.Println("- [ txpool ] Look at the transactions in the pool")
	fmt.Println("- [ mine ] Start mining")
	fmt.Println("- [ blockchain ] Look at block information")
	fmt.Println("- [ mptrie ] Look at state of different blocks")
	fmt.Println("- [ balance ] Check your account balance")
	fmt.Println()
	fmt.Println("To exit, press ctrl-d or input quit")
}

func Balance(acc string, mptDB *pebble.DB) {
	core.ShowAccountBalance(acc, mptDB)
	fmt.Println()
}

func Transaction(acc string, txDB, mptDB *pebble.DB, priKey *ecdsa.PrivateKey, pubKey *ecdsa.PublicKey, line *liner.State) {
	ok, loc := core.TxIsFull(txDB)
	if ok {
		color.Yellow("The current txpool is full!")
		fmt.Println()
		return
	}
	color.Blue("Welcome to Transaction Mode!")
	fmt.Println("To exit, press ctrl-d or input quit")
	mptBytes := db.Get([]byte{byte(99)}, mptDB)
	var e []interface{}
	err := rlp.DecodeBytes(mptBytes, &e)
	if err != nil {
		log.Panic("Failed to DecodeBytes:", err)
	}
	trie := mptrie.NewTrieWithDecodeData(e)
	accByte := common.Hex2Bytes(acc[2:])
	stateByte, _ := trie.Get(accByte)
	state := core.DeserializeState(stateByte)
	balance := state.Balance
	var finish bool
	for {
		fmt.Println("Which acc do you want do transfer skc to?")
		if to, err := line.Prompt("> "); err != nil {
			fmt.Println()
			fmt.Println()
			Prompts(acc)
			break
		} else {
			if strings.Compare("quit", to) == 0 {
				fmt.Println()
				Prompts(acc)
				break
			}
			if !core.AccountIsExist(to) {
				color.Red("The acc you entered does net exist!")
				fmt.Println()
				continue
			}
			if strings.Compare(acc, to) == 0 {
				color.Red("Do not transfer skc to yourself, It meaningless!")
				fmt.Println()
				continue
			}
			toB := common.Hex2Bytes(to[2:])
			fmt.Println()
			for {
				fmt.Println("How many skc do you want to transfer?")
				if input, err := line.Prompt("> "); err != nil {
					fmt.Println()
					fmt.Println()
					Prompts(acc)
					finish = true
					break
				} else {
					if strings.Compare("quit", input) == 0 {
						fmt.Println()
						Prompts(acc)
						finish = true
						break
					}
					if !common.StringIsAllNumber(input) {
						color.Red("The amount you entered is illegal!")
						fmt.Println()
						continue
					}
					amount, _ := strconv.ParseFloat(input, 64)
					if balance.Sub(decimal.NewFromFloat(amount)).LessThan(decimal.NewFromFloat(0)) {
						color.Red("Your balance is insufficient!")
						fmt.Println()
						continue
					}
					state.Balance = balance.Sub(decimal.NewFromFloat(amount))
					state.Nonce += 1
					trie.Update(accByte, state.Serialize())
					state2Bytes, _ := trie.Get(toB)
					state2 := core.DeserializeState(state2Bytes)
					state2.Balance = state2.Balance.Add(decimal.NewFromFloat(amount))
					trie.Update(toB, state2.Serialize())
					db.Set([]byte{byte(99)}, mptrie.Serialize(trie.Root), mptDB)
					tx := core.NewTransaction(decimal.NewFromFloat(amount), uint64(time.Now().Unix()),
						common.BytesToAddress(accByte), common.BytesToAddress(toB), account.EncodePubKey(pubKey))
					tx.Hash()
					tx.Sign(priKey)
					core.PushTxToPool(loc, tx, txDB)
					times := strings.Split(common.GetCurrentTime(), " ")
					color.Green("INFO [%s|%s] Successful transaction!", times[0], times[1])
					fmt.Println("TransactionHash: ", common.Bytes2Hex(tx.Hash()))
					fmt.Println()
					Prompts(acc)
					finish = true
					break
				}
			}
		}
		if finish {
			break
		}
	}
}

func Txpool(acc string, txDB *pebble.DB, line *liner.State) {
	_, loc := core.TxIsFull(txDB)
	if loc[0] == 0 {
		color.Yellow("There is no transaction in the Txpool.")
		fmt.Println()
		return
	}
	color.Blue("Welcome to the Txpool Mode!")
	fmt.Printf("There are now %d txs in the pool.\n", loc[0])
	fmt.Printf("You can enter 0 ~ %d to see the info of tx.\n", loc[0]-1)
	fmt.Println()
	fmt.Println("To exit, press ctrl-d or input quit")
	for {
		if input, err := line.Prompt("> "); err != nil {
			fmt.Println()
			fmt.Println()
			Prompts(acc)
			break
		} else {
			if strings.Compare("quit", input) == 0 {
				fmt.Println()
				Prompts(acc)
				break
			}
			i, _ := strconv.Atoi(input)
			if !common.StringIsAllNumber(input) || i < 0 || i > int(loc[0])-1 {
				color.Yellow("Your input is not valid!")
				continue
			}
			txBytes := db.Get([]byte{byte(i)}, txDB)
			tx := core.DeserializeTx(txBytes)
			fmt.Printf("TxHash: %x\n", tx.Hash())
			fmt.Printf("From: %s\n", tx.From.Hex())
			fmt.Printf("To: %s\n", tx.To.Hex())
			fmt.Println("Amount: ", tx.Amount, "skc")
			fmt.Println("State: pending")
			fmt.Println("Time: ", common.TimestampToTime(int64(tx.Time)))
			fmt.Println()
		}
	}
}

func Mine(acc string, txDB, chainDB, mptDB *pebble.DB) {
	_, loc := core.TxIsFull(txDB)
	if loc[0] == 0 {
		color.Yellow("There is no transaction in the Txpool.")
		fmt.Println()
		return
	}
	var txs []*core.Transaction
	for i := 0; i < int(loc[0]); i++ {
		txBytes := db.Get([]byte{byte(i)}, txDB)
		tx := core.DeserializeTx(txBytes)
		if !tx.Verity() {
			color.Red("Transaction verification failed!")
			fmt.Println()
			return
		}
		tx.State = 1
		txs = append(txs, tx)
		db.Set([]byte{byte(i)}, tx.Serialize(), txDB)
	}
	accBytes := common.Hex2Bytes(acc[2:])
	core.NewBlock(common.BytesToAddress(accBytes), chainDB, mptDB, txs)
	mptBytes := db.Get([]byte{byte(99)}, mptDB)
	var e []interface{}
	err := rlp.DecodeBytes(mptBytes, &e)
	if err != nil {
		log.Panic("Failed to DecodeBytes: ", err)
	}
	trie := mptrie.NewTrieWithDecodeData(e)
	stateBytes, _ := trie.Get(accBytes)
	state := core.DeserializeState(stateBytes)
	state.Balance = state.Balance.Add(decimal.NewFromFloat(10))
	trie.Update(accBytes, state.Serialize())
	db.Set([]byte{byte(99)}, mptrie.Serialize(trie.Root), mptDB)
	times := strings.Split(common.GetCurrentTime(), " ")
	fmt.Println()
	color.Green("INFO [%s|%s] A block was successfully mined!", times[0], times[1])
	fmt.Println("Account", acc, "will be awarded 10 skc.")
	fmt.Println()
}

func Blockchain(acc string, chainDB *pebble.DB, line *liner.State) {
	lastBlockHashBytes := db.Get([]byte{byte(99)}, chainDB)
	lastBlockBytes := db.Get(lastBlockHashBytes, chainDB)
	lastBlock := core.DeserializeBlock(lastBlockBytes)
	num := lastBlock.Header.Number
	color.Blue("Welcome to the Blockchain Mode!")
	fmt.Printf("There are now %d blocks in blockchain.\n", num+1)
	fmt.Printf("You can enter 0 ~ %d to see the info of block.\n", num)
	fmt.Println()
	fmt.Println("To exit, press ctrl-d or input quit")
	for {
		if input, err := line.Prompt("> "); err != nil {
			fmt.Println()
			fmt.Println()
			Prompts(acc)
			break
		} else {
			if strings.Compare("quit", input) == 0 {
				fmt.Println()
				Prompts(acc)
				break
			}
			in, _ := strconv.Atoi(input)
			if !common.StringIsAllNumber(input) || in < 0 || in > int(num) {
				color.Yellow("Your input is not valid!")
				continue
			}
			var temp *core.Block
			temp = lastBlock
			for i := 0; i < int(num)-in; i++ {
				prevBlockHash := db.Get(temp.Header.PrevBlockHash.Bytes(), chainDB)
				temp = core.DeserializeBlock(prevBlockHash)
			}
			color.Green("BlockHeaderInformation")
			fmt.Printf("Number: %d\n", temp.Header.Number)
			fmt.Printf("Nonce: %d\n", temp.Header.Nonce)
			fmt.Printf("Difficulty: %d\n", temp.Header.Difficulty)
			fmt.Printf("Time: %s\n", common.TimestampToTime(int64(temp.Header.Time)))
			fmt.Printf("Coinbase: %s\n", temp.Header.Coinbase.Hex())
			fmt.Printf("BlockHash: %x\n", temp.Header.BlockHash.Bytes())
			fmt.Printf("PrevBlockHash: %x\n", temp.Header.PrevBlockHash.Bytes())
			fmt.Printf("StateTreeRoot: %x\n", temp.Header.StateTreeRoot.Bytes())
			fmt.Printf("MerkleTreeRoot: %x\n", temp.Header.MerkleTreeRoot.Bytes())
			fmt.Println()
			if temp.Header.Number == 0 {
				continue
			}
			color.Green("BlockBodyInformation")
			for i := 0; i < len(temp.Body.Txs); i++ {
				color.Green("Tx%d", i)
				fmt.Printf("TxHash: %x\n", temp.Body.Txs[i].Hash())
				fmt.Printf("From: %s\n", temp.Body.Txs[i].From.Hex())
				fmt.Printf("To: %s\n", temp.Body.Txs[i].To.Hex())
				fmt.Println("Amount: ", temp.Body.Txs[i].Amount, "skc")
				fmt.Println("State: stored")
				fmt.Println("Time: ", common.TimestampToTime(int64(temp.Body.Txs[i].Time)))
				fmt.Println()
			}
		}
	}
}

func MPTrie(acc string, chainDB, mptDB *pebble.DB, line *liner.State) {
	lastBlockHashBytes := db.Get([]byte{byte(99)}, chainDB)
	lastBlockBytes := db.Get(lastBlockHashBytes, chainDB)
	lastBlock := core.DeserializeBlock(lastBlockBytes)
	num := lastBlock.Header.Number
	color.Blue("Welcome to the MPTrie Mode!")
	fmt.Printf("You can enter 1 ~ %d to see state of different blocks.\n", num)
	fmt.Println()
	fmt.Println("To exit, press ctrl-d or input quit")
	for {
		if input, err := line.Prompt("> "); err != nil {
			fmt.Println()
			fmt.Println()
			Prompts(acc)
			break
		} else {
			if strings.Compare("quit", input) == 0 {
				fmt.Println()
				Prompts(acc)
				break
			}
			in, _ := strconv.Atoi(input)
			if !common.StringIsAllNumber(input) || in < 1 || in > int(num) {
				color.Yellow("Your input is not valid!")
				continue
			}
			mptBytes := db.Get([]byte{byte(in)}, mptDB)
			var e []interface{}
			err := rlp.DecodeBytes(mptBytes, &e)
			if err != nil {
				log.Panic("Failed to DecodeBytes: ", err)
			}
			trie := mptrie.NewTrieWithDecodeData(e)
			files, _ := filepath.Glob(db.KeystorePath + "/*")
			for i := 0; i < len(files); i++ {
				stateBytes, _ := trie.Get(common.Hex2Bytes(files[i][16:]))
				if stateBytes == nil {
					continue
				}
				state := core.DeserializeState(stateBytes)
				fmt.Println(files[i][14:], "{nonce:", state.Nonce, ";balance:", state.Balance, "}")
			}
		}
	}
}
