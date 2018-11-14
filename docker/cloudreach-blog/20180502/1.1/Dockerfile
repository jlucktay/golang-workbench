FROM golang:alpine as builder
RUN mkdir /build
COPY src /build/
WORKDIR /build
RUN go build -o main .
# hadolint ignore=DL3007
FROM alpine:latest
RUN adduser -S -D -H -h /app appuser
USER appuser
COPY --from=builder /build/main /app/
WORKDIR /app
ENV VERSION v1.1
CMD ["./main"]
