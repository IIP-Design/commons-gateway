/**
 * Returns a random integer between 0 and the provided value.
 * @param max The upper limit (non-inclusive) of the range.
 */
export const randomInt = ( max: number ) => Math.floor( Math.random() * max );

/**
 * Returns a random string of the provided length with the first letter capitalized.
 * @param length How many characters to include in the name.
 */
export const generateName = ( length: number ) => {
  const chars = 'abcdefghijklmnopqrstuvwxyz';

  let str = '';

  for ( let i = 0; i < length; i++ ) {
    str += chars.charAt( randomInt( chars.length ) );
  }

  return str.charAt( 0 ).toUpperCase() + str.slice( 1 );
};
