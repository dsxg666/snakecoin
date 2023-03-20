package console

import (
	"strings"

	"github.com/peterh/liner"
)

var (
	names = []string{"balance", "blockchain", "mining", "txpool", "transaction"}
)

func Line() *liner.State {
	line := liner.NewLiner()
	defer line.Close()
	line.SetCtrlCAborts(false)
	line.SetCompleter(func(line string) (c []string) {
		for _, n := range names {
			if strings.HasPrefix(n, strings.ToLower(line)) {
				c = append(c, n)
			}
		}
		return
	})
	return line
}
