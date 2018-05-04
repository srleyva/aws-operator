BINARY = aws-operator
VERSION=0.0.1
COMMIT=$(shell git rev-parse HEAD)
BRANCH=$(shell git rev-parse --abbrev-ref HEAD)
GITHUB_USERNAME=srleyva
GITHUB_REPO=${GITHUB_USERNAME}/${BINARY}

LDFLAGS = -ldflags "-X main.VERSION=${VERSION} -X main.COMMIT=${COMMIT} -X main.BRANCH=${BRANCH}"
pkgs = $(shell go list ./... | grep -v /vendor/ | grep -v /test/)
gobuild = go build ${LDFLAGS} -o ${BINARY}



all: format build-linux

format:
	go fmt $(pkgs)

test:
	go test `go list ./... | grep -v apis | grep -v client`

build:
	$(gobuild) main.go

build-linux:
	CGO_ENABLED=0 GOOS=linux $(gobuild) main.go

docker-build:
	docker build . -t ${GITHUB_REPO}:${COMMIT}
