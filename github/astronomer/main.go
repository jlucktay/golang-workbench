package main

import (
	"container/heap"
	"context"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/shurcooL/githubv4"
	"golang.org/x/oauth2"
)

const (
	fmtPkgImp = "https://pkg.go.dev/%s?tab=importedby"
	envToken  = "GITHUB_TOKEN"
)

const (
	ExitSuccess int = iota
	ExitNoGitHubToken
	ExitBadArgs
	ExitParsingURL
	ExitGettingPackages
	ExitReadingResponseBody
)

func main() {
	// Make sure the GitHub token is set in the environment.
	token, tokenSet := os.LookupEnv(envToken)
	if !tokenSet {
		fmt.Fprintf(os.Stderr, "token not set in environment: %s\n", envToken)

		os.Exit(ExitNoGitHubToken)
	}

	// Make sure we were given exactly one string argument.
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "must provide exactly one string argument\n")

		os.Exit(ExitBadArgs)
	}

	rawURL := fmt.Sprintf(fmtPkgImp, os.Args[1])
	impURL, err := url.Parse(rawURL)
	if err != nil {
		fmt.Fprintf(os.Stderr, "could not parse URL '%s': %v\n", rawURL, err)

		os.Exit(ExitParsingURL)
	}

	// Get the URL using the arg.
	http.DefaultClient.Timeout = 5 * time.Second
	resp, err := http.Get(impURL.String())
	if err != nil || resp.StatusCode != 200 {
		fmt.Fprintf(os.Stderr, "getting other packages that import this one from '%s'; status '%s'; error: %v\n",
			impURL, resp.Status, err)

		os.Exit(ExitGettingPackages)
	}

	defer resp.Body.Close()

	// Parse the HTML and get a list of packages.
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		fmt.Fprintf(os.Stderr, "reading HTML document from response: %v\n", err)

		os.Exit(ExitReadingResponseBody)
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
	httpClient.Timeout = 5 * time.Second
	client := githubv4.NewClient(httpClient)

	stargazers := make(map[string]int)

	for i, impBy := range importedBy {
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

		if err := client.Query(context.TODO(), &query, vars); err != nil {
			if strings.Contains(err.Error(), "API rate limit exceeded") {
				fmt.Println("Rate limit exhausted.")
				fmt.Println("https://docs.github.com/en/graphql/overview/resource-limitations")
				return
			}

			fmt.Printf(" is a dud (error: %v)\n", err)
			stargazers[sgKey] = 0
			continue
		}

		fmt.Printf("[%04d] %v", i, vars)
		fmt.Printf(" %d stars", query.Repository.StargazerCount)
		stargazers[sgKey] = query.Repository.StargazerCount

		untilReset := time.Until(query.RateLimit.ResetAt.Time)
		if int64(query.RateLimit.Remaining) <= 0 {
			fmt.Printf("; rate limit exhausted. Resets at: %s\n", query.RateLimit.ResetAt.Format(time.RFC3339))
			fmt.Println("https://docs.github.com/en/graphql/overview/resource-limitations")
			fmt.Printf("Sleeping %s until reset...", untilReset)
			time.Sleep(untilReset)
			fmt.Println()
		} else {
			sleepDuration := time.Duration(untilReset.Nanoseconds() / int64(query.RateLimit.Remaining))
			fmt.Printf(" (cost %d, %d/%d remaining, sleep for %s)\n",
				query.RateLimit.Cost, query.RateLimit.Remaining, query.RateLimit.Limit,
				sleepDuration.Truncate(time.Millisecond))
			time.Sleep(sleepDuration)
		}
	}

	h := getHeap(stargazers)

	for i := 1; i <= 100 && i <= h.Len(); i++ {
		popped := heap.Pop(h).(kv)
		fmt.Printf("%3d. %d %s\n", i, popped.Value, popped.Key)
	}
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
	*h = append(*h, x.(kv))
}

func (h *KVHeap) Pop() any {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]

	return x
}

func getHeap(m map[string]int) *KVHeap {
	h := &KVHeap{}
	heap.Init(h)

	for k, v := range m {
		heap.Push(h, kv{Key: k, Value: v})
	}

	return h
}
