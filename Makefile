REDIS_NAME ?= redis-smolserv
DOCKER_RUN_NAME ?= smol
SMOL_IMAGE ?= smol
SMOL_TAG ?= latest
SMOL_DOCKER_FULL ?= ${SMOL_IMAGE}:${SMOL_TAG}
SMOL_STORAGE ?= boltdb

all: lint test build

run: build serve

clean: build-clean redis-clean boltdb-clean

docker: docker-build docker-run

lint:
	-golangci-lint run ./...

serve:
	./smolserv --storage ${SMOL_STORAGE}

build:
	go build -ldflags "-s -w" -o smolserv ./cmd/smolserv/

test:
	go test ./...

redis:
	docker run --name ${REDIS_NAME} -d -p 6379:6379 redis:5.0.7-buster

docker-build:
	docker build --cache-from ${SMOL_DOCKER_FULL} -t ${SMOL_DOCKER_FULL} .

docker-run:
	docker run --rm --name ${DOCKER_RUN_NAME} -p 8080:8080 ${SMOL_DOCKER_FULL}

# Clean up tasks
build-clean:
	-@rm ./smolserv 2>/dev/null || true

redis-clean:
	-@docker stop ${REDIS_NAME} 2>/dev/null || true
	-@docker rm ${REDIS_NAME} 2>/dev/null || true

boltdb-clean:
	-@rm ./boltdb 2>/dev/null || true
