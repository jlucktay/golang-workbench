package main

import (
	"io"
	"log/slog"
	"os"

	"github.com/lmittmann/tint"
	"golang.org/x/term"
)

func main() {
	defer slog.Info("outer deferred goodbye") // This log won't show; it's being sent to a logger that stops existing.

	status := run(os.Stdout)

	slog.Info("goodbye before os.Exit call", slog.String("emoji", "👋"))

	os.Exit(status)
}

func run(stdout io.Writer) int {
	defer slog.Info("inner deferred goodbye", slog.String("emoji", "⏯️"))

	var handler slog.Handler

	stdoutFile, isFile := stdout.(*os.File)
	if isFile && term.IsTerminal(int(stdoutFile.Fd())) {
		handler = tint.NewHandler(stdout, nil)
	} else {
		handler = slog.NewJSONHandler(stdout, nil)
	}

	slogger := slog.New(handler)

	slogger.Info("direct call to new structured logger", slog.String("emoji", "👀"))

	defer slogger.Info("deferred call made directly to new structured logger", slog.String("emoji", "🐌"))

	slog.SetDefault(slogger)

	slog.Info("hello via slog package default logger", slog.String("emoji", "🙂"))

	return 0
}
