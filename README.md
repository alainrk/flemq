# Flexible Message Queue

![Example](https://github.com/alainrk/seraph/raw/main/assets/flemq.gif)

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

## Quick and dirty usage

```sh
nc localhost 22123
> PUSH topic1 message1
> PICK topic1 0
> EXIT
```
