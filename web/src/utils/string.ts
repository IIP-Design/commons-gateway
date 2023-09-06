export const capitalize = ( str: string ): string => {
  if ( !str || !str.length ) {
    return '';
  }

  const lower = str.toLowerCase();

  return lower.substring( 0, 1 ).toUpperCase() + lower.substring( 1, lower.length );
};

export const titleCase = ( str: string ) => {
  const parts
      = str
        ?.replace( /([A-Z])+/g, capitalize )
        // eslint-disable-next-line no-useless-escape
        ?.split( /(?=[A-Z])|[\.\-\s_]/ )
        .map( x => x.toLowerCase() ) ?? [];

  if ( parts.length === 0 ) {
    return '';
  }

  parts[0] = capitalize( parts[0] );

  return parts.reduce( ( acc, part ) => `${acc} ${part.charAt( 0 ).toUpperCase()}${part.slice( 1 )}` );
};
