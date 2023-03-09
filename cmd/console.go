package cmd

import (
	"github.com/dsxg666/snakecoin/cmd/console"
	"github.com/dsxg666/snakecoin/db"
	"github.com/dsxg666/snakecoin/logs"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"strings"
)

// consoleCmd represents the console command
var consoleCmd = &cobra.Command{
	Use:   "console",
	Short: "Start an interactive environment",
	Long: `
The snakecoin console is an interactive shell for
blockchain runtime environment which can interact
with the blockchain.`,
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
				console.Start(logger, account, txDB, accountDB, chainDB)
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
