# FleMQ

**Flexible Message Queue**

<p align="center">
  <img src='assets/logo.png' width='400'>
</p>

## Dev

```sh
# Generate TLS certificate/key for local testing (if you want to use TLS)
# You have to set SAN to localhost in the config (see cert.conf)
make cert

# Run the server in live reload
make dev

# Build and Run the server
make build
make run

# Send messages to the server (plaintext - tcp)
go run example/tcp/*

# Send messages to the server (tls)
go run example/producer-tls/*
```

## Quick and dirty usage

It uses a redis-like protocol, completely text based, so you can use netcat to test it.

```sh
nc localhost 22123

> push topic_x message
:0

> push topic_x message_2
:1

> pick topic_y 0
-topic topic_y does not exist

> pick topic_x 1
+message_2

> subscribe topic_x
+message
+message_2
...
```

![Example](https://github.com/alainrk/flemq/raw/main/assets/flemq.gif)
