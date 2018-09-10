FROM golang:alpine as builder
RUN mkdir /build
ADD src /build/
WORKDIR /build
ENV CGO_ENABLED 0
ENV GOOS linux
RUN go build -o main -a -installsuffix cgo -ldflags '-extldflags "-static"' .
FROM scratch
COPY --from=builder /build/main /app/
WORKDIR /app
CMD ["./main"]
