import { FlemQ, FlepResponse } from "../src/flemq";
import { promisify } from "util";

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
