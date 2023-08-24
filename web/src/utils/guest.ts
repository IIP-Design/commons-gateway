/**
 * Compare the guest's invite expiration date with the current date.
 * If the expiration is past the current date the user is active.
 * @param exp The date when the user's invite expires.
 */
export const isGuestActive = ( exp: string ) => new Date( exp ).getTime() > new Date().getTime();
