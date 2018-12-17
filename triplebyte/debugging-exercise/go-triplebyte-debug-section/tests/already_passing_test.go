package main

import (
	"testing"

	"net/url"

	"../crawler"
	"../htmlhelp"
)

func TestHtmlhelp(t *testing.T) {
	root, err := url.Parse("http://example.com")
	if err != nil {
		t.Fatal("Broken test, can't parse root url")
	}

	doc := `
<!DOCTYPE html>
<html>
  <body>
    <h1>Test Case 1</h1>

    <p>I am a paragraph! <a href="javascript:doThing">blah</a></p>

    <p>Sometimes I am <a href="./cynical.html">overly cynical</a>, but sometimes I am
      <a href="./page2.html">overly na&#xEFve.</a></p>
  </body>
</html>
`
	ns, errors := htmlhelp.Neighbors(doc, *root)
	if len(errors) > 0 {
		t.Errorf("unexpected errors: %v", err)
	}

	expect := []string{"javascript:doThing", "http://example.com/cynical.html", "http://example.com/page2.html"}
	if len(ns) != len(expect) {
		t.Errorf("unexpected neighbors: %v", ns)
		return
	}

	for i, u := range ns {
		if u.String() != expect[i] {
			t.Errorf("neighbor mismatch: %s %s", u.String(), expect[i])
		}
	}

}

func TestCrawler(t *testing.T) {
	c := crawler.Crawler{100, crawler.Silent()}
	graph, err := c.Crawl("http://triplebyte.github.io/web-crawler-test-site/already-passing-tests/", "")
	if err != nil {
		t.Fatalf("can't crawl: %v", err)
	}

	tests := []struct {
		u      string
		status crawler.NodeStatus
		code   int
	}{
		{"http://triplebyte.github.io/web-crawler-test-site/already-passing-tests/page2", crawler.SUCCESS, 200},
		{"http://triplebyte.github.io/web-crawler-test-site/already-passing-tests/page2-real", crawler.SUCCESS, 200},
		{"http://triplebyte.github.io/web-crawler-test-site/already-passing-tests/page2-fake", crawler.SUCCESS, 404},
	}

	for _, tc := range tests {
		u, err := url.Parse(tc.u)
		if err != nil {
			t.Fatalf("Broken test: %q not a URL", tc.u)
		}

		n, ok := graph.Nodes[*u]
		switch {
		case !ok:
			t.Errorf("Crawl(%q): not found", u)
		case n.Status != tc.status:
			t.Errorf("Crawl(%q).Status: '%d' '%d'", u, tc.status, n.Status)
		case n.Code != tc.code:
			t.Errorf("Crawl(%q).Code: '%d' '%d'", u, tc.code, n.Code)
		}

	}
}
