package console

import (
	"fmt"
	"github.com/fatih/color"
)

// FirstMeet 第一次进入控制台主页的提示符
func FirstMeet(account string) {
	color.Blue("Welcome to the Snath console!")
	fmt.Println()
	fmt.Printf("CurrentAccountAddress: %s\n", account)
	fmt.Println("You can enter the following instruction to use blockchain:")
	fmt.Println("- [ transaction ] Conduct a transfer transaction")
	fmt.Println("- [ txpool ] Look at the transactions in the pool")
	fmt.Println("- [ mining ] Enter the mining program")
	fmt.Println("- [ blockchain ] See block information")
	fmt.Println("- [ balance ] Check your account balance")
	fmt.Println()
	fmt.Println("To exit, press ctrl-d or input quit")
}

// MeetAgain 第一次之后进行控制台主页的提示符
func MeetAgain(account string) {
	color.Blue("Welcome back to the Snath console!")
	fmt.Println()
	fmt.Printf("CurrentAccountAddress: %s\n", account)
	fmt.Println("You can enter the following instruction to use blockchain:")
	fmt.Println("- [ transaction ] Conduct a transfer transaction")
	fmt.Println("- [ txpool ] You can view the situation in the txpool")
	fmt.Println("- [ mining ] Enter the mining program")
	fmt.Println("- [ blockchain ] See block information")
	fmt.Println("- [ balance ] Check your account balance")
	fmt.Println()
	fmt.Println("To exit, press ctrl-d or input quit")
}
