FROM golang:alpine as builder
RUN mkdir /build
COPY src /build/
WORKDIR /build
RUN go build -o main .
FROM alpine:3.8
RUN adduser -S -D -H -h /app appuser
USER appuser
COPY --from=builder /build/main /app/
WORKDIR /app
CMD ["./main"]
