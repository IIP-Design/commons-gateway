/**
 * Converts any date to a string in the `YYYY-MM-DD` format
 * @param date The date to be transformed.
 */
export const getYearMonthDay = ( date: Date ) => {
  const year = date.getFullYear();
  const month = date.getMonth() + 1; // Add one to adjust for zero indexing.
  const day = date.getDate();

  return `${year}-${month > 9 ? month : `0${month}`}-${day}`;
};

/**
 * Calculates the current date plus the provided number of days.
 * Returns the result as a Date value.
 * @param days The number of days to add to the current date.
 */
export const addDaysToNow = ( days: number ) => {
  const now = new Date();
  const future = new Date( now ).setDate( now.getDate() + days );

  return new Date( future );
};
