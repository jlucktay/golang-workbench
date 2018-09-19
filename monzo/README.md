# Monzo

See the [design](DESIGN.md) and [spec](SPEC.md) Markdown files for some more details.

You can supply a URL as a command line argument like so:

``` shell
$ go run . -url monzo.com
[genCrawl] Fetched 'https://monzo.com'.
[genCrawl] 'https://monzo.com/' is a child of 'https://monzo.com'.
[genCrawl] 'https://monzo.com/about' is a child of 'https://monzo.com'.
[genCrawl] 'https://monzo.com/blog' is a child of 'https://monzo.com'.
...
```
