# Development Notes

## 2023-12-22

- I want to use a simple text-based protocol over TCP, the idea is to start with some sort of message queue, but I don't want to limit to that
- Creating simple example that sends plain text to test during dev
- Creating simple TCP server accepting connection on 22123 port
  - Handling sigterm termination
  - Tests with netcat suffer from a stuck issue (>= 1024 chars blocks the write and need connection to be closed), not happening using example
  - Handling different types of buffering on the server, it doesn't seem to affect the nc issue though, learning sth anyway
- Implementing TLS on the server based on constant (TODO: Move all that stuff to a config file / env vars)

## 2023-12-23

- I'll start creating some structs to handle server/connection stuff
- Refactoring a bit the config stuff as it started to get messy
- I'm considering adopting a simplified version of redis [RESP](https://redis.io/docs/reference/protocol-spec/) protocol for the communication, it's simple and already implemented in many languages
  - I don't think I need the full resp for now, but it seems a good approach to start with something simple like that

## 2023-12-24

- Writing a simple FLEP (flemq procotol) parser with a basic command
