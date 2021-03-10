package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

func main() {
	// Set account keys & information
	const (
		accountSid = "ACXXXX"
		authToken  = "XXXX"
		urlStr     = "https://api.twilio.com/2010-04-01/Accounts/" + accountSid + "/Messages.json"
	)

	// Create possible message bodies
	quotes := [7]string{
		"I urge you to please notice when you are happy, and exclaim or murmur or think at some point, " +
			"'If this isn't nice, I don't know what is.'",
		"Peculiar travel suggestions are dancing lessons from God.",
		"There's only one rule that I know of, babies—God damn it, you've got to be kind.",
		"Many people need desperately to receive this message: " +
			"'I feel and think much as you do, care about many of the things you care about, " +
			"although most people do not care about them. You are not alone.'",
		"That is my principal objection to life, I think: It's too easy, when alive, to make perfectly horrible mistakes.",
		"So it goes.",
		"We must be careful about what we pretend to be.",
	}

	// Set up rand
	rand.Seed(time.Now().Unix())

	// Pack up the data for our message
	msgData := url.Values{}
	msgData.Set("To", "NUMBER_TO")
	// msgData.Set("From", "NUMBER_FROM")
	msgData.Set("MessagingServiceSid", "XXXX")
	msgData.Set("Body", quotes[rand.Intn(len(quotes))])
	msgDataReader := strings.NewReader(msgData.Encode())

	// Create HTTP request
	req, err := http.NewRequest(http.MethodPost, urlStr, msgDataReader)
	if err != nil {
		fmt.Fprintf(os.Stderr, "could not create request: %v\n", err)
		os.Exit(1)
	}

	req.SetBasicAuth(accountSid, authToken)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	// Make HTTP POST request and return message SID
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Fprintf(os.Stderr, "could not perform request: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Response status: (%d) %s\n", resp.StatusCode, resp.Status)

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		fmt.Fprintln(os.Stderr, "non-OK status code from response")

		b := []byte{}
		defer resp.Body.Close()

		if _, err := resp.Body.Read(b); err != nil {
			fmt.Fprintf(os.Stderr, "could not read response body: %v\n", err)
			os.Exit(1)
		}

		fmt.Fprintf(os.Stderr, "Response body: '%s'\n", b)
		os.Exit(1)
	}

	var data map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		fmt.Fprintf(os.Stderr, "could not decode response body: %v\n", err)
		os.Exit(1)
	}

	fmt.Println(data["sid"])
}
