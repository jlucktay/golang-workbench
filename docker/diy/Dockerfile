FROM golang:stretch as builder

RUN mkdir /build
COPY src /build/

WORKDIR /build
ENV CGO_ENABLED 0
ENV GOOS linux

RUN go build -o main -a -installsuffix cgo -ldflags '-extldflags "-static"' .
FROM scratch

COPY --from=builder /build/main /app/

WORKDIR /app
ENV VERSION v1.2

CMD [ "./main" ]

# TODO
# Make some use of this:
# https://cloud.google.com/appengine/docs/flexible/custom-runtimes/build
