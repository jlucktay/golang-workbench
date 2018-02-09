# Compiling

Compile this with the following, to see the Go compiler make the decision to inline the function:

```golang
go build -gcflags -m .\inlining.go
```
