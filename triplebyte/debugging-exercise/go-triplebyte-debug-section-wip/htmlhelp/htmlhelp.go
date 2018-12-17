package htmlhelp

import (
	"fmt"
	"net/url"
	"regexp"
	"strings"
)

// Browsers do all kinds of crazy things to make their hrefs valid even when
// they're sketchy. For example, spaces get autoescaped on Chrome.
func cleanup(href string) string {
	return strings.Replace(href, " ", "%20", 0)
}

func neighborsFromPaths(paths []string, ref url.URL) (neighbors []url.URL, errors []string) {
	for _, p := range paths {
		p = cleanup(p)

		// url.Parse silently fixes exactly the sort of thing
		// that we want to report on, so we have to check not
		// only if there was a parse error, but if the parse
		// fixed an error that we want to note and report.
		check, err := url.Parse(p)
		if err != nil || check.String() != p {
			msg := fmt.Sprintf("The page %v has an href of %s, which is not a valid URI", &ref, p)
			errors = append(errors, msg)
			continue
		}

		u, err := ref.Parse(p)
		if u.Scheme == "mailto" {
			continue
		}

		u.Fragment = "" // Don't send to the server in the next trip

		neighbors = append(neighbors, *u)
	}

	return
}

func paths(doc string) []string {
	var ret []string

	// patterns must have the url as the first capturing subgroup
	aPattern := regexp.MustCompile(`<a [^>]*href="([^"#]*)#?[^"]*"`)
	scriptPattern := regexp.MustCompile(`<script [^>]*src="([^"]*)"`)
	linkPattern := regexp.MustCompile(`<link [^>]*href="([^"]*)"`)

	for _, pat := range []*regexp.Regexp{aPattern, linkPattern, scriptPattern} {
		found := pat.FindAllStringSubmatch(doc, -1)
		chunk := make([]string, len(found))
		for i, match := range found {
			chunk[i] = match[1]
		}
		ret = append(ret, chunk...)
	}

	return ret
}

// Neighbors processes a url
func Neighbors(doc string, ref url.URL) (neighbors []url.URL, errors []string) {
	ps := paths(doc)
	return neighborsFromPaths(ps, ref)
}
