package main

import (
	"context"
	"os"

	botCmd "github.com/myyang/imbot/commands"
)

func main() {
	ctx := context.Background()

	botCmd.Execute(ctx, os.Args[1:])
}
