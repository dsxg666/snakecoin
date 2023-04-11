package cmd

import (
	"github.com/dsxg666/snakecoin/app/snath/logic"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Version information",
	Long:  `Print version numbers.`,
	Run: func(cmd *cobra.Command, args []string) {
		logic.Version()
	},
}
