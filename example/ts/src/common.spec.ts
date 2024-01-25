import tap from "tap";

import { assertPositiveInteger } from "./common";

tap.test("asserting positive integers", (t) => {
  const inputs = [0, 1, 42, 100, 666];
  t.plan(inputs.length);
  try {
    for (const i of inputs) {
      assertPositiveInteger(i);
      t.pass();
    }
  } catch (err) {
    t.fail("Should not throw error:", err);
  }
  t.end();
});

tap.test("asserting non-positive integers", (t) => {
  const inputs = [-666, -1, 0.6, 0.42, 123.456];
  t.plan(inputs.length);
  for (const i of inputs) {
    try {
      assertPositiveInteger(i);
      t.fail("Should throw error");
    } catch (err) {
      t.pass();
    }
  }
  t.end();
});
