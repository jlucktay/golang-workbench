all: lint build run

### Setup

build: build-1.0 build-1.1 build-1.2 show

build-1.0:
	docker build --tag hello-docker-main:1.0 -f 1.0-single/Dockerfile .

build-1.1:
	docker build --tag hello-docker-main:1.1 -f 1.1-multi/Dockerfile .

build-1.2:
	docker build --tag hello-docker-main:1.2 -f 1.2-scratch/Dockerfile .

show:
	@docker images hello-docker-main

### Execute

run: run-1.0 run-1.1 run-1.2

run-1.0:
	@docker run hello-docker-main:1.0

run-1.1:
	@docker run hello-docker-main:1.1

run-1.2:
	@docker run hello-docker-main:1.2

### Linting

lint: lint-1.0 lint-1.1 lint-1.2

lint-1.0:
	hadolint 1.0-single/Dockerfile

lint-1.1:
	hadolint 1.1-multi/Dockerfile

lint-1.2:
	hadolint 1.2-scratch/Dockerfile

### Cleanup

clean: clean-1.0 clean-1.1 clean-1.2

clean-1.0:
	docker image rm --force hello-docker-main:1.0

clean-1.1:
	docker image rm --force hello-docker-main:1.1

clean-1.2:
	docker image rm --force hello-docker-main:1.2
