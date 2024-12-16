include .env
export

.PHONY: test format compile clean install serve auth build cdc

.ONESHELL:

test:
	@TOKEN=$(shell cat .token) go test \
		-count=1 \
		-v ./... | sed -e "/PASS/s//$$(printf "\033[32mPASS\033[0m")/" \
			-e "/FAIL/s//$$(printf "\033[31mFAIL\033[0m")/" \
			-e "/SKIP/s//$$(printf "\033[33mSKIP\033[0m")/" \
			-e "/^=== RUN/d" \
			-e "/^\?/d"
format:
	@go fmt

compile:
	@GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build \
		-ldflags="-w -s" \
		-o app \
		main.go

clean:
	@go clean -cache

install:
	@go mod tidy

serve:
	@go run main.go serve

auth:
	@go run main.go auth -u $${TEST_USERNAME} -p $${TEST_PASSWORD} > .token

build:
	@docker build --progress=plain . -t github.com/z411392/golab

cdc:
	@go run main.go cdc