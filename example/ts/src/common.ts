/**
 * Asserts that the given number is a positive integer.
 * @param n The number to assert.
 * @throws {Error} If the given number is not a positive integer.
 */
export const assertPositiveInteger = (n: number): void => {
  if (!(n >>> 0 === parseFloat(n.toString()))) {
    throw new Error("Not a positive integer");
  }
};
