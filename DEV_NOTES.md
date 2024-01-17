# Development Notes

## 2023-12-22

- I want to use a simple text-based protocol over TCP, the idea is to start with some sort of message queue, but I don't want to limit to that
- Creating simple example that sends plain text to test during dev
- Creating simple TCP server accepting connection on 22123 port
  - Handling sigterm termination
  - Tests with netcat suffer from a stuck issue (>= 1024 chars blocks the write and need connection to be closed), not happening using example
  - Handling different types of buffering on the server, it doesn't seem to affect the nc issue though, learning sth anyway
- Implementing TLS on the server based on constant

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
- I need to:
  - Write back to the client the response instead of just logging
  - Implement the FLEP stuff for response (Look at RESP for inspiration on data type handling, I could really just steal what I need when I need it)
- I need to create a queue store for each topic (implicit topic creation btw for now, will be customizable in the future)
- Topic creation: for now is implicit, I'm gonna do the above point for now and improve later

## 2023-12-28

- Adding a subscribe command to the protocol:
  - it should behave leave the connection open
  - no timeout for the connection
  - client specifies the topic and the offset to start from (included, i.e. 0 from the beginning)
- For now I'm gonna add a stupid sleep/retry loop that polls the queue for new messages and send them as they come
- Hiding the store type behind topic interface, I may change my mind later, but for now I don't need to expose it
- I'd like to have some sort of broadcast/fan-out mechanism to distribute messages coming to a topic
  - e.g. Each topic once created (or at startup) have a goroutine that keeps track of each subscriber (a channel?) and sends the message to each of them as they come
- Creating a simple broker to be used internally, so that every topic can be subscribed and implements a publish over it at each write (let's see, I want to experiment this way)

## 2023-12-31

- Being dumb I didn't remember that unblocking channel send drops the message, so broker sends were broker for the current tests (only noticeable with more messages coming)
  - I could either enforce a buffered channel or accept the best-effort/use a config to specify the behaviour depending on the case
  - I'll start with blocking send for now
- TODO: In many places I'm going with io.ReadAll instead of io.Copy, consider to switch to that in the future if necessary
- Was using a reader in the topic writer multiple times, having the store not saving data anymore, using tee instead

## 2023-01-02

- Just for fun opening the possibility to have a non blocking broker as well, this will imply other issues anyway
  - TODO: Test non-blocking broker
  - TODO: Address non-blocking issues (also, there could be a best-effort mode of the topic, where you just stream the messages and don't care if there are lost messages)
    - Buffered channels on the client side could be a way to go, but what's the size?
- I'm going to centralize the response creation itself in the flep module, it doesn't send it, just compose it.
  - TODO: For big reponses (if any), should I think about streaming the response instead? How to deal with this?

## 2023-01-03

- I need to refactor some stuff to create proper interfaces and allow tests and better decoupling

## 2023-01-05

- DONE: Handle connection termination (e.g. client disconnects) during subscription handler (timeout is disabled in this case)

## 2023-01-11

- DONE: Fix EOF error when using fqueue store on subscribe (it gets returned as last element of subscribe topic), troubleshoot why it happens
- TODO: When using fqueue store, I need to restore the existing topics, otherwise weird stuff happens when pushing to a topic name that was already used

## 2023-01-13

- I've decided to don't implement any cache for now in the fqueue store, thinking about it I have the broker that is already providing all the subscriber with the current written message. For most use cases previously written messages are not needed as frequently or at least not in the immediate future once writter so it doesn't make sense to keep them in memory. I would change my mind if needed later.
- DONE: (fqueue persistence) I've partially fixed the re-open issue, but it only works if at least someone pushes before subscribing/picking as the broker isn't initialized yet in the topic (this because only happens implicitly at the first write (New() call)), so the broker must start at the startup or at each read I have to check if the topic exists, otherwise create it.

## 2023-01-14

- I do really need to improve some config passing around, so I'll add the store setting and passing it down to the topic level to set persistence and enable testing
- I'll also implement the above topic restore as it's not possible to read before writing when fqueue is used otherwise
