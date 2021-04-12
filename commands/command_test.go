package commands

import (
	"context"
	"os"
	"testing"
)

func TestCommand(t *testing.T) {
	ctx := context.Background()

	Execute(ctx, os.Args[1:])
}
