# Ref

## General

- <https://github.com/avelino/awesome-go>
- [Writing Great Go Code](https://scene-si.org/2018/07/24/writing-great-go-code/)
- [Practical Go: Real world advice for writing maintainable Go programs](https://dave.cheney.net/practical-go/presentations/qcon-china.html)
- [Learning Resources for Go](https://www.stuartellis.name/articles/golang-learning-resources/)

``` go
mux := http.NewServeMux()
mux.Handle("/dog/", d) // trailing slash so will also handle /dog/something/else
mux.Handle("/cat", c) // no trailing slash; /cat/something will 404
```

### Sample input data

``` shell
curl --silent --get "http://mockbin.org/bin/41ca3269-d8c4-4063-9fd5-f306814ff03f" --header Accept:application/json | jq
```

## Design

- [The Three Principles of Excellent API Design](https://nordicapis.com/the-three-principles-of-excellent-api-design/)
- [Using Golang to Build Microservices at The Economist: A Retrospective](https://www.infoq.com/articles/golang-the-economist)
- [Notes on API design in Go](https://xyrillian.de/thoughts/posts/golang-api-design.html)
- [Build A Go API](https://medium.com/commonbond-engineering/build-a-go-api-eb27e6663d78)

## Patterns

- [Context](https://blog.golang.org/context)
  - [Please explain Go context to me like I'm a five year old what are the main benefits?](https://www.reddit.com/r/golang/comments/afuh8f/please_explain_go_context_to_me_like_im_a_five/)
  - [Context Deadlines and How to Set Them](https://engineering.grab.com/context-deadlines-and-how-to-set-them)
- [Aspects of a good Go library](https://medium.com/@cep21/aspects-of-a-good-go-library-7082beabb403)
- [How I write Go HTTP services after seven years](https://medium.com/statuscode/how-i-write-go-http-services-after-seven-years-37c208122831)
- [How to not use an http-router in go](https://blog.merovius.de/2017/06/18/how-not-to-use-an-http-router.html)

## Libraries

- [Fast HTTP package for Go. Up to 10x faster than net/http.](https://github.com/valyala/fasthttp)
- Here is [a good third-party ServeMux](https://godoc.org/github.com/julienschmidt/httprouter) that allows easy access to methods for routing & path parameters.

### stdlib

- [`func CanonicalHeaderKey`](https://golang.org/pkg/net/http/#CanonicalHeaderKey)

## Internals

- [How the Go runtime implements maps efficiently (without generics)](https://dave.cheney.net/2018/05/29/how-the-go-runtime-implements-maps-efficiently-without-generics)
- [Computing and plotting π with Gonum and a zest of Monte Carlo](https://blog.gopheracademy.com/advent-2018/montecarlo/)
- [A visual guide to Go Memory Allocator from scratch (Golang)](https://blog.learngoprogramming.com/a-visual-guide-to-golang-memory-allocator-from-ground-up-e132258453ed)

### Garbage Collection

- [Avoiding high GC overhead with large heaps](https://blog.gopheracademy.com/advent-2018/avoid-gc-overhead-large-heaps/)
- [Getting to Go: The Journey of Go's Garbage Collector](https://blog.golang.org/ismmkeynote)
- [Golang’s Real-time GC in Theory and Practice](https://making.pusher.com/golangs-real-time-gc-in-theory-and-practice/)

### Interfaces

- [Go Pointers: Why I Use Interfaces (in Go)](https://medium.com/@kent.rancourt/go-pointers-why-i-use-interfaces-in-go-338ae0bdc9e4)
- [Go interfaces and type assertions](https://marcofranssen.nl/go-interfaces-and-type-assertions/)

### Concurrency

- [On concurrency in Go HTTP servers](https://eli.thegreenplace.net/2019/on-concurrency-in-go-http-servers/)

## Security

- <https://github.com/shieldfy/API-Security-Checklist>
- [How to Hash and Verify Passwords With Argon2 in Go](https://www.alexedwards.net/blog/how-to-hash-and-verify-passwords-with-argon2-in-go)
- [Automate your Go TLS certificates with Vault PKI](https://talks.godoc.org/github.com/johanbrandhorst/presentations/certify/certify.slide)

### Auth

- [Ask HN: What do you use for authentication and authorization?](https://news.ycombinator.com/item?id=18767767)

## Stability

- [Gracefully restarting a Golang web server](https://tomaz.lovrec.eu/posts/graceful-server-restart/)

## Performance

- [go-perfbook](https://github.com/dgryski/go-perfbook)

### Caching

- [Build a Go Cache in 10 Minutes](https://hackernoon.com/build-a-go-cache-in-10-minutes-c908a8255568)

### Profiling

- [Fixing Service Performance with PProf in Go](https://nicholasjackson.io/2019/03/01/fixing-bugs-with-pprof/)

## Database

- [MongoDB - Go Migration Guide](https://www.mongodb.com/blog/post/go-migration-guide)
- [In MySQL, never use “utf8”. Use “utf8mb4”.](https://medium.com/@adamhooper/in-mysql-never-use-utf8-use-utf8mb4-11761243e434)
- [DynamoDB, explained.](https://www.dynamodbguide.com)

## Logging

- [Logrus is a structured logger for Go (golang), completely API compatible with the standard library logger](https://github.com/Sirupsen/logrus)

## Versioning

- [Versioning your API in Go](https://dev.to/geosoft1/versioning-your-api-in-go-1g4h)

## Examples

- [Write and Deploy a Golang Web App](https://vpsranked.com/write-and-deploy-a-golang-web-app/)
- [AWS DynamoDB with Go SDK](https://github.com/aws/aws-sdk-go-v2/tree/master/example/service/dynamodb)
- [Serverless Reference Architecture: Vote Application](https://github.com/aws-samples/lambda-refarch-voteapp)
- [Serverless Reference Architecture: Web Application](https://github.com/aws-samples/lambda-refarch-webapp)
- [Serverless Golang API With AWS Lambda](https://dzone.com/articles/serverless-golang-api-with-aws-lambda)
- [Ultimate Service](https://github.com/ardanlabs/service)
- [Web Service Architecture for Golang Developers](https://boobo94.xyz/web-service/webservice-architecture-golang/)
- [🔥 Idiomatic Golang Restful Starter Kit](https://github.com/ribice/gorsk)
- [Droplets is a platform for Gophers](https://github.com/spy16/droplets)
  - start from [graceful.Server](https://github.com/spy16/droplets/blob/master/pkg/graceful/graceful.go)
- [WTF Dial - A tour of building a silly app in Go](https://medium.com/wtf-dial)

- ~~[Developing A RESTful API With Golang And A MongoDB NoSQL Database](https://www.thepolyglotdeveloper.com/2019/02/developing-restful-api-golang-mongodb-nosql-database/)~~

### Idempotency

- [Cloud Functions pro tips: Building idempotent functions](https://cloud.google.com/blog/products/serverless/cloud-functions-pro-tips-building-idempotent-functions)
- [Idempotency key](https://stripe.com/blog/idempotency)
- [Avoiding Double Payments in a Distributed Payments System](https://medium.com/airbnb-engineering/avoiding-double-payments-in-a-distributed-payments-system-2981f6b070bb)

## Frameworks

- [Building a Web App With Go, Gin and React](https://hakaselogs.me/2018-04-20/building-a-web-app-with-go-gin-and-react/)

## Modules

- [This is my 2019 development setup for Go](https://medium.com/@pettersenkim/a-2019-example-of-a-go-development-environment-with-go-modules-f4a4c188b3c2)

## Testing

- [Test-Driven Development in Go](https://medium.com/@pierreprinetti/test-driven-development-in-go-baeab5adb468)
- [GoConvey](http://goconvey.co)
  - <https://github.com/smartystreets/goconvey>
- [mockify](https://github.com/brianmoran/mockify)
  - [Testing and mocking - what clicked for me](https://javorszky.co.uk/2019/02/09/testing-and-mocking-what-clicked-for-me/)
- [Learn Go with tests](https://github.com/quii/learn-go-with-tests)
  - [Learn Go with tests - Context](https://dev.to/quii/learn-go-with-tests---context-mi)
- [mkcert: Valid HTTPS certificates for localhost](https://blog.filippo.io/mkcert-valid-https-certificates-for-localhost/)
- [When Writing Unit Tests, Don’t Use Mocks](https://sendgrid.com/blog/when-writing-unit-tests-dont-use-mocks/)
- [Unit testing and interfaces](https://blog.andreiavram.ro/golang-unit-testing-interfaces/)

### Testing Examples

- [Testable Examples in Go](https://blog.golang.org/examples)
- [Adhere to this table!](https://www.restapitutorial.com/lessons/httpmethods.html)

## Error handling

- [errorx](https://github.com/joomcode/errorx)
- [Exploring Error Handling Patterns in Go](https://8thlight.com/blog/kyle-krull/2018/08/13/exploring-error-handling-patterns-in-go.html)
- [Return nil or custom error in Go](https://stackoverflow.com/questions/46022517/return-nil-or-custom-error-in-go/46022840)

## Tools

- [JSON-to-Go - Convert JSON to Go struct](https://mholt.github.io/json-to-go/)

### Containerisation

- [Containerize a golang app using private GitHub repos as import packages](https://medium.com/@anuragdhingra/containerize-a-golang-app-using-private-github-repos-as-import-packages-51e993afe394)

## Documentation

- <https://goswagger.io/>
- [swag](https://github.com/swaggo/swag)
  - [[Golang] – How to generate swagger API doc from the source code](https://dev4devs.com/2019/02/08/golang-how-to-generate-swagger-api-doc-from-the-source-code/)

## Notes

- If your API uses POST to create a resource, be sure to include a Location header in the response that includes the URL of the newly-created resource, along with a 201 status code — that is part of the HTTP standard.
