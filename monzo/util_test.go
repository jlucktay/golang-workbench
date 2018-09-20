package main

import (
	"fmt"
	"net/url"
	"testing"
)

func TestConvertURL(t *testing.T) {
	// Domain filter given to each convertURL() call
	domain := "monzo.com"

	testCases := []struct {
		href      string
		urlResult url.URL
	}{
		{
			href:      "/",
			urlResult: url.URL{Scheme: "https", Host: domain, Path: "/"},
		},
		{
			href:      "/about",
			urlResult: url.URL{Scheme: "https", Host: domain, Path: "/about"},
		},
		{
			href:      "//monzo.com/about",
			urlResult: url.URL{Scheme: "https", Host: domain, Path: "/about"},
		},
		{
			href:      "http://monzo.com/about",
			urlResult: url.URL{Scheme: "http", Host: domain, Path: "/about"},
		},
		{
			href:      "//facebook.com/about",
			urlResult: url.URL{Scheme: "https", Host: "facebook.com", Path: "/about"},
		},
		{
			href:      "http://twitter.com/monzo",
			urlResult: url.URL{Scheme: "http", Host: "twitter.com", Path: "/monzo"},
		},
		{
			href:      "/fragment#foo",
			urlResult: url.URL{Scheme: "https", Host: domain, Path: "/fragment", Fragment: "foo"},
		},
		{
			href:      "/query?foo=bar",
			urlResult: url.URL{Scheme: "https", Host: domain, Path: "/query", RawQuery: "foo=bar"},
		},
		{
			href:      "monzo.com/bare-domain-same",
			urlResult: url.URL{Scheme: "https", Host: domain, Path: "/bare-domain-same"},
		},
		{
			href:      "monzo.com/bare/domain/same",
			urlResult: url.URL{Scheme: "https", Host: domain, Path: "/bare/domain/same"},
		},
		{
			href:      "twitter.com/bare-domain-different",
			urlResult: url.URL{Scheme: "https", Host: "twitter.com", Path: "/bare-domain-different"},
		},
		{
			href:      "twitter.com/bare/domain/different",
			urlResult: url.URL{Scheme: "https", Host: "twitter.com", Path: "/bare/domain/different"},
		},
		{
			href:      "samelevel",
			urlResult: url.URL{Scheme: "https", Host: domain, Path: "/samelevel"},
		},
	}
	for _, tC := range testCases {
		desc := fmt.Sprintf("[%s -> %s]", tC.href, tC.urlResult.String())
		t.Run(desc, func(t *testing.T) {
			actual := convertURL(tC.href, domain)
			if *actual != tC.urlResult {
				t.Errorf("Got '%q', want '%q'.\n[%#v]\n[%#v]\n", actual.String(), tC.urlResult.String(), actual, tC.urlResult)
			}
		})
	}
}
