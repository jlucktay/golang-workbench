package main

import (
	"bytes"
	"errors"
	"fmt"
	"strings"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

const (
	envPrefix = "gbhc"
)

func gatherConfig(arguments []string) error {
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	viper.SetDefault("api-key", "")
	viper.SetConfigName(envPrefix)
	viper.SetConfigType("json")
	viper.AddConfigPath(".")

	if errRead := viper.ReadInConfig(); errRead != nil {
		if errWrite := viper.SafeWriteConfig(); errWrite != nil {
			return errWrite
		}
	}

	viper.SetEnvPrefix(envPrefix)
	viper.AutomaticEnv()

	fs := pflag.NewFlagSet(envPrefix, pflag.ContinueOnError)

	fs.String("api-key", "", "Your Giant Bomb API key")
	fs.Bool("help", false, "Help with usage")

	if errParse := fs.Parse(arguments); errParse != nil {
		return errParse
	}

	if errBind := viper.BindPFlags(fs); errBind != nil {
		return errBind
	}

	if viper.GetBool("help") {
		buf := bytes.NewBufferString(fmt.Sprintf("Usage of '%s':\n", arguments[0]))
		fs.SetOutput(buf)
		fs.PrintDefaults()

		return errors.New(buf.String())
	}

	if viper.Get("api-key") == "" {
		return errors.New(
			`No API key was provided. You can solve this by doing one of the following:
  - Invoke the application with the '--api-key="<value>"' flag
  - Set a GBHC_API_KEY environment variable
  - Set the 'api-key' string in your 'gbhc.json' config file`)
	}

	return nil
}
