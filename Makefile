CONTAINER_REGISTRY ?= localhost:5001
PROJECT ?= generic-dns-controller
TAG ?= 0.0.1
REPO ?= middlewaregruppen
SHELL ?= /usr/bin/bash
.ONESHELL:
MODULE   = $(shell env GO111MODULE=on go list -m)
BIN      = $(CURDIR)/bin
BUILDPATH ?= $(BIN)/$(shell basename $(MODULE))

.PHONY: \
	Makefile \
	test \
	binary

binary:
	@go build \
		-tags release \
		-ldflags '-X main.VERSION=${VERSION} -X main.COMMIT=${COMMIT} -X main.BRANCH=${BRANCH} -X main.GOVERSION=${GOVERSION}' \
		-o $(BUILDPATH) $(CURDIR)/main.go

build:
	@docker build --network=host --tag ${CONTAINER_REGISTRY}/${REPO}/${PROJECT}:${TAG} . && \
	docker push ${CONTAINER_REGISTRY}/${REPO}/${PROJECT}:${TAG}

test: 
	@go test -v ./dns
