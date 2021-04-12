package commands

import (
	"context"
	"flag"

	"github.com/google/subcommands"
)

func newCmd(ctx context.Context) (topLevelFS *flag.FlagSet, cmd *subcommands.Commander) {
	topLevelFS = flag.NewFlagSet("top-level", flag.ContinueOnError)

	cmd = subcommands.NewCommander(topLevelFS, "root")
	return
}

// Execute k
func Execute(
	ctx context.Context,
	args []string,
) (
	exitStatus subcommands.ExitStatus,
	err error,
) {
	fs, cmd := newCmd(ctx)

	err = fs.Parse(args)
	if err != nil {
		exitStatus = subcommands.ExitUsageError
		return
	}

	exitStatus = cmd.Execute(ctx)
	return
}
