// ////////////////////////////////////////////////////////////////////////////
// Local Imports
// ////////////////////////////////////////////////////////////////////////////
import currentUser from '../stores/current-user';

// ////////////////////////////////////////////////////////////////////////////
// Helpers
// ////////////////////////////////////////////////////////////////////////////
const userIsExpired = () => {
  const { exp } = currentUser.get();

  if ( !exp ) {
    return true;
  }

  const expTime = Number.parseInt( exp as unknown as string || '', 10 );

  return Number.isNaN( expTime ) || expTime < Date.now() / 1000;
};

// ////////////////////////////////////////////////////////////////////////////
// Implementation
// ////////////////////////////////////////////////////////////////////////////
export const isLoggedIn = ( additionalCheck?: boolean ) => !userIsExpired() && ( additionalCheck ?? true );

export const userIsExternalPartner = () => {
  const { role } = currentUser.get();

  return role === 'externalPartner';
};

export const userIsAdmin = () => {
  const { role } = currentUser.get();

  return role === 'superAdmin' || role === 'admin';
};

export const userIsSuperAdmin = () => {
  const { role } = currentUser.get();

  return role === 'superAdmin';
};

export const isLoggedInAsSuperAdmin = () => isLoggedIn( userIsSuperAdmin() );

export const isLoggedInAsAdmin = () => isLoggedIn( userIsAdmin() );

export const isLoggedInAsExternalPartner = () => isLoggedIn( userIsExternalPartner() );

/**
 * Checks whether the current user is authenticated and if not,
 * redirects them to the specified page.
 */
const protectPage = ( protectionFn: () => boolean, redirect: string ) => () => {
  const authenticated = protectionFn();

  if ( !authenticated ) {
    window.location.replace( redirect.startsWith( '/' ) ? redirect : `/${redirect}` );
  }
};

export const adminOnlyPage = protectPage( isLoggedInAsAdmin, 'adminLogin' );
export const superAdminOnlyPage = protectPage( isLoggedInAsAdmin, 'adminLogin' );
export const partnerOnlyPage = protectPage( isLoggedInAsAdmin, 'partnerLogin' );
