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
  // Wait a bit for the publisher to start
  await sleep(2000);

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
