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

	slog.Info("goodbye", slog.String("emoji", "👋"))

	os.Exit(status)
}

func run(stdout io.Writer) int {
	defer slog.Info("inner deferred goodbye", slog.String("emoji", "⏯️"))

	var handler slog.Handler

	stdoutFile, isFile := stdout.(*os.File)
	if isFile && term.IsTerminal(int(stdoutFile.Fd())) {
		handler = tint.NewHandler(stdout, &tint.Options{NoColor: false})
	} else {
		handler = slog.NewJSONHandler(stdout, &slog.HandlerOptions{})
	}

	slog.SetDefault(slog.New(handler))

	slog.Info("hello", slog.String("emoji", "🙂"))

	return 0
}
