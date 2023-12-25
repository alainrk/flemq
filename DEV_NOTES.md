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
- Adding stupid push command and simple memory storage (cache + mutex)
  - No persistence, no real queue (offset based on internal counter in memory), no store per topic, nothing, just testing stuff
- Adding as stupid as the push command, a pick command to retrieve from a specific offset

## 2023-12-25

- I'd maybe like to keep the connection open until timeout or explicit termination (EXIT command)
- TODO: I need to:
  - Write back to the client the response instead of just logging
  - Implement the FLEP stuff for response (Look at RESP for inspiration on data type handling, I could really just steal what I need when I need it)
- TODO: I need to create a queue store for each topic (implicit topic creation btw for now, will be customizable in the future)
- Topic creation: for now is implicit, I'm gonna do the above point for now and improve later
