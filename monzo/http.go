package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
)

func getResponse(get url.URL) io.Reader {
	req, reqErr := http.NewRequest("GET", get.String(), nil)
	if reqErr != nil {
		log.Fatal(reqErr)
	}

	req.Header.Add("User-Agent", "jlucktay (monzo-crawler)")

	res, resErr := http.DefaultClient.Do(req)
	if resErr != nil {
		log.Fatal(resErr)
	}
	defer res.Body.Close()

	buf := new(bytes.Buffer)

	if res.StatusCode != 200 {
		fmt.Printf("[getResponse] URL '%+v': status code error: [%d] %s", get, res.StatusCode, res.Status)
		return buf
	}

	buf.ReadFrom(res.Body)
	return buf
}
