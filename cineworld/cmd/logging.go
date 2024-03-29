package cmd

import (
	"io"
	"log/slog"
	"os"

	"github.com/lmittmann/tint"
	"golang.org/x/term"
)

// setUpLogging will detect whether stderr is a terminal.
// If so it configures human-readable colourful logs, or JSON logs if not.
// Per [these CLI guidelines], logging is assumed to be sent to stderr, as implied by this argument's name.
//
// [these CLI guidelines]: https://clig.dev/#the-basics
func setUpLogging(stderr io.Writer) *slog.Logger {
	var handler slog.Handler

	stderrFile, isFile := stderr.(*os.File)
	if isFile && term.IsTerminal(int(stderrFile.Fd())) {
		handler = tint.NewHandler(stderr, nil)
	} else {
		handler = slog.NewJSONHandler(stderr, nil)
	}

	return slog.New(handler)
}
