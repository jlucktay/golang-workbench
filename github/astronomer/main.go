package main

import (
	"container/heap"
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/shurcooL/githubv4"
	"golang.org/x/oauth2"
)

const (
	envToken = "GITHUB_TOKEN"
)

func main() {
	// Set up GitHub GraphQL API v4 client
	token, tokenSet := os.LookupEnv(envToken)
	if !tokenSet {
		fmt.Fprintf(os.Stderr, "token not set in environment: %s\n", envToken)

		return
	}

	src := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
	httpClient := oauth2.NewClient(context.TODO(), src)
	httpClient.Timeout = 5 * time.Second
	client := githubv4.NewClient(httpClient)

	content, err := os.ReadFile("list.txt")
	if err != nil {
		log.Fatal(err)
	}

	lines := strings.Split(string(content), "\n")

	stargazers := map[string]int{}

	for i := range lines {
		if !strings.HasPrefix(lines[i], "github.com/") {
			continue
		}

		lines[i] = strings.TrimPrefix(lines[i], "github.com/")

		xLine := strings.Split(lines[i], "/")

		if len(xLine) < 2 {
			continue
		}

		vars := map[string]interface{}{
			"owner": githubv4.String(xLine[0]),
			"name":  githubv4.String(xLine[1]),
		}

		fmt.Printf("[%04d] %v ", i, vars)

		var query queryStargazers

		if err := client.Query(context.TODO(), &query, vars); err != nil {
			fmt.Printf("is a dud (error: %v)\n", err)
			continue
		} else if query.RateLimit.Remaining <= 100 {
			fmt.Printf("- rate limit remaining is low: %d/%d\n", query.RateLimit.Remaining, query.RateLimit.Limit)
			time.Sleep(time.Second)
		}

		fmt.Printf("%d stars", query.Repository.StargazerCount)
		stargazers[strings.ToLower(xLine[0]+"/"+xLine[1])] = query.Repository.StargazerCount

		untilReset := time.Until(query.RateLimit.ResetAt.Time)
		sleepDuration := time.Duration(untilReset.Nanoseconds() / int64(query.RateLimit.Remaining))
		fmt.Printf(" (cost %d, %d/%d remaining, sleep for %s)\n",
			query.RateLimit.Cost, query.RateLimit.Remaining, query.RateLimit.Limit,
			sleepDuration.Truncate(time.Millisecond))
		time.Sleep(sleepDuration)
	}

	h := getHeap(stargazers)

	for i := 1; i <= 100; i++ {
		popped := heap.Pop(h).(kv)
		fmt.Printf("%3d. %d %s\n", i, popped.Value, popped.Key)
	}
}

type queryStargazers *struct {
	Repository *struct {
		StargazerCount int
	} `graphql:"repository(name: $name, owner: $owner)"`
	RateLimit struct {
		ResetAt   githubv4.DateTime
		Cost      githubv4.Int
		Limit     githubv4.Int
		Remaining githubv4.Int
	}
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

func (h *KVHeap) Push(x interface{}) {
	*h = append(*h, x.(kv))
}

func (h *KVHeap) Pop() interface{} {
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
