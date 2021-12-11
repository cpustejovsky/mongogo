SHELL := /bin/bash

run:
	cd app;\
	go run main.go

update:
	go get -u -t -d -v ./...
	go mod vendor
	
# build:
# 	go build -ldflags "-X main.build=local"

# Building containers

# $(shell git rev-parse --short HEAD)
VERSION := 1.0

all: service

service:
	docker build \
		-f zarf/docker/dockerfile \
		-t service-amd64:$(VERSION) \
		--build-arg BUILD_REF=$(VERSION) \
		--build-arg BUILD_DATE=`date -u +"%Y-%m-%dT%H:%M:%SZ"` \
		.
