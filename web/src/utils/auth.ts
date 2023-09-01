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

export const userIsAdmin = () => {
  const role = currentUser.get().role;
  return role === 'superAdmin' || role === 'admin';
}

export const isLoggedInAsSuperAdmin = () => {
  const superAdmin = currentUser.get().role === 'superAdmin';
  return isLoggedIn( superAdmin );
}

export const isLoggedInAsAdmin = () => {
  return isLoggedIn( userIsAdmin() );
};

/**
 * Checks whether the current user is authenticated and if not,
 * redirects them to the admin login page.
 */
export const adminOnlyPage = () => {
  const authenticated = isLoggedInAsAdmin();

  if ( !authenticated ) {
    window.location.replace( '/admin' );
  }
};
