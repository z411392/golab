include .env
export

.PHONY: test format build clean install

test:
	@cd lib && go test \
		-count=1 \
		-v ./... | sed -e "/PASS/s//$$(printf "\033[32mPASS\033[0m")/" \
	    	-e "/FAIL/s//$$(printf "\033[31mFAIL\033[0m")/" \
			-e "/SKIP/s//$$(printf "\033[33mSKIP\033[0m")/" \
			-e "/^=== RUN/d" \
			-e "/^\?/d"
format:
	@cd lib && go fmt

build:
	@GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build \
		-ldflags="-w -s" \
		-o bin/main \
		lib/main.go

clean:
	@cd lib && go clean -cache

install:
	@cd lib && go mod tidy