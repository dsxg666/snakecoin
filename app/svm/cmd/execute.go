package cmd

import (
	"context"
	"flag"
	"fmt"

	"github.com/dsxg666/cli"
	"github.com/dsxg666/snakecoin/vm/cpu"
)

type ExecuteCmd struct{}

func (*ExecuteCmd) Name() string     { return "execute" }
func (*ExecuteCmd) Synopsis() string { return "Executed a compiled program" }
func (*ExecuteCmd) Usage() string {
	return `Execute the bytecodes contained in the given input file.`
}

func (p *ExecuteCmd) SetFlags(*flag.FlagSet) {
}

func (p *ExecuteCmd) Execute(_ context.Context, f *flag.FlagSet, _ ...interface{}) cli.ExitStatus {
	for _, file := range f.Args() {

		c := cpu.NewCPU()

		err := c.LoadFile(file)
		if err != nil {
			fmt.Printf("Error loading file: %s\n", err)
		}

		err = c.Run()
		if err != nil {
			fmt.Printf("Error running file: %s\n", err)
			return cli.ExitFailure
		}
	}
	return cli.ExitSuccess
}
