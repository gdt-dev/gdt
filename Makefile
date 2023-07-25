VERSION ?= $(shell git describe --tags --always --dirty)

.PHONY: test

test:
	go test -v ./...
