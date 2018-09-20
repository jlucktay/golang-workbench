package main

import (
	"fmt"
	"net/url"
	"testing"
)

func TestConvertURL(t *testing.T) {
	domain := "monzo.com"

	testCases := []struct {
		href      string
		urlResult url.URL
	}{
		{
			href:      "/",
			urlResult: url.URL{Scheme: "https", Host: "monzo.com", Path: "/"},
		},
		{
			href:      "/about",
			urlResult: url.URL{Scheme: "https", Host: "monzo.com", Path: "/about"},
		},
		{
			href:      "//monzo.com/about",
			urlResult: url.URL{Scheme: "https", Host: "monzo.com", Path: "/about"},
		},
		{
			href:      "http://monzo.com/about",
			urlResult: url.URL{Scheme: "http", Host: "monzo.com", Path: "/about"},
		},
		{
			href:      "//facebook.com/about",
			urlResult: url.URL{Scheme: "https", Host: "facebook.com", Path: "/about"},
		},
		{
			href:      "http://twitter.com/monzo",
			urlResult: url.URL{Scheme: "http", Host: "twitter.com", Path: "/monzo"},
		},
	}
	for _, tC := range testCases {
		desc := fmt.Sprintf("[%s][%s]", tC.href, tC.urlResult.String())
		t.Run(desc, func(t *testing.T) {
			actual := convertURL(tC.href, domain)
			if *actual != tC.urlResult {
				t.Errorf("Got '%q', want '%q'.\n", actual.String(), tC.urlResult.String())
			}
		})
	}
}
