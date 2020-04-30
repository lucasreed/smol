REDIS_NAME=redis-smolserv
DOCKER_RUN_NAME=smol
SMOL_IMAGE=smol
SMOL_TAG=latest
SMOL_DOCKER_FULL=${SMOL_IMAGE}:${SMOL_TAG}

all: lint build run

clean: build-clean redis-clean docker-clean

docker: docker-build docker-run

lint:
	-golangci-lint run ./...

run:
	./smolserv --storage redis

build:
	go build -ldflags "-s -w" -o smolserv ./cmd/smolserv/

build-clean:
	rm ./smolserv || true

redis:
	docker run --name ${REDIS_NAME} -d -p 6379:6379 redis:5.0.7-buster

redis-clean:
	docker stop ${REDIS_NAME} || true
	docker rm ${REDIS_NAME} || true

docker-build:
	docker build --cache-from ${SMOL_DOCKER_FULL} -t ${SMOL_DOCKER_FULL} .

docker-run:
	docker run -d --name ${DOCKER_RUN_NAME} -p 8080:8080 ${SMOL_DOCKER_FULL}
	docker logs ${DOCKER_RUN_NAME} -f

docker-clean:
	docker stop ${DOCKER_RUN_NAME} || true
	docker rm ${DOCKER_RUN_NAME} || true
