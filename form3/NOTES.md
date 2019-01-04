# Ref

## General

- <https://github.com/avelino/awesome-go>

``` golang
mux := http.NewServeMux()
mux.Handle("/dog/", d) // trailing slash so will also handle /dog/something/else
mux.Handle("/cat", c) // no trailing slash; /cat/something will 404
```

## Design

- [The Three Principles of Excellent API Design](https://nordicapis.com/the-three-principles-of-excellent-api-design/)
- [Using Golang to Build Microservices at The Economist: A Retrospective](https://www.infoq.com/articles/golang-the-economist)

## Patterns

- [Context](https://blog.golang.org/context)
- [Aspects of a good Go library](https://medium.com/@cep21/aspects-of-a-good-go-library-7082beabb403)
- [How I write Go HTTP services after seven years](https://medium.com/statuscode/how-i-write-go-http-services-after-seven-years-37c208122831)

## Libraries

- [Fast HTTP package for Go. Up to 10x faster than net/http.](https://github.com/valyala/fasthttp)
- Here is [a good third-party ServeMux](https://godoc.org/github.com/julienschmidt/httprouter) that allows easy access to methods for routing & path parameters.

### stdlib

- [`func CanonicalHeaderKey`](https://golang.org/pkg/net/http/#CanonicalHeaderKey)

## Security

- <https://github.com/shieldfy/API-Security-Checklist>

### Auth

- [Ask HN: What do you use for authentication and authorization?](https://news.ycombinator.com/item?id=18767767)

## Versioning

- [Versioning your API in Go](https://dev.to/geosoft1/versioning-your-api-in-go-1g4h)

## Examples

- [Write and Deploy a Golang Web App](https://vpsranked.com/write-and-deploy-a-golang-web-app/)
- [AWS DynamoDB with Go SDK](https://github.com/aws/aws-sdk-go-v2/tree/master/example/service/dynamodb)
- [Serverless Reference Architecture: Vote Application](https://github.com/aws-samples/lambda-refarch-voteapp)
- [Serverless Reference Architecture: Web Application](https://github.com/aws-samples/lambda-refarch-webapp)

### Idempotency

- [Cloud Functions pro tips: Building idempotent functions](https://cloud.google.com/blog/products/serverless/cloud-functions-pro-tips-building-idempotent-functions)
- [Idempotency key](https://stripe.com/blog/idempotency)

## Testing

- [Test-Driven Development in Go](https://medium.com/@pierreprinetti/test-driven-development-in-go-baeab5adb468)
- [GoConvey](http://goconvey.co)
  - <https://github.com/smartystreets/goconvey>
- [Learn Go with tests](https://github.com/quii/learn-go-with-tests)

## Documentation

- [swag](https://github.com/swaggo/swag)
