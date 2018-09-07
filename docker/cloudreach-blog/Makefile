all: build run

### Setup

build: build-1.0 build-1.1 build-1.2 show

build-1.0:
	cd 1.0-single && docker build --tag hello-docker-main:1.0 .

build-1.1:
	cd 1.1-multi && docker build --tag hello-docker-main:1.1 .

build-1.2:
	cd 1.2-scratch && docker build --tag hello-docker-main:1.2 .

show:
	docker images hello-docker-main

### Execute

run: run-1.0 run-1.1 run-1.2

run-1.0:
	docker run hello-docker-main:1.0

run-1.1:
	docker run hello-docker-main:1.1

run-1.2:
	docker run hello-docker-main:1.2

### Cleanup

clean: clean-1.0 clean-1.1 clean-1.2

clean-1.0:
	docker image rm --force hello-docker-main:1.0

clean-1.1:
	docker image rm --force hello-docker-main:1.1

clean-1.2:
	docker image rm --force hello-docker-main:1.2
