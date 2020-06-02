package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"text/tabwriter"
)

func main() {
	raw, errRead := ioutil.ReadFile("messages.json")
	if errRead != nil {
		panic(errRead)
	}

	var skypeExport SkypeExport

	if errUm := json.Unmarshal(raw, &skypeExport); errUm != nil {
		panic(errUm)
	}

	fmt.Printf("Conversations: %d\n\n", len(skypeExport.Conversations))

	tw := tabwriter.NewWriter(os.Stdout, 1, 2, 2, ' ', tabwriter.TabIndent)

	fmt.Fprintln(tw, "From\tOriginalArrivalTime")

	callSum := 0

	for _, conversation := range skypeExport.Conversations {
		if strings.Contains(strings.ToLower(conversation.DisplayName), secret) {
			for _, message := range conversation.MessageList {
				if message.MessageType == "Event/Call" {
					fmt.Fprintf(tw, "%s\t%s\n", message.From, message.OriginalArrivalTime)

					callSum++
				}
			}
		}
	}

	tw.Flush()

	fmt.Printf("\nCall count: %d\n", callSum)
}
