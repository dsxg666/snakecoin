package cmd

import (
	"context"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/dsxg666/cli"
	"github.com/dsxg666/snakecoin/vm/compiler"
	"github.com/dsxg666/snakecoin/vm/lexer"
)

type CompileCmd struct{}

func (*CompileCmd) Name() string     { return "compile" }
func (*CompileCmd) Synopsis() string { return "Compile program" }
func (*CompileCmd) Usage() string {
	return `Compile the given input file to a series of bytecodes.`
}

func (p *CompileCmd) SetFlags(*flag.FlagSet) {}

func (p *CompileCmd) Execute(_ context.Context, f *flag.FlagSet, _ ...interface{}) cli.ExitStatus {
	for _, file := range f.Args() {
		input, err := os.ReadFile(file)
		if err != nil {
			fmt.Printf("Error reading %s - %s\n", file, err.Error())
			return cli.ExitFailure
		}

		l := lexer.New(string(input))
		c := compiler.New(l)
		c.Compile()
		name := strings.TrimSuffix(file, filepath.Ext(file))
		c.Write(name + ".raw")
	}
	return cli.ExitSuccess
}
