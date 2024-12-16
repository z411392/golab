include .env
export

.PHONY: test

test:
	@go test -v ./...