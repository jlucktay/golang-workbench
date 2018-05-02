# Containerize This! How to build Golang Dockerfiles

## Single-stage

``` shell
docker build -t hello-docker-main:1.0 . -f Dockerfile.single
```

## Multi-stage

``` shell
docker build -t hello-docker-main:1.1 . -f Dockerfile.multi
```

## FROM scratch

``` shell
docker build -t hello-docker-main:1.2 . -f Dockerfile.scratch
```
