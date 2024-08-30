# Docker Tag Monitor (`dtm`)

## Purpose

The overall goal is to prune older versioned images/tags, and keep only the leading edge.

Specifically, for any given repository/image:

- If there is a `latest`, keep it
- If there is, e.g. a `1` _and_ a `1.2`, keep both, even if they do not point to the same hash
- If there is, e.g. `1` _and_ `1.2` _and_ `1.2.3`, keep all, as above
- Remove everything else that doesn't qualify

## Usage

If you are using a non-standard container engine (CE) host, i.e. anything that isn't `unix:///var/run/docker.sock`,
then the `DOCKER_HOST` environment variable needs to be set when running `dtm`:

```shell
export DOCKER_HOST="$(docker context inspect --format="{{ .Endpoints.docker.Host }}")"
dtm
```
