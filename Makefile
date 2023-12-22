.PHONY: dev
dev:
	air

.PHONY: test
test:
	@go test -v -race -coverprofile /tmp/c.out ./...
