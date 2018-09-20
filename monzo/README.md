# Monzo

See the [other](DESIGN.md) [Markdown](SPEC.md) [files](NOTES.md) for some more details.

You can supply a URL as a command line argument like so:

``` shell
$ go run . -url monzo.com
Pages crawled: 1222
Pages outside target 'https://monzo.com' domain: 9465
$
```

If this argument is not supplied at runtime, it will default and fall back to my very small work-in-progress personal site. 🙂
