package cmd

import (
	"context"
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/dsxg666/cli"
)

var out io.Writer = os.Stdout

var (
	version = "SVM-Version: v1.0.0-stable"
)

type VersionCmd struct{}

func (*VersionCmd) Name() string     { return "version" }
func (*VersionCmd) Synopsis() string { return "Version information" }
func (*VersionCmd) Usage() string {
	return `Print version numbers.`
}

func (p *VersionCmd) SetFlags(f *flag.FlagSet) {}

func showVersion() {
	fmt.Fprintf(out, "%s\n", version)
}

func (p *VersionCmd) Execute(_ context.Context, f *flag.FlagSet, _ ...interface{}) cli.ExitStatus {
	showVersion()
	return cli.ExitSuccess
}
