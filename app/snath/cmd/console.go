package cmd

import (
	"github.com/dsxg666/snakecoin/app/snath/logic"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(consoleCmd)
	consoleCmd.Flags().StringP("account", "a", "", "an account for mining and transferring skc")
	consoleCmd.Flags().StringP("password", "p", "", "password used to unlock the account")
}

var consoleCmd = &cobra.Command{
	Use:   "console",
	Short: "Start an interactive environment",
	Long: `
The snath is an interactive shell for blockchain
runtime environment witch user can interactive
with snakecoin blockchain.`,
	Run: func(cmd *cobra.Command, args []string) {
		logic.Console(cmd)
	},
}
