import * as net from "net";

const client = new net.Socket();

client.connect(22123, "127.0.0.1", () => {
  console.log("Connected to server");
  client.write("push topicA xxxx");
});

client.on("data", (data) => {
  console.log("Received: " + data);
  client.destroy();
});

client.on("close", () => {
  console.log("Connection closed");
});
