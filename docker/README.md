# Containerize This! How to build Golang Dockerfiles

Adapted from [this blog post](https://www.cloudreach.com/blog/containerize-this-golang-dockerfiles/).

## Building

Run the below `docker build` command lines from within each of the respective sub-directories.

### Single-stage

``` shell
cd 1.0-single
docker build -t hello-docker-main:1.0 .
```

### Multi-stage

``` shell
cd 1.1-multi
docker build -t hello-docker-main:1.1 .
```

### FROM scratch

``` shell
cd 1.2-scratch
docker build -t hello-docker-main:1.2 .
```

## Checking

The `docker images` command will show all built images.

``` shell
$ docker images
REPOSITORY          TAG   IMAGE ID       CREATED          SIZE
hello-docker-main   1.0   f878922f1271   13 minutes ago   378MB
hello-docker-main   1.1   379235791891   8 minutes ago    6.16MB
hello-docker-main   1.2   1840671fa652   7 minutes ago    2.01MB
```

## Running

Run these `docker run` command lines from anywhere to execute the built images.

### 1.0

``` shell
$ docker run -it hello-docker-main:1.0
Hello Docker v1.0!
```

### 1.1

``` shell
$ docker run -it hello-docker-main:1.1
Hello Docker v1.1!
```

### 1.2

``` shell
$ docker run -it hello-docker-main:1.2
Hello Docker v1.2!
```
