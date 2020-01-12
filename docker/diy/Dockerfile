FROM golang:stretch as builder

RUN mkdir /build
COPY src /build/

WORKDIR /build
ENV CGO_ENABLED 0
ENV GOARCH amd64
ENV GOOS linux

# Adapted from here:
# https://drewdevault.com/2020/01/04/Slow.html#testgo
#
# Note: "We don't currently support PIE when not using cgo. There is some background at #18968."
# https://github.com/golang/go/issues/30986#issuecomment-475626018

RUN go build -o main -a -ldflags '-extldflags "-fno-PIC -static"' -tags 'osusergo netgo static_build'

FROM scratch

COPY --from=builder /build/main /app/

WORKDIR /app
ENV VERSION v1.2

CMD [ "./main" ]

# TODO
# Make some use of this:
# https://cloud.google.com/appengine/docs/flexible/custom-runtimes/build
