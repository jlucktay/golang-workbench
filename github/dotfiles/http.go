package main

import (
	"bytes"
	"io"
	"log"
	"net/http"
	"net/url"
)

func getResponse(get url.URL) io.Reader {
	res, err := http.DefaultClient.Do(newRequest(get))
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		log.Fatalf("URL '%s': status code error: %d %s", get.String(), res.StatusCode, res.Status)
	}

	buf := new(bytes.Buffer)
	buf.ReadFrom(res.Body)

	return buf
}

func newRequest(u url.URL) (req *http.Request) {
	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Add("User-Agent", "jlucktay (dotfiles)")
	req.SetBasicAuth("jlucktay", ghpaToken)

	return
}
