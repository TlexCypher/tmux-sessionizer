package cmd

import (
	"context"
	"errors"
	"os"
	"os/signal"

	"github.com/TlexCypher/my-tmux-sessionizer/handler"
	"github.com/urfave/cli/v3"
)

const (
	ExitCodeOK        = 0
	ExitCodeError int = iota
)

var (
	ErrNoSuchCmd = errors.New("no such command")
)

func Core() int {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	cmd := newCmd()
	if err := cmd.Run(ctx, os.Args); err != nil {
		return ExitCodeError
	}
	return ExitCodeOK
}

func newCmd() *cli.Command {
	return &cli.Command{
		Name:   "tmux-sessionizer",
		Usage:  "tmux session manager",
		Action: run,
	}
}

func run(ctx context.Context, cmd *cli.Command) error {
	sh := handler.NewSessionHandler()
	return runWithHandler(sh, ctx, cmd)
}

func runWithHandler(h handler.ISessionHandler, ctx context.Context, cmd *cli.Command) error {
	args := cmd.Args().Slice()
	if len(args) > 0 && args[0] == "list" {
		return h.GrabExistingSession()
	} else if len(args) > 0 {
		return ErrNoSuchCmd
	} else {
		return h.NewSession()
	}
}
