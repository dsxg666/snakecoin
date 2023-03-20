package main

import (
	"os"
	"strings"

	"github.com/dsxg666/snakecoin/core"
	"github.com/dsxg666/snakecoin/db"
	"github.com/dsxg666/snakecoin/logs"
	"github.com/dsxg666/snakecoin/util"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

// initCmd 代表 init 命令
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize a new genesis block",
	Long: `
The init command initializes a new genesis block and
initializes the folder structure of the data. It
expects the file path as argument.`,
	Run: func(cmd *cobra.Command, args []string) {
		InitDataDir()
		logger := logs.InitLogger()
		chainDB := db.GetDB(db.ChainDataPath)
		defer db.CloseDB(chainDB)
		core.NewBlockchain(chainDB)
		s := strings.Split(util.CurrentTimeFormant(), " ")
		color.Green("INFO [%s|%s] Congratulations, the Genesis block was successfully initialized!", s[0], s[1])
		logger.Infof("INFO [%s|%s] Congratulations, the Genesis block was successfully initialized!", s[0], s[1])
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}

func InitDataDir() {
	_, err := os.Stat("./data")
	if os.IsNotExist(err) {
		_ = os.Mkdir("./data", 0666)
	}
	_ = os.Mkdir(db.AccountDataPath, 0666)
	_ = os.Mkdir(db.ChainDataPath, 0666)
	_ = os.Mkdir(db.KeystoreDataPath, 0666)
	_ = os.Mkdir(db.LogDataPath, 0666)
	_ = os.Mkdir(db.TxDataPath, 0666)
}
