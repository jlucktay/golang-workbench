# Compiling

Compile this with the following, to see the Go compiler make the decision to inline the function:

```go
go build -gcflags -m .\inlining.go
```
