package cmd

import (
	"github.com/dsxg666/snakecoin/app/snath/logic"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(runCmd)
}

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Start a Snakecoin node",
	Long:  `Start and connect to the Snakecoin main network.`,
	Run: func(cmd *cobra.Command, args []string) {
		logic.Run()
	},
}
