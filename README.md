# Flexible Message Queue

## Dev

```sh
# Generate TLS certificate/key for local testing
# It has SAN to localhost in the config (see cert.conf)
make cert

# Run the server in live reload
make run

# Send messages to the server (plaintext - tcp)
go run example/tcp/*

# Send messages to the server (tls)
go run example/producer-tls/*
```
