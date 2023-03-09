package cmd

import (
	"fmt"
	"os"
	"runtime"
	"strings"

	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Version information",

	Long: `Print version numbers`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("SnakeCoin")
		fmt.Println("Version: 1.0.0-stable")
		fmt.Printf("Architecture: %s\n", runtime.GOARCH)
		fmt.Printf("Go Version: %s\n", runtime.Version())
		fmt.Printf("Operating System: %s\n", runtime.GOOS)
		fmt.Printf("GOPATH=%s\n", os.Getenv("GOPATH"))
		fmt.Printf("GOROOT=%s\n", strings.ReplaceAll(os.Getenv("GOROOT"), "/", "\\"))
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
