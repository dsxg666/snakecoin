package console

import (
	"strings"

	"github.com/peterh/liner"
)

var (
	instructions = []string{"balance", "blockchain", "mine", "txpool", "transaction", "quit", "mptrie", "lookupInfoByHash"}
)

func GetLiner() *liner.State {
	line := liner.NewLiner()
	line.SetCtrlCAborts(false)
	line.SetCompleter(func(line string) (c []string) {
		for _, instruction := range instructions {
			if strings.HasPrefix(instruction, strings.ToLower(line)) {
				c = append(c, instruction)
			}
		}
		return
	})
	return line
}
