
# Image URL to use all building/pushing
DOCKER_IMAGE ?= herochinese/go-bumblebee-locust
DOCKER_TAG ?= $(shell git rev-parse --abbrev-ref HEAD)
# Docker Hub username / password
DOCKER_PASS =
DOCKER_USER =


#generate:
#	go run github.com/99designs/gqlgen generate ./...

# Build the docker image
docker-build:
	docker build . -t ${DOCKER_IMAGE}:${DOCKER_TAG}

# Push the docker image
docker-push:
	# echo ${DOCKER_PASS} | docker login --username ${DOCKER_USER} --password-stdin
	docker push ${DOCKER_IMAGE}:${DOCKER_TAG}

