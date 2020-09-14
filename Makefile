PWD?=$(shell pwd)
APP?=$(shell basename ${PWD})

GOOS?=linux
GOARCH?=amd64
COMMIT?=$(shell git rev-parse --short HEAD)
BUILD_TIME?=$(shell date -u '+%Y-%m-%d_%H:%M:%S')

.PHONY: build
build:
	go build -o ${APP} -v *.go

linux_build: 
	GOOS=${GOOS} GOARCH=${GOARCH} make build
.DEFAULT_GOAL := build

.PHONY: rebuild
rebuild: clean build

.PHONY: clean
clean:
	@rm -f ${APP}