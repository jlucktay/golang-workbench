package cmd

import (
	"errors"
	"fmt"
	"io"
	"log/slog"
	"path/filepath"
	"strings"
	"time"

	"github.com/carlmjohnson/versioninfo"
	"github.com/spf13/cobra"
)

// Exit status codes returned by the CLI.
const (
	// ExitSuccess when everything goes to plan.
	ExitSuccess = iota

	// ExitUnknown if the cause of the error is not defined.
	ExitUnknown

	// ExitNoCommandName if Execute is passed a zero-length string slice, without a command name as the first element.
	ExitNoCommandName

	// ExitParsingArguments if parsing arguments goes awry.
	ExitParsingArguments
)

var ErrUnknownArguments = errors.New("unknown arguments passed in")

// Execute creates a new root command, sets flags appropriately, connects the two given [io.Writer] interfaces to
// stdout and stderr, and passes the given string slice forward to be parsed for flags and args.
// This is called by main.main(). It only needs to happen once to the root command.
func Execute(stdout, stderr io.Writer, args []string) int {
	if len(args) < 1 {
		return ExitNoCommandName
	}

	cmdName := filepath.Base(args[0])

	slogger := setUpLogging(stderr)

	version := fmt.Sprintf("%s built on %s from git SHA %s",
		versioninfo.Version, versioninfo.LastCommit.UTC().Format(time.RFC3339), versioninfo.Revision)

	if versioninfo.DirtyBuild {
		version += " (dirty)"
	}

	// rootCmd represents the base command when called without any subcommands
	rootCmd := &cobra.Command{
		Use: cmdName,

		Short: "Fetch, parse, and display screenings at a given Cineworld cinema",
		Long: `Fetch, parse, and display screenings at a given Cineworld cinema.
Shows screenings for the rest of today, and excludes those in 3D, by default.
Flags can be set to show screenings for additional days into the future, include 3D screenings, and start from a future
day after today.
The cinema to show screenings for can also be set, by its ID.`,

		Example: `  ` + cmdName + ` --number-days 7
  ` + cmdName + ` -3 -f=1`,

		Version: version,

		PersistentPreRun: func(_ *cobra.Command, _ []string) { setUpLogging(stderr) },
		RunE:             root,
	}

	// Wire in the arguments passed to this func.
	rootCmd.SetOut(stdout)
	rootCmd.SetErr(stderr)
	rootCmd.SetArgs(args[1:])

	// Add persistent flags to the root command, and then execute it.
	rootCmd.PersistentFlags().AddFlagSet(rootPersistentFlags())

	if err := rootCmd.Execute(); err != nil {
		defer slogger.Error("executing root command",
			slog.Any("err", err))

		switch {
		case errors.Is(err, ErrUnknownArguments):
			return ExitParsingArguments
		default:
			return ExitUnknown
		}
	}

	return ExitSuccess
}

// root is the core command logic that is wrapped and run by [Execute].
func root(cmd *cobra.Command, args []string) error {
	// Check for any remaining unparsed arguments.
	if len(args) > 0 {
		return fmt.Errorf("%w: '%s'", ErrUnknownArguments, strings.Join(args, "', '"))
	}

	// Display help text when passed no options.
	// Cf. https://clig.dev/#help
	if cmd.Flags().HasFlags() && cmd.Flags().NFlag() == 0 {
		fmt.Fprintf(cmd.OutOrStdout(), "%s\n\n%s", cmd.Long, cmd.UsageString())
		return nil
	}

	return nil
}
