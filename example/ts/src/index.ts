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
    console.log(`Sending ${i}`);
    const res = await flemq.push("ts_tests", `Hello from TS ${i}`);
    console.log("Push Result:", res);
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
