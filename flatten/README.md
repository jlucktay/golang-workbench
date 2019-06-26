# Flatten

Flatten will take an arbitrarily nested slice of ints (and interfaces) and return a flat, single-dimensional slice.

[![GoDoc](https://godoc.org/github.com/jlucktay/golang-workbench/flatten?status.svg)][badge-godoc]

## Installation

### Prerequisites

You should have a [working Go environment](https://golang.org/doc/install) and have `$GOPATH/bin` in your `$PATH`.

### Compiling

To download the source, compile, and install the demo binary, run:

``` shell
go get github.com/jlucktay/golang-workbench/flatten/...
```

The source code will be located in `$GOPATH/src/github.com/jlucktay/golang-workbench/flatten/`.

A newly-compiled `flatten` binary will be in `$GOPATH/bin/`.

## Usage

Launching the demo binary:

``` shell
$ flatten
Started with: []interface {}{[]interface {}{1, 2, []int{3}}, 4}
Finished with: []int{1, 2, 3, 4}
```

### Importing the package

To import the `flatten` package for use elsewhere, add the following to your import declarations:

``` go
import (
    "github.com/jlucktay/golang-workbench/flatten/pkg/flatten"
)
```

The `Flatten()` func can then be called:

``` go
flatten.Flatten([]int{[]int{1},2,3})
```

Further details are [available on godoc.org](https://godoc.org/github.com/jlucktay/golang-workbench/flatten/pkg/flatten).

## References and inspirations

- [The Go Programming Language Specification](https://golang.org/ref/spec)
- From the [Go wiki](https://github.com/golang/go/wiki):
  - [Interface slice](https://github.com/golang/go/wiki/InterfaceSlice)
  - [Slice tricks](https://github.com/golang/go/wiki/SliceTricks)

## License

[MIT](https://choosealicense.com/licenses/mit/)

[badge-godoc]: https://godoc.org/github.com/jlucktay/golang-workbench/flatten
