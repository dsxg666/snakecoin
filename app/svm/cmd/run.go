package cmd

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/dsxg666/cli"
	"github.com/dsxg666/snakecoin/vm/compiler"
	"github.com/dsxg666/snakecoin/vm/cpu"
	"github.com/dsxg666/snakecoin/vm/lexer"
)

type RunCmd struct{}

func (*RunCmd) Name() string     { return "run" }
func (*RunCmd) Synopsis() string { return "Run the given source program" }
func (*RunCmd) Usage() string {
	return `run :
  The run sub-command compiles the given source program, and then executes
  it immediately.
`
}

func (p *RunCmd) SetFlags(f *flag.FlagSet) {}

func (p *RunCmd) Execute(_ context.Context, f *flag.FlagSet, _ ...interface{}) cli.ExitStatus {
	// For each file on the command-line both compile and execute it.
	for _, file := range f.Args() {
		fmt.Printf("Parsing file: %s\n", file)

		// Read the file.
		input, err := os.ReadFile(file)
		if err != nil {
			fmt.Printf("Error reading %s - %s\n", file, err.Error())
			return cli.ExitFailure
		}

		// Lex it
		l := lexer.New(string(input))

		// Compile it.
		e := compiler.New(l)
		e.Compile()

		// Now create a machine to run the compiled program in
		c := cpu.NewCPU()

		// Load the program
		c.LoadBytes(e.Output())

		// Run the machine
		err = c.Run()
		if err != nil {
			fmt.Printf("Error running file: %s\n", err)
			return cli.ExitFailure
		}
	}
	return cli.ExitSuccess
}
