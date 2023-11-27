/* eslint-disable @typescript-eslint/no-explicit-any */

/**
 * Retrieve a subset of a given length from a list of items.
 * @param arr The list of items to split up.
 * @param count The number of items to retrieve.
 * @param offset The depth of scroll into the total (by increments of the count).
 */
export const selectSlice = ( arr: any[], count: number, offset: number ) => {
  const startingIndex = offset * count;

  return arr.slice( startingIndex, startingIndex + count );
};
