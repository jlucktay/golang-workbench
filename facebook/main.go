package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"text/tabwriter"
	"time"
)

func main() {
	raw, errRead := ioutil.ReadFile("message_1.json")
	if errRead != nil {
		panic(errRead)
	}

	var fbExport FacebookExport

	if errUm := json.Unmarshal(raw, &fbExport); errUm != nil {
		panic(errUm)
	}

	fmt.Printf("Message count: %d\n\n", len(fbExport.Messages))

	tw := tabwriter.NewWriter(os.Stdout, 1, 2, 2, ' ', tabwriter.TabIndent)

	fmt.Fprintln(tw, "Caller\tDuration\tStart time\tFinish time")

	for _, msg := range fbExport.Messages {
		switch msg.Type {
		case FMTCall:
			if !msg.Missed {
				finishTime := epochMsToHumanReadable(msg.TimestampMs)

				duration, errParse := time.ParseDuration(fmt.Sprintf("%ds", msg.CallDuration))
				if errParse != nil {
					panic(errParse)
				}

				startTime := finishTime.Add(-1 * duration)

				fmt.Fprintf(tw, "%s\t%s\t%s\t%s\n", msg.SenderName, duration, startTime, finishTime)
			}

		case FMTGeneric:

		case FMTShare:

		default:
			panic("unknown Facebook message type")
		}
	}

	tw.Flush()
}

func epochMsToHumanReadable(epoch int64) time.Time {
	return time.Unix(epoch/1000, 0)
}
