package main

import (
	"context"
	"os"

	botCmd "github.com/myyang/imbot/commands"
	botLog "github.com/myyang/imbot/log"
)

func main() {
	ctx := context.Background()
	ctx = context.WithValue(ctx, botLog.CtxLogger, &botLog.DebugLogger{})

	botCmd.Execute(ctx, os.Args[1:])
}
