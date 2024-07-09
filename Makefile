SHELL := /bin/bash

TAG_SHA = agora-backend-middleware

# set the build args from env file
DECONARGS = $(shell sed -e 's/^/--build-arg /' .env | tr '\n' ' ' | sed 's/ *$$//')
# Generate the Arguments using DECONARGS
GEN_ARGS = $(eval BARGS=$(DECONARGS))
# Set the SERVER_PORT from .env, run target checks if set and defaults to 8080 if needed
SERVER_PORT = $(shell grep SERVER_PORT .env | cut -d '=' -f2 | tr -d '[:space:]' || echo "8080")

# Docker and Golang source files
DOCKER_FILES := Dockerfile
GO_SOURCE_FILES := $(shell find . -type f -name '*.go')
GO_MOD_FILES := go.mod go.sum

.PHONY: all check-env build run clean

all: build run

check-env: 
	@if [ ! -f .env ]; then \
		echo ".env file not found. Please create one."; \
		exit 1;\
	fi

build_marker: $(DOCKER_FILES) $(GO_SOURCE_FILES) $(GO_MOD_FILES) .env
	@echo "Running docker build with tag: ${TAG_SHA}"
	$(GEN_ARGS)
	docker build -t $(TAG_SHA) $(BARGS) .
	@touch build_marker

build: check-env build_marker

run:
	@SERVER_PORT=$$(grep SERVER_PORT .env | cut -d '=' -f2 | tr -d '[:space:]' || echo "8080"); \
	echo "Running docker container on port: $$SERVER_PORT"; \
	docker run --env-file .env -p $$SERVER_PORT:$$SERVER_PORT $(TAG_SHA)

clean:
	@echo "Stopping and removing containers using $(TAG_SHA) image..."
	@docker ps -a -q --filter ancestor=$(TAG_SHA) | xargs -r docker stop
	@docker ps -a -q --filter ancestor=$(TAG_SHA) | xargs -r docker rm
	@echo "Removing $(TAG_SHA) image..."
	@docker rmi $(TAG_SHA) || true
	@echo "Cleanup complete."

dev: check-env
	@SERVER_PORT=$${SERVER_PORT:-8080}; \
	echo "Running in development mode on port: $$SERVER_PORT"; \
	go run ./cmd/main.go

test:
	go test ./...

test-verbose:
	go test -v ./...

test-coverage:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

test-no-cache:
	go test -count=1 ./...

test-coverage-func:
	go test -coverprofile=coverage.out ./...
	go tool cover -func=coverage.out

test-race:
	go test -race ./...

benchmark:
	go test -bench=. ./...

lint:
	golangci-lint run
