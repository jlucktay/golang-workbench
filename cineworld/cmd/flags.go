package cmd

import (
	"github.com/spf13/pflag"
)

//nolint:gochecknoglobals // Flags to pass parsed values forward to command logic.
var (
	flagCinemaID                   string
	flagFutureDays, flagNumberDays int
	flagInclude3D                  bool
)

// rootPersistentFlags returns a [pflag.FlagSet] to be added to the root command.
func rootPersistentFlags() *pflag.FlagSet {
	pfs := &pflag.FlagSet{}

	pfs.StringVarP(&flagCinemaID, "cinema-id", "c", "073", "ID of cinema to pull screenings for")
	pfs.IntVarP(&flagFutureDays, "future-days", "f", 0, "start listing from this many days into the future")
	pfs.BoolVarP(&flagInclude3D, "include-3d", "3", false, "include screenings in 3D")
	pfs.IntVarP(&flagNumberDays, "number-days", "n", 1, "retrieve screenings for this many days")

	return pfs
}
