# FleMQ
TEST
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

## Node.js client example

```ts
import { FlemQ } from "../src/flemq";

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
    const res = await flemq.push("ts_tests", `Pushing message ${i}`);
    console.log("Result:", res);
    await sleep(1500);
  }
})();

// Subscriber
(async () => {
  const flemq = new FlemQ({
    port: 22123,
    serder: "base64",
  });

  await flemq.connect();
  await flemq.subscribe("ts_tests", (data: string): void => {
    console.log("Received:", data);
  });
})();

// Pick
(async () => {
  const flemq = new FlemQ({
    port: 22123,
    serder: "base64",
  });

  await flemq.connect();
  const res = await flemq.pick("ts_tests", 1000);
  console.log("Pick Result:", res);
})();
```
