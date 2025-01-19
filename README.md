# FleMQ

![](https://img.shields.io/badge/-Work%20in%20Progress-f00?style=for-the-badge)

**Flexible Message Queue**

FleMQ is a message queue system implemented in Go, drawing inspiration from Kafka's offset-based message persistence and Redis's simplicity in protocol design. At its core, it implements a custom protocol called FLEP (FleMQ Protocol), which borrows it design from [Redis RESP](https://redis.io/docs/latest/develop/reference/protocol-spec/) maintaining a simple approach for message queue operations.

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

It uses **FLEP** (FleMQ Protocol), a simple text-based command protocol, so you can use `netcat` to test it.
Commands are _case-insensitive_. So, for example, you can use `PUSH` and `push` interchangeably.

```
# Push a message into a topic
PUSH <topic> <message>

# Pick a message from a topic at offset
PICK <topic> <offset>

# Subscribe to a topic (optionally starting from an offset)
SUBSCRIBE <topic> [<start_offset>]

# Quit connection to the server
EXIT
```

Session example:

```sh
nc localhost 22123

> PUSH topic_x message
:0

> PUSH topic_x message_2
:1

> PICK topic_y 0
-topic topic_y does not exist

> PICK topic_x 1
+message_2

> SUBSCRIBE topic_x
+message
+message_2
...
```

![Example](https://github.com/alainrk/flemq/raw/main/assets/flemq.gif)

## Node.js client example

```ts
import { FlemQ, FlepResponse } from "../src/flemq";

const sleep = async (msec: number) => {
  return new Promise((resolve) => {
    setTimeout(resolve, msec);
  });
};

// Publisher
(async () => {
  const flemq = new FlemQ({
    port: 22123,
    serder: "base64",
  });

  await flemq.connect();
  for (let i = 0; i < 100; i++) {
    flemq.push("ts_tests", `Pushing message ${i}`, (data: FlepResponse) => {
      console.log("Received from push:", data);
    });
    await sleep(1500);
  }
})();

// Subscriber
(async () => {
  // Wait a bit for the publisher to start
  await sleep(2000);

  const flemq = new FlemQ({
    port: 22123,
    serder: "base64",
  });

  await flemq.connect();

  flemq.subscribe("ts_tests", (data: FlepResponse) => {
    console.log("Received from subscribe:", data);
  });
})();

// Pick
(async () => {
  const flemq = new FlemQ({
    port: 22123,
    serder: "base64",
  });

  await flemq.connect();
  flemq.pick("ts_tests", 3, (data: FlepResponse) => {
    console.log("Received from pick:", data);
  });
})();
```
