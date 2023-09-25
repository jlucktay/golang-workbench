package main

import (
	"bytes"
	"container/heap"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/shurcooL/githubv4"
	"golang.org/x/oauth2"
	"golang.org/x/time/rate"
)

const (
	fmtPkgImp     = "https://pkg.go.dev/%s?tab=importedby"
	envToken      = "GITHUB_TOKEN"
	tenSecTimeout = 10 * time.Second
)

const (
	ExitSuccess int = iota
	ExitNoGitHubToken
	ExitBadArgs
	ExitParsingURL
	ExitCreatingRequest
	ExitGettingPackages
	ExitReadingResponseBody
	ExitParsingResponseBody
	ExitRateLimitExceeded
)

func main() {
	os.Exit(run())
}

func run() int {
	// Make sure the GitHub token is set in the environment.
	token, tokenSet := os.LookupEnv(envToken)
	if !tokenSet {
		fmt.Fprintf(os.Stderr, "token not set in environment: %s\n", envToken)

		return ExitNoGitHubToken
	}

	// Make sure we were given exactly one string argument.
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "must provide exactly one string argument\n")

		return ExitBadArgs
	}

	rawURL := fmt.Sprintf(fmtPkgImp, os.Args[1])

	impURL, err := url.Parse(rawURL)
	if err != nil {
		fmt.Fprintf(os.Stderr, "could not parse URL '%s': %v\n", rawURL, err)

		return ExitParsingURL
	}

	// Get the URL using the arg.
	timeoutCtx, cancelFunc := context.WithTimeout(context.TODO(), tenSecTimeout)
	defer cancelFunc()

	getImpURLs, err := http.NewRequestWithContext(timeoutCtx, http.MethodGet, impURL.String(), nil)
	if err != nil {
		fmt.Fprintf(os.Stderr, "could not create HTTP request: %v\n", err)

		return ExitCreatingRequest
	}

	resp, err := http.DefaultClient.Do(getImpURLs)
	if err != nil || resp.StatusCode < 200 || resp.StatusCode >= 300 {
		var urlError *url.Error
		if errors.As(err, &urlError) && urlError.Timeout() {
			fmt.Fprint(os.Stderr, "request timed out; ")
		}

		if resp != nil {
			fmt.Fprintf(os.Stderr, "response status '%v'; ", resp.Status)
		}

		fmt.Fprintf(os.Stderr, "getting other packages that import this one from '%s'; error: %v\n", impURL, err)

		return ExitGettingPackages
	}
	defer resp.Body.Close()

	bodyBytes := &bytes.Buffer{}
	fmt.Print("reading bytes from response body")

	for {
		_, err = io.CopyN(bodyBytes, resp.Body, 1024)
		fmt.Print(".")

		if errors.Is(err, io.EOF) {
			break
		} else if err != nil {
			fmt.Fprintf(os.Stderr, "reading bytes from response body: %v\n", err)

			return ExitReadingResponseBody
		}
	}

	fmt.Println()

	// Parse the HTML and get a list of packages.
	doc, err := goquery.NewDocumentFromReader(bodyBytes)
	if err != nil {
		fmt.Fprintf(os.Stderr, "parsing HTML document from response: %v\n", err)

		return ExitParsingResponseBody
	}

	importedBy := make([]*url.URL, 0)

	doc.Find("li.ImportedBy-detailsIndent a.u-breakWord").Each(func(_ int, s *goquery.Selection) {
		val, exists := s.Attr("href")
		if !exists {
			return
		}

		val = strings.TrimPrefix(val, "/")
		if !strings.HasPrefix(val, "github.com/") {
			return
		}

		parsed, err := url.Parse(val)
		if err != nil {
			return
		}

		parsed.Scheme = "https"
		importedBy = append(importedBy, parsed)
	})

	// Set up a GitHub GraphQL client to get stars.
	src := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
	httpClient := oauth2.NewClient(context.TODO(), src)
	httpClient.Timeout = tenSecTimeout
	client := githubv4.NewClient(httpClient)
	stargazers := make(map[string]int)

	// This limit is accurate as of 2023-09-25.
	// cf. https://docs.github.com/en/graphql/overview/resource-limitations#rate-limit
	rl := rate.NewLimiter(rate.Every(1388*(time.Millisecond)), 1)

	for index, impBy := range importedBy {
		trimmed := strings.TrimPrefix(impBy.String(), "https://github.com/")
		xLine := strings.Split(trimmed, "/")

		if len(xLine) < 2 {
			continue
		}

		sgKey := strings.ToLower(xLine[0] + "/" + xLine[1])
		if _, alreadySeen := stargazers[sgKey]; alreadySeen {
			continue
		}

		vars := map[string]any{
			"owner": githubv4.String(xLine[0]),
			"name":  githubv4.String(xLine[1]),
		}

		query := &queryStargazers{}

		fmt.Printf("[%04d] %v", index, vars)

		if err := client.Query(context.TODO(), &query, vars); err != nil {
			if strings.Contains(err.Error(), "API rate limit exceeded") {
				fmt.Fprintln(os.Stderr, "Rate limit exhausted.")
				fmt.Fprintln(os.Stderr, "https://docs.github.com/en/graphql/overview/resource-limitations")

				return ExitRateLimitExceeded
			}

			fmt.Printf(" is a dud (error: %v)\n", err)

			_ = rl.WaitN(context.Background(), int(query.RateLimit.Cost))
			stargazers[sgKey] = 0

			continue
		}

		fmt.Printf(" %d stars", query.Repository.StargazerCount)
		stargazers[sgKey] = query.Repository.StargazerCount

		rlWaitStart := time.Now()
		err := rl.WaitN(context.Background(), int(query.RateLimit.Cost))
		rlWaitFinish := time.Now()

		if err != nil {
			fmt.Fprintf(os.Stderr, "rate limiter waiting: %v\n", err)

			return ExitRateLimitExceeded
		}

		sleptDuration := rlWaitFinish.Sub(rlWaitStart)

		fmt.Printf(" (cost %d, %d/%d remaining, slept for %s)\n",
			query.RateLimit.Cost, query.RateLimit.Remaining, query.RateLimit.Limit,
			sleptDuration.Truncate(time.Millisecond))
	}

	h := getHeap(stargazers)

	for i := 1; i <= 100 && i <= h.Len(); i++ {
		if popped, ok := heap.Pop(h).(kv); ok {
			fmt.Printf("%3d. %d %s\n", i, popped.Value, popped.Key)
		}
	}

	return ExitSuccess
}

type queryStargazers struct {
	Repository struct {
		StargazerCount int
	} `graphql:"repository(name: $name, owner: $owner)"`
	RateLimit queryRateLimit
}

type queryRateLimit struct {
	Limit     githubv4.Int
	Cost      githubv4.Int
	Remaining githubv4.Int
	ResetAt   githubv4.DateTime
}

type kv struct {
	Key   string
	Value int
}

type KVHeap []kv

// Less uses greater-than so we can pop *larger* items.
func (h KVHeap) Less(i, j int) bool { return h[i].Value > h[j].Value }
func (h KVHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }
func (h KVHeap) Len() int           { return len(h) }

func (h *KVHeap) Push(x any) {
	if y, ok := x.(kv); ok {
		*h = append(*h, y)
	}
}

func (h *KVHeap) Pop() any {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]

	return x
}

func getHeap(m map[string]int) *KVHeap {
	kvh := &KVHeap{}
	heap.Init(kvh)

	for k, v := range m {
		heap.Push(kvh, kv{Key: k, Value: v})
	}

	return kvh
}
