package logic

import (
	"fmt"
	"runtime"
)

func Version() {
	fmt.Println("Snath")
	fmt.Println("Version: 1.4.0-stable")
	fmt.Printf("Architecture: %s\n", runtime.GOARCH)
	fmt.Printf("Go Version: %s\n", runtime.Version())
	fmt.Printf("Operating System: %s\n", runtime.GOOS)
	fmt.Println()
}
