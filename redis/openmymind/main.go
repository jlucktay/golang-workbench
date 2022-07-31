package main

import (
	"context"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/caarlos0/env/v6"
	"github.com/go-redis/redis/v9"
)

type config struct {
	RSSURL string `env:"RSS_URL"`
}

func main() {
	//
	// config

	cfg := config{}
	opt := env.Options{
		Prefix:          "OPENMYMIND_",
		RequiredIfNoDef: true,
	}
	if err := env.Parse(&cfg, opt); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return
	}

	//
	// redis

	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	ctx := context.Background()

	ping, err := rdb.Ping(ctx).Result()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return
	}
	if ping != "PONG" {
		fmt.Fprintf(os.Stderr, "%s\n", ping)
		return
	}

	defer rdb.Close()

	//
	// xml

	resp, err := http.Get(cfg.RSSURL)
	if err != nil {
		fmt.Fprintf(os.Stderr, "http.Get error: %v\n", err)
		return
	}

	rss := Rss{}
	if err := xml.NewDecoder(resp.Body).Decode(&rss); err != nil {
		fmt.Fprintf(os.Stderr, "xml.Decode error: %v\n", err)
		return
	}

	resp.Body.Close()

	//
	// loop

	for i := range rss.Channel.Item {
		tj, err := rss.Channel.Item[i].ToJSON()
		if err != nil {
			fmt.Fprintf(os.Stderr, "%v\n", err)

			return
		}

		score, err := rss.Channel.Item[i].CreatedAt()
		if err != nil {
			fmt.Fprintf(os.Stderr, "%v\n", err)

			return
		}

		key := strings.Join([]string{
			"item",
			rss.Channel.Item[i].Guid.Text,
		}, ":")

		if err := rdb.Set(ctx, key, tj, 0).Err(); err != nil && err != redis.Nil {
			fmt.Fprintf(os.Stderr, "%v\n", err)

			return
		}

		if err := rdb.ZAdd(ctx, "itemz", redis.Z{Score: score, Member: key}).Err(); err != nil && err != redis.Nil {
			fmt.Fprintf(os.Stderr, "%v\n", err)

			return
		}
	}

	zrevrange := rdb.ZRevRange(ctx, "itemz", 0, 2)
	if zrevrange.Err() != nil && zrevrange.Err() != redis.Nil {
		fmt.Fprintf(os.Stderr, "%v\n", zrevrange.Err())

		return
	}

	for _, j := range zrevrange.Val() {
		get := rdb.Get(ctx, j)
		if get.Err() != nil && get.Err() != redis.Nil {
			fmt.Fprintf(os.Stderr, "%v\n", err)

			return
		}

		dec := json.NewDecoder(strings.NewReader(get.Val()))
		for {
			var i Item
			if err := dec.Decode(&i); err == io.EOF {
				break
			} else if err != nil {
				fmt.Fprintf(os.Stderr, "%v\n", err)

				return
			}

			fmt.Printf("%s (%s)\n", i.Title, i.PubDate)
		}
	}
}
