package commands

import (
	"context"
	"os"
	"testing"

	botLog "github.com/myyang/imbot/log"
)

func TestCommand(t *testing.T) {
	ctx := context.Background()
	ctx = context.WithValue(ctx, botLog.CtxLogger, &botLog.DebugLogger{})

	Execute(ctx, os.Args[1:])
}
