package main

import (
	"fmt"
	"net/url"
	"strings"
)

// TODO: reinstate the error handling for the Parse() calls
func convertURL(input, domain string) (urlOut *url.URL) {
	// The assumption for relative URLs on the same domain are that they are all secure, hence prepending 'https'
	if strings.HasPrefix(input, "/") {
		if strings.HasPrefix(input, "//") {
			urlOut, _ = url.Parse(fmt.Sprintf("https:%s", input))
			// if errOne != nil {
			// 	log.Fatal(errOne)
			// }
		} else {
			urlOut, _ = url.Parse(fmt.Sprintf("https://%s%s", domain, input))
			// if errTwo != nil {
			// 	log.Fatal(errTwo)
			// }
		}
	} else {
		urlOut, _ = url.Parse(input)
		// if errThree != nil {
		// 	log.Fatal(errThree)
		// }
	}

	return
}
