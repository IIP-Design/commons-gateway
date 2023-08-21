/**
 * Determines what to display in the scroll controls range.
 * @param total Count of all users.
 * @param show Number of users to show in the table.
 * @param offset The depth of scroll into the total (by increments of the viewCount).
 * @returns The string that makes up the top of the range.
 */
const setUpperBound = ( total: number, show: number, offset: number ) => {
  const remaining = total - ( show * offset );
  const theoreticalTop = show * ( offset + 1 );

  if ( remaining === 1 ) {
    return '';
  }

  if ( total < show || remaining < show ) {
    return `-${total}`;
  }

  return `-${theoreticalTop}`;
};

/**
 * Generates the string showing the pagination status of the table.
 * @param userCount The total number of users.
 * @param viewCount The number of items shown at any one time.
 * @param viewOffset The depth of scroll into the total (by increments of the viewCount).
 * @returns The full count string.
 */
export const renderCountWidget = ( userCount: number, viewCount: number, viewOffset: number ) => {
  const start = viewOffset * viewCount + 1;
  const end = setUpperBound( userCount, viewCount, viewOffset );

  return `${start}${end} out of ${userCount}`;
};

/**
 * Calculates how many pages of results the table has.
 * @param total The total number of users.
 * @param show The number of items shown at any one time.
 */
export const setIntermediatePagination = ( total: number, show: number ) => {
  const divisions = Math.floor( total / show );
  const remainder = total % show;
  const segments = remainder > 0 ? divisions + 1 : divisions;

  const pages = [];

  for ( let i = 0; i < segments; i++ ) {
    pages.push( i + 1 );
  }

  return pages;
};
