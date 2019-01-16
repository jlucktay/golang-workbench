# Ref

## General

- <https://github.com/avelino/awesome-go>
- [Writing Great Go Code](https://scene-si.org/2018/07/24/writing-great-go-code/)

``` golang
mux := http.NewServeMux()
mux.Handle("/dog/", d) // trailing slash so will also handle /dog/something/else
mux.Handle("/cat", c) // no trailing slash; /cat/something will 404
```

### Sample input data

```sh
curl --silent --get "http://mockbin.org/bin/41ca3269-d8c4-4063-9fd5-f306814ff03f" --header Accept:application/json | jq
```

## Design

- [The Three Principles of Excellent API Design](https://nordicapis.com/the-three-principles-of-excellent-api-design/)
- [Using Golang to Build Microservices at The Economist: A Retrospective](https://www.infoq.com/articles/golang-the-economist)
- [Notes on API design in Go](https://xyrillian.de/thoughts/posts/golang-api-design.html)

## Patterns

- [Context](https://blog.golang.org/context)
- [Aspects of a good Go library](https://medium.com/@cep21/aspects-of-a-good-go-library-7082beabb403)
- [How I write Go HTTP services after seven years](https://medium.com/statuscode/how-i-write-go-http-services-after-seven-years-37c208122831)

## Libraries

- [Fast HTTP package for Go. Up to 10x faster than net/http.](https://github.com/valyala/fasthttp)
- Here is [a good third-party ServeMux](https://godoc.org/github.com/julienschmidt/httprouter) that allows easy access to methods for routing & path parameters.

### stdlib

- [`func CanonicalHeaderKey`](https://golang.org/pkg/net/http/#CanonicalHeaderKey)

## Internals

- [How the Go runtime implements maps efficiently (without generics)](https://dave.cheney.net/2018/05/29/how-the-go-runtime-implements-maps-efficiently-without-generics)
- [Computing and plotting π with Gonum and a zest of Monte Carlo](https://blog.gopheracademy.com/advent-2018/montecarlo/)
- [Avoiding high GC overhead with large heaps](https://blog.gopheracademy.com/advent-2018/avoid-gc-overhead-large-heaps/)

## Security

- <https://github.com/shieldfy/API-Security-Checklist>

### Auth

- [Ask HN: What do you use for authentication and authorization?](https://news.ycombinator.com/item?id=18767767)

## Stability

- [Gracefully restarting a Golang web server](https://tomaz.lovrec.eu/posts/graceful-server-restart/)

## Performance

- [go-perfbook](https://github.com/dgryski/go-perfbook)

## Database

- [MongoDB - Go Migration Guide](https://www.mongodb.com/blog/post/go-migration-guide)
- [In MySQL, never use “utf8”. Use “utf8mb4”.](https://medium.com/@adamhooper/in-mysql-never-use-utf8-use-utf8mb4-11761243e434)

## Versioning

- [Versioning your API in Go](https://dev.to/geosoft1/versioning-your-api-in-go-1g4h)

## Examples

- [Write and Deploy a Golang Web App](https://vpsranked.com/write-and-deploy-a-golang-web-app/)
- [AWS DynamoDB with Go SDK](https://github.com/aws/aws-sdk-go-v2/tree/master/example/service/dynamodb)
- [Serverless Reference Architecture: Vote Application](https://github.com/aws-samples/lambda-refarch-voteapp)
- [Serverless Reference Architecture: Web Application](https://github.com/aws-samples/lambda-refarch-webapp)
- [Serverless Golang API With AWS Lambda](https://dzone.com/articles/serverless-golang-api-with-aws-lambda)

### Idempotency

- [Cloud Functions pro tips: Building idempotent functions](https://cloud.google.com/blog/products/serverless/cloud-functions-pro-tips-building-idempotent-functions)
- [Idempotency key](https://stripe.com/blog/idempotency)

## Testing

- [Test-Driven Development in Go](https://medium.com/@pierreprinetti/test-driven-development-in-go-baeab5adb468)
- [GoConvey](http://goconvey.co)
  - <https://github.com/smartystreets/goconvey>
- [Learn Go with tests](https://github.com/quii/learn-go-with-tests)
- [mkcert: Valid HTTPS certificates for localhost](https://blog.filippo.io/mkcert-valid-https-certificates-for-localhost/)
- [When Writing Unit Tests, Don’t Use Mocks](https://sendgrid.com/blog/when-writing-unit-tests-dont-use-mocks/)
- [Unit testing and interfaces](https://blog.andreiavram.ro/golang-unit-testing-interfaces/)

### Examples

- [Testable Examples in Go](https://blog.golang.org/examples)

## Error handling

- [errorx](https://github.com/joomcode/errorx)
- [Exploring Error Handling Patterns in Go](https://8thlight.com/blog/kyle-krull/2018/08/13/exploring-error-handling-patterns-in-go.html)

## Documentation

- [swag](https://github.com/swaggo/swag)
