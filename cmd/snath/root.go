package main

import (
	"os"

	"github.com/spf13/cobra"
)

// rootCmd 根命令
var rootCmd = &cobra.Command{
	Use:   "snath",
	Short: "The snakecoin command line interface",
	Long: `
The snath is the benchmark ethereum, and although
there are few features implemented so far,I will keep
updating iteratively.`,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
