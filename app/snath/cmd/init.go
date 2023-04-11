package cmd

import (
	"github.com/dsxg666/snakecoin/app/snath/logic"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(initCmd)
}

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Complete the initialization work of the blockchain",
	Long: `
The init command initializes the directory structure 
of the data store(generated data directory for you in
the current directory) and initializes a genesis block.`,
	Run: func(cmd *cobra.Command, args []string) {
		logic.Init()
	},
}
