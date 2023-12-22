.PHONY: dev
dev:
	air

.PHONY: example
example:
	@go run example/main.go

.PHONY: test
test:
	@go test -v -race -coverprofile /tmp/c.out ./...
