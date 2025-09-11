package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/DataDog/datadog-api-client-go/v2/api/datadog"
	"github.com/DataDog/datadog-api-client-go/v2/api/datadogV1"
)

const (
	datadogAPIKey         = "DD_CLIENT_API_KEY"
	datadogApplicationKey = "DD_CLIENT_APP_KEY"
)

var targetSLOName string

func init() {
	const usage = "the name of the SLO to emit (required)"

	flag.StringVar(&targetSLOName, "slo-name", "", usage)
	flag.StringVar(&targetSLOName, "s", "", usage+" (shorthand)")
}

func main() {
	flag.Parse()

	if targetSLOName == "" {
		fmt.Fprintf(os.Stderr, "The name of the SLO to emit must be supplied with the -s/--slo-name flag.\n\n")

		flag.Usage()

		fmt.Fprintf(os.Stderr, `
Example:
	%[1]s -s='SLO name'
	%[1]s --slo-name="SLO name"
`, os.Args[0])

		return
	}

	requiredEnvVars := map[string]string{
		datadogAPIKey:         "API",
		datadogApplicationKey: "Application",
	}

	for evKey, evType := range requiredEnvVars {
		if evValue, set := os.LookupEnv(evKey); !set || evValue == "" {
			fmt.Fprintf(os.Stderr, "The %s environment variable must be set to a valid Datadog %s key.\n", evKey, evType)

			return
		}
	}

	// Furnish the context with the necessary auth keys.
	ctx := context.WithValue(
		context.Background(),
		datadog.ContextAPIKeys,
		map[string]datadog.APIKey{
			"apiKeyAuth": {Key: os.Getenv(datadogAPIKey)},
			"appKeyAuth": {Key: os.Getenv(datadogApplicationKey)},
		},
	)

	// Hardcode to the Europe instance.
	ctx = context.WithValue(ctx,
		datadog.ContextServerVariables,
		map[string]string{
			"site": "datadoghq.eu",
		})

	// Set up a baseline client of the Datadog API.
	configuration := datadog.NewConfiguration()
	apiClient := datadog.NewAPIClient(configuration)

	// Prepare an SLO client and a request to list all SLOs.
	sloClient := datadogV1.NewServiceLevelObjectivesApi(apiClient)

	listResp, listR, err := sloClient.ListSLOs(ctx, datadogV1.ListSLOsOptionalParameters{})
	if err != nil {
		fmt.Fprintf(os.Stderr, "error listing SLOs: %v\n", err)
		fmt.Fprintf(os.Stderr, "full HTTP response: %v\n", listR)

		os.Exit(1)
	}

	defer listR.Body.Close()

	if listResp.HasErrors() {
		fmt.Fprintf(os.Stderr, "response has errors:\n")

		for index, err := range listResp.GetErrors() {
			fmt.Fprintf(os.Stderr, "[%03d] %s\n", index, err)
		}

		os.Exit(2)
	}

	// Look through the list of SLOs for a match by name.
	var targetSLOID string

	for _, slo := range listResp.GetData() {
		if strings.EqualFold(slo.GetName(), targetSLOName) {
			targetSLOID = slo.GetId()
			break
		}
	}

	if targetSLOID == "" {
		fmt.Fprintf(os.Stderr, "could not find SLO with name '%s'\n\nNOTE: this tool does not (yet) implement pagination, so if there are too many SLOs on Datadog then bother this tool's author to wire that up!\nhttps://github.com/DataDog/datadog-api-client-go?tab=readme-ov-file#pagination\n", targetSLOName)

		os.Exit(3)
	}

	// Get all details of the specific SLO by ID.
	getResp, getR, err := sloClient.GetSLO(ctx, targetSLOID, datadogV1.GetSLOOptionalParameters{})
	if err != nil {
		fmt.Fprintf(os.Stderr, "error getting SLO: %v\n", err)
		fmt.Fprintf(os.Stderr, "full HTTP response: %v\n", getR)

		os.Exit(4)
	}

	defer getR.Body.Close()

	if getResp.HasErrors() {
		fmt.Fprintf(os.Stderr, "response has errors:\n")

		for index, err := range getResp.GetErrors() {
			fmt.Fprintf(os.Stderr, "[%03d] %s\n", index, err)
		}

		os.Exit(5)
	}

	// Pretty-print the SLO data.
	sloData := getResp.GetData()

	indentedSLOBytes, err := json.MarshalIndent(sloData, "  ", "  ")
	if err != nil {
		fmt.Fprintf(os.Stderr, "marshalling indented SLO JSON: %v\n", err)

		os.Exit(6)
	}

	fmt.Fprintf(os.Stdout, "[\n  %s", indentedSLOBytes)

	// This should round off the pretty-print, regardless of whether the SLO has any monitors.
	defer fmt.Fprintf(os.Stdout, "\n]\n")

	// If the SLO does not have any monitors, we're done, so exit cleanly.
	if !sloData.HasMonitorIds() {
		return
	}

	// Prepare a Monitors client, get details of all monitors listed against the SLO, and pretty-print each one as we go.
	monitorsClient := datadogV1.NewMonitorsApi(apiClient)

	for _, monitor := range sloData.GetMonitorIds() {
		monResp, monR, err := monitorsClient.GetMonitor(ctx, monitor, datadogV1.GetMonitorOptionalParameters{})
		if err != nil {
			fmt.Fprintf(os.Stderr, "error getting Monitor: '%v'\n", err)

			os.Exit(7)
		}

		defer monR.Body.Close()

		indentedMonitorBytes, err := json.MarshalIndent(monResp, "  ", "  ")
		if err != nil {
			fmt.Fprintf(os.Stderr, "marshalling indented Monitor JSON: %v\n", err)

			os.Exit(8)
		}

		fmt.Fprintf(os.Stdout, ",\n  %s", indentedMonitorBytes)
	}
}
