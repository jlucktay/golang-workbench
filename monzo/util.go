package main

import (
	"fmt"
	"net/url"
	"strings"
)

func convertURL(base, input string) *url.URL {
	// Found an odd youtube.com link when crawling google.com that started with a space 🤔
	input = strings.TrimSpace(input)

	// If they are just anchor links on the same page, disregard and return nil
	if strings.HasPrefix(input, "#") {
		return nil
	}

	// We need a base URL to resolve potential relative links further down the line
	urlBase, errBaseParse := url.Parse(base)
	if errBaseParse != nil {
		errorLog.Printf("Error parsing base URL '%s': %v\n", base, errBaseParse)
		return nil
	}

	prefix := ""

	// The assumption for relative URLs on the same domain is that they are all secure, hence prepending 'https'
	if strings.HasPrefix(input, "/") {
		if strings.HasPrefix(input, "//") {
			if urlBase.IsAbs() {
				prefix = urlBase.Scheme + ":"
			} else {
				prefix = "https:"
			}
		} else {
			if urlBase.IsAbs() {
				prefix = fmt.Sprintf("%s://%s", urlBase.Scheme, flagURL)
			} else {
				prefix = fmt.Sprintf("https://%s", flagURL)
			}
		}
	}

	urlOut, errParse := url.Parse(prefix + input)
	if errParse != nil {
		errorLog.Printf("Error parsing input '%s': %v\n", prefix+input, errParse)
		return nil
	}

	// Make non-absolute URLs into absolute
	if !urlOut.IsAbs() {
		split := strings.SplitN(urlOut.String(), "/", 2)

		if split[0] == "." || split[0] == ".." {
			urlOut.Host = ""
		} else {
			urlOut.Scheme = "https"

			if len(split) < 2 {
				urlOut.Host = flagURL
				urlOut.Path = "/" + split[0]
			} else {
				urlOut.Host = split[0]
				urlOut.Path = "/" + split[1]
			}
		}
	}

	urlOut = urlBase.ResolveReference(urlOut)

	for strings.Contains(urlOut.Path, "//") {
		urlOut.Path = strings.Replace(urlOut.Path, "//", "/", -1)
	}

	return urlOut
}
