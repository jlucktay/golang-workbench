# Notes

## Timing

Upon reaching feature-complete status:

``` shell
$ time go run . -url monzo.com
Pages crawled: 1225
Pages outside target 'https://monzo.com' domain: 9477

real    0m4.941s
user    0m4.246s
sys     0m2.553s
$
```

After adding unit test coverage to `convertURL()` and refactoring/adding logic:

``` shell
$ time go run . -url monzo.com
Pages crawled: 1190
Pages outside target 'https://monzo.com' domain: 8280

real    0m8.197s
user    0m4.381s
sys     0m2.462s
```

Some further refactoring in various spots across the whole package gained some speed back:

``` shell
$ time go run . -url monzo.com
Pages crawled: 1202
Pages outside target 'https://monzo.com' domain: 8767

real    0m4.982s
user    0m4.058s
sys     0m2.204s
```

Adding the relative path and HTTP/HTTPS protocol handling to `convertURL()` did slow it back down, although the number of pages covered/crawled went up significantly:

``` shell
$ time go run . -url monzo.com
Pages crawled: 2315
Pages outside target 'https://monzo.com' domain: 16805

real    0m10.043s
user    0m7.713s
sys     0m4.701s
```

Following some more refactoring around the error handling and locking section in `crawl()`:

``` shell
$ time go run . -url monzo.com
Pages crawled: 2283
Pages outside target 'https://monzo.com' domain: 16530

real    0m8.934s
user    0m7.990s
sys     0m4.663s
```

## References

- [Golang.org](https://golang.org)
  - Notably the [language spec](https://golang.org/ref/spec) and the [package docs and examples](https://golang.org/pkg/)
- [Another webcrawler I wrote previously](https://github.com/jlucktay/golang-workbench/tree/master/github/dotfiles)
- [An example webcrawler from the Golang tour](https://github.com/golang/tour/blob/master/solutions/webcrawler.go)
- [Using The Log Package In Go - ArdanLabs](https://www.ardanlabs.com/blog/2013/11/using-log-package-in-go.html)
- [Go: Marshal and unmarshal JSON with time and URL data](https://ukiahsmith.com/blog/go-marshal-and-unmarshal-json-with-time-and-url-data/)
- [Stack Overflow](https://stackoverflow.com)
  - Not going to list all of the specific pages, but check out [my recent upvotes](https://stackoverflow.com/users/380599/jlucktay?tab=votes) for a sample 😅
  - Special mention to [this page](https://stackoverflow.com/questions/38362631/go-error-non-constant-array-bound) though, for detailing one issue I ran into when optimising

## Magefile target candidates

- `check latest error log for non-404 errors:`
- `$ cat $(ls -rt1 *error* | tail -n 1) | grep -v "\[404\]"`
- `$ go test -count 1 -v .`
- `$ time go run .`
- `$ time go run . -url monzo.com`
- `$ rm -f *.json *.log`
