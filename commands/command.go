package commands

import (
	"context"
	"flag"
	"fmt"
	"strings"

	"github.com/google/subcommands"

	botLog "github.com/myyang/imbot/log"
)

var (
	defaultCommandsCmd = &commandsCmd{}

	customCommands = []subcommands.Command{
		defaultCommandsCmd,
	}
)

type commandsCmd struct {
	usage string
}

func (c *commandsCmd) Name() string             { return "commands" }
func (c *commandsCmd) Synopsis() string         { return "show available commands" }
func (c *commandsCmd) SetFlags(f *flag.FlagSet) {}
func (c *commandsCmd) Usage() string            { return c.Name() + ":\n    " + c.Synopsis() }

func (c *commandsCmd) Execute(ctx context.Context, f *flag.FlagSet, args ...interface{}) subcommands.ExitStatus {
	// add commands here for revealing
	logger := ctx.Value(botLog.CtxLogger).(botLog.Logger)
	if c.usage != "" {
		logger.Infof("%s\n", c.usage)
		return subcommands.ExitSuccess
	}

	b := strings.Builder{}
	b.WriteString("available commands:\n\n")
	for _, cc := range customCommands {
		b.WriteString(fmt.Sprintf("%-30v%s\n", cc.Name(), cc.Synopsis()))
	}

	b.WriteString("\nOr use `CMD help` for more details\n")

	c.usage = b.String()
	logger.Infof("%s", c.usage)
	return subcommands.ExitSuccess
}

func newCmd(ctx context.Context) (topLevelFS *flag.FlagSet, cmd *subcommands.Commander) {
	topLevelFS = flag.NewFlagSet("top-level", flag.ContinueOnError)

	cmd = subcommands.NewCommander(topLevelFS, "root")

	// custom commands
	for _, cc := range customCommands {
		cmd.Register(cc, "op")
	}

	return
}

// Execute is entrypoint for imbot commands.
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
		defaultCommandsCmd.Execute(ctx, fs)
		exitStatus = subcommands.ExitUsageError
		return
	}

	exitStatus = cmd.Execute(ctx)
	if exitStatus == subcommands.ExitUsageError {
		defaultCommandsCmd.Execute(ctx, fs)
	}
	return
}
