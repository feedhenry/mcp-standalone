DOCKERHOST = docker.io
DOCKERORG = feedhenry
USER=$(shell id -u)
PWS=$(shell pwd)
TAG=latest

build_and_push: apb_build docker_push

.PHONY: apb_build
apb_build:
	docker run --rm -u $(USER) -v $(PWD):/mnt:z feedhenry/apb prepare
	docker build -t $(DOCKERHOST)/$(DOCKERORG)/android-app-apb:$(TAG) .

.PHONY: docker_push
docker_push:
	docker push $(DOCKERHOST)/$(DOCKERORG)/android-app-apb:$(TAG)
