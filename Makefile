.PHONY: dev
dev:
	air

.PHONY: test
test:
	@go test -v -race -coverprofile /tmp/c.out ./...

.PHONY: cert
# Generate self-signed certificate
cert:
	@mkdir -p cert
	@openssl req -x509 -nodes -days 99999 -newkey rsa:2048 -keyout cert/key.pem -out cert/cert.pem -config cert.conf -extensions 'v3_req'

.PHONY: example-plain example-tls
example-plain:
	@go run ./example/plain/*.go
example-tls:
	@go run ./example/tls/*.go
