/**
 * Manipulates the given string value so that it begins with
 * a capital letter followed by all lower case letters.
 * @param str Any string value.
 * @returns A capitalized version of the provided string.
 */
export const capitalize = ( str: string ): string => {
  if ( !str || !str.length ) {
    return '';
  }

  const lower = str.toLowerCase();

  return lower.substring( 0, 1 ).toUpperCase() + lower.substring( 1, lower.length );
};

/**
 * Manipulate the given string so that it follows title casing,
 * i.e. the first letter of each word is capitalized.
 * @param str Any string value.
 * @returns A title-case version of the provided string.
 */
export const titleCase = ( str: string ) => {
  const parts = str
    ?.replace( /([A-Z])+/g, capitalize )
    ?.split( /(?=[A-Z])|[\.\-\s_]/ ) // eslint-disable-line no-useless-escape
    .map( x => x.toLowerCase() ) ?? [];

  if ( parts.length === 0 ) {
    return '';
  }

  parts[0] = capitalize( parts[0] );

  return parts.reduce( ( acc, part ) => `${acc} ${part.charAt( 0 ).toUpperCase()}${part.slice( 1 )}` );
};

/**
 * Generates a pseudorandom alphanumeric string of a given length.
 * @param length The number of desired characters in the random string.
 * @returns A random string value.
 */
export const randomString = ( length: number ) => {
  const chars = 'ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789';

  let str = '';

  for ( let i = 0; i < length; i++ ) {
    str += chars.charAt( Math.floor( Math.random() * chars.length ) );
  }

  return str;
};
