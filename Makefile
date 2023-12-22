.PHONY: dev
dev:
	air

.PHONY: example
example:
	@go run example/main.go

.PHONY: test
test:
	@go test -v -race -coverprofile /tmp/c.out ./...

.PHONY: cert
# Generate self-signed certificate
cert:
	@mkdir -p cert
	@openssl req -newkey rsa:2048 -new -nodes -x509 -days 3650 -keyout cert/key.pem -out cert/cert.pem
