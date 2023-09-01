// ////////////////////////////////////////////////////////////////////////////
// Local Imports
// ////////////////////////////////////////////////////////////////////////////
import currentUser from '../stores/current-user';

// ////////////////////////////////////////////////////////////////////////////
// Helpers
// ////////////////////////////////////////////////////////////////////////////
const userIsExpired = () => {
  const exp = currentUser.get().exp;
  if( !exp ) {
    return true;
  }

  const expTime = Number.parseInt( exp as unknown as string || '' );
  return isNaN( expTime ) || expTime < Date.now()/1000;
}

// ////////////////////////////////////////////////////////////////////////////
// Implementation
// ////////////////////////////////////////////////////////////////////////////
export const isLoggedIn = ( additionalCheck?: boolean ) => {
  return !userIsExpired() && ( additionalCheck ?? true );
}

export const userIsExternalPartner = () => {
  const role = currentUser.get().role;
  return role === 'externalPartner';
}

export const userIsAdmin = () => {
  const role = currentUser.get().role;
  return role === 'superAdmin' || role === 'admin';
}

export const userIsSuperAdmin = () => {
  const role = currentUser.get().role;
  return role === 'superAdmin';
}

export const isLoggedInAsSuperAdmin = () => {
  return isLoggedIn( userIsSuperAdmin() );
}

export const isLoggedInAsAdmin = () => {
  return isLoggedIn( userIsAdmin() );
};

export const isLoggedInAsExternalPartner = () => {
  return isLoggedIn( userIsExternalPartner() );
};

/**
 * Checks whether the current user is authenticated and if not,
 * redirects them to the specified page.
 */
const protectPage = ( protectionFn: () => boolean, redirect: string ) => {
  return () => {
    const authenticated = protectionFn();

    if ( !authenticated ) {
      console.log( redirect );
      window.location.replace( redirect.startsWith( '/' ) ? redirect : `/${redirect}` );
    }
  };
}

export const adminOnlyPage = protectPage( isLoggedInAsAdmin, 'adminLogin' );
export const superAdminOnlyPage = protectPage( isLoggedInAsAdmin, 'adminLogin' );
export const partnerOnlyPage = protectPage( isLoggedInAsAdmin, 'partnerLogin' );