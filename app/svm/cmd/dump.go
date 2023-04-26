package cmd

import (
	"context"
	"flag"
	"fmt"
	"io/ioutil"

	"github.com/dsxg666/cli"
	"github.com/dsxg666/snakecoin/vm/compiler"
	"github.com/dsxg666/snakecoin/vm/lexer"
)

type DumpCmd struct{}

func (*DumpCmd) Name() string     { return "dump" }
func (*DumpCmd) Synopsis() string { return "Show the lexed output of the given program" }
func (*DumpCmd) Usage() string {
	return `dump :
  Demonstrate how our lexer performed by dumping the given input file, as a
  stream of tokens.
`
}

func (p *DumpCmd) SetFlags(f *flag.FlagSet) {}

func (p *DumpCmd) Execute(_ context.Context, f *flag.FlagSet, _ ...interface{}) cli.ExitStatus {
	for _, file := range f.Args() {

		// Read the file.
		input, err := ioutil.ReadFile(file)
		if err != nil {
			fmt.Printf("Error reading %s - %s\n", file, err.Error())
			return cli.ExitFailure
		}

		// Lex it
		l := lexer.New(string(input))

		// Dump it
		e := compiler.New(l)
		e.Dump()
	}
	return cli.ExitSuccess
}
