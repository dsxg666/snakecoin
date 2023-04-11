package cmd

import (
	"github.com/dsxg666/snakecoin/app/snath/logic"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(accountCmd)
	accountCmd.AddCommand(newCmd)
	accountCmd.AddCommand(listCmd)
}

var accountCmd = &cobra.Command{
	Use:   "account",
	Short: "Manage accounts",
	Long: `
Account command can create a new account and list all
your existing accounts.`,
}

var newCmd = &cobra.Command{
	Use:   "new",
	Short: "New account",
	Long: `
New Command can create a new account for you.`,
	Run: func(cmd *cobra.Command, args []string) {
		logic.New()
	},
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all accounts",
	Long: `
List command can list all your accounts.`,
	Run: func(cmd *cobra.Command, args []string) {
		logic.List()
	},
}
