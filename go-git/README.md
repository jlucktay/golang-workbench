# `go-git`

## Purpose

Take the following `git` alias and make it faster with Go:

```bash
lines=!f() { git ls-files -- '*.go' ':!bindata.go' ':!vendor/**' | while read f; do git --no-pager blame -w -M -C -C --line-porcelain "$f" | grep '^author '; done | sort -f | uniq -ci | sort -nr; }; cd -- ${GIT_PREFIX:-.} && f
```

## Pipe breakdown

### File list (`git ls-files`)

`git ls-files` is taking a few glob arguments in the alias, so this section should handle the same.

### Blame with porcelain (`git blame --line-porcelain`)

Pulling out the author for each line, basically.

### `sort | uniq | sort`

Collating authors and line counts.

### `cd -- ${GIT_PREFIX:-.}`

Default to running from current directory if no other directory is specified.
