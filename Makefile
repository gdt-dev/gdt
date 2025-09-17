VERSION ?= $(shell git describe --tags --always --dirty)

.PHONY: test clean build run

bin/gdt:
	@cd cmd/gdt && go build -o ../../bin/gdt main.go && cd ../../

# If the first argument is "run"...
ifeq (run,$(firstword $(MAKECMDGOALS)))
  # use the rest as arguments for "run"
  RUN_ARGS := $(wordlist 2,$(words $(MAKECMDGOALS)),$(MAKECMDGOALS))
  # ...and turn them into do-nothing targets
  $(eval $(RUN_ARGS):;@:)
endif

build: clean bin/gdt

run: build
	@bin/gdt $(RUN_ARGS)

test:
	@go test -cover -v ./...

clean:
	@rm -f bin/gdt
