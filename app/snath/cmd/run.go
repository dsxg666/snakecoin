package cmd

import (
	"github.com/dsxg666/snakecoin/app/snath/logic"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(runCmd)
	runCmd.Flags().StringP("rpc.port", "r", "8545", "password used to unlock the account")
}

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Start a Snakecoin node",
	Long:  `Start and connect to the Snakecoin main network.`,
	Run: func(cmd *cobra.Command, args []string) {
		logic.Run(cmd)
	},
}
