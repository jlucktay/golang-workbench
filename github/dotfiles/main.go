package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

var (
	ghpaToken = readTokenFromSecrets()
)

func main() {
	dotfilesURL, _ := url.Parse("https://dotfiles.github.io")
	ghRegex := `^https://github.com/[A-Za-z0-9_\-\.]+/[A-Za-z0-9_\-\.]+$`

	doc, err := goquery.NewDocumentFromReader(getResponse(*dotfilesURL))
	if err != nil {
		log.Fatal(err)
	}

	for x := range genGoRepos(genAPILang(genLinks(doc, ghRegex))) {
		fmt.Println(x)
	}
}

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

type secrets struct {
	GitHubPersonalAccessToken string
}

func readTokenFromSecrets() (token string) {
	fileContents, err := ioutil.ReadFile("./secrets.json")
	if err != nil {
		log.Fatal(err)
	}

	var tokenMap map[string]string

	if err := json.Unmarshal(fileContents, &tokenMap); err != nil {
		log.Fatal(err)
	}

	token = tokenMap["GitHubPersonalAccessToken"]

	fmt.Println(token)

	return
}

func genLinks(input *goquery.Document, filter string) (output chan url.URL) {
	output = make(chan url.URL)

	go func() {
		sentOne := false // TODO: delete this

		input.Find("a").Each(func(i int, s *goquery.Selection) {
			src := s.AttrOr("href", "")
			if u, _ := url.Parse(src); u.IsAbs() && u.Scheme != "data" && matchStringCompiled(filter, u.String()) && !sentOne {
				output <- *u
				sentOne = true // TODO: delete this
			}
		})

		close(output)
	}()

	return
}

func genAPILang(in chan url.URL) (output chan url.URL) {
	output = make(chan url.URL)

	go func() {
		for i := range in {
			p := strings.Split(i.Path, "/")
			langURL, err := url.Parse(fmt.Sprintf("https://api.github.com/repos/%s/%s/languages", p[1], p[2]))
			if err != nil {
				log.Fatal(err)
			}

			output <- *langURL
		}

		close(output)
	}()

	return
}

func genGoRepos(input chan url.URL) (output chan string) {
	// maybe do this with a map instead? use the repo URL as the key?
	output = make(chan string)

	go func() {
		for i := range input {
			// make API request
			// print language(s)

			output <- fmt.Sprint(getResponse(i))
		}

		close(output)
	}()

	return
}

func matchStringCompiled(needle, haystack string) bool {
	r, err := regexp.Compile(needle)
	if err != nil {
		panic(err)
	}

	return r.MatchString(haystack)
}
