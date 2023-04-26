package main

import (
	"context"
	"flag"
	"os"

	"github.com/dsxg666/cli"
	"github.com/dsxg666/snakecoin/app/svm/cmd"
)

func main() {
	cli.Register(cli.HelpCommand(), "")
	cli.Register(cli.FlagsCommand(), "")
	cli.Register(cli.CommandsCommand(), "")
	cli.Register(&cmd.CompileCmd{}, "")
	cli.Register(&cmd.DumpCmd{}, "")
	cli.Register(&cmd.ExecuteCmd{}, "")
	cli.Register(&cmd.RunCmd{}, "")
	cli.Register(&cmd.VersionCmd{}, "")

	flag.Parse()
	ctx := context.Background()
	os.Exit(int(cli.Execute(ctx)))
}
