package cmd

import (
	"errors"
	"fmt"
	"io"
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
)

//nolint:gochecknoglobals // Flags to pass parsed values forward to command logic.
var (
	flagCinemaID                   string
	flagFutureDays, flagNumberDays int
	flagInclude3D                  bool
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

		RunE: root,
	}

	// Wire in the arguments passed to this func.
	rootCmd.SetOut(stdout)
	rootCmd.SetErr(stderr)
	rootCmd.SetArgs(args[1:])

	// Set up flags.
	rootCmd.PersistentFlags().StringVarP(&flagCinemaID, "cinema-id", "c", "073", "ID of cinema to pull screenings for")
	rootCmd.PersistentFlags().IntVarP(&flagFutureDays, "future-days", "f", 0,
		"start listing from this many days into the future")
	rootCmd.PersistentFlags().BoolVarP(&flagInclude3D, "include-3d", "3", false, "include screenings in 3D")
	rootCmd.PersistentFlags().IntVarP(&flagNumberDays, "number-days", "n", 1, "retrieve screenings for this many days")

	if err := rootCmd.Execute(); err != nil {
		return ExitUnknown
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
