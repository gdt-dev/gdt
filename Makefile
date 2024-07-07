VERSION ?= $(shell git describe --tags --always --dirty)

.PHONY: test

test:
	go test -timeout 10s -v ./...
