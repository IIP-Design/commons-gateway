import { isAfter, isBefore, addDays, parse, differenceInDays } from 'date-fns';

import { MAX_ACCESS_GRANT_DAYS } from './constants';

/**
 * Converts any date to a string in the `YYYY-MM-DD` format
 * @param date The date to be transformed.
 */
export const getYearMonthDay = ( date: Date ) => {
  const year = date.getFullYear();
  const month = date.getMonth() + 1; // Add one to adjust for zero indexing.
  const day = date.getDate();

  const fmtMonth = month > 9 ? month : `0${month}`;
  const fmtDay = day > 9 ? day : `0${day}`;

  return `${year}-${fmtMonth}-${fmtDay}`;
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

/**
 * Checks whether a given date falls after the present and before max grant days.
 * @param dateStr A provided date string.
 */
export const dateSelectionIsValid = ( dateStr?: string ) => {
  const now = new Date();
  const date = parse( dateStr || '', 'yyyy-MM-dd', new Date() );

  return date && isAfter( date, now ) && isBefore( date, addDays( now, MAX_ACCESS_GRANT_DAYS ) );
};

export const userWillNeedNewPassword = (
  grantDateStr: string,
  endDateStr: string,
  expired: boolean,
  passwordResetLastAuth: boolean,
) => {
  if ( expired || !passwordResetLastAuth ) {
    return true;
  }

  const now = new Date();
  const grantDate = parse( grantDateStr, 'yyyy-MM-dd', now );
  const endDate = parse( endDateStr, 'yyyy-MM-dd', now );

  return isAfter( endDate, addDays( grantDate, MAX_ACCESS_GRANT_DAYS ) );
};

export const daysUntil = ( dateStr: string ) => {
  const dt = new Date( dateStr );

  return differenceInDays( dt, new Date() );
};
