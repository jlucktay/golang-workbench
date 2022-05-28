FROM golang:1.17 AS builder

# Set some shell options for using pipes and such.
SHELL [ "/bin/bash", "-euo", "pipefail", "-c" ]

# Don't call any C code; the 'scratch' base image used later won't have any libraries to reference.
ENV CGO_ENABLED=0

# Use Go modules
ENV GO111MODULE=on

WORKDIR /go/src/go.jlucktay.dev/golang-workbench/jam-gp

# Add the sources.
COPY . .

# Compile! With the '--mount' flags below, Go's build cache is kept between builds.
# https://github.com/golang/go/issues/27719#issuecomment-514747274
RUN --mount=type=cache,target=/go/pkg/mod \
  --mount=type=cache,target=/root/.cache/go-build \
  go build -ldflags="-X 'go.jlucktay.dev/version.builtBy=Docker'" -trimpath -v -o /bin/jam-gp

FROM gcr.io/distroless/base:nonroot AS deployable
USER 65532

# Bring binary over.
COPY --from=builder /bin/jam-gp /bin/jam-gp

VOLUME /workdir
WORKDIR /workdir

ENTRYPOINT [ "/bin/jam-gp" ]
CMD [ "--help" ]
