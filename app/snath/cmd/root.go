package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "snath",
	Short: "The snakecoin command line interface",
	Long: `
The snakecoin is a simple blockchain technology implementation
and the snath is a command line interface to snakecoin.`,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
