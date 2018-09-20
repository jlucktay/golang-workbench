package main

import (
	"fmt"
	"log"
	"net/url"
	"strings"
)

func convertURL(input, domain string) *url.URL {
	prefix := ""

	// The assumption for relative URLs on the same domain is that they are all secure, hence prepending 'https'
	if strings.HasPrefix(input, "/") {
		if strings.HasPrefix(input, "//") {
			prefix = "https:"
		} else {
			prefix = fmt.Sprintf("https://%s", domain)
		}
	}

	urlOut, errParse := url.Parse(prefix + input)
	if errParse != nil {
		log.Fatalf("Error parsing '%s': %v\n", prefix+input, errParse)
	}

	return urlOut
}
