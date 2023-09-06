// ////////////////////////////////////////////////////////////////////////////
// Local Imports
// ////////////////////////////////////////////////////////////////////////////
import currentUser from '../stores/current-user';
import type { TUserRole } from '../stores/current-user';
import { buildQuery, constructUrl } from './api';

// ////////////////////////////////////////////////////////////////////////////
// Types and Interfaces
// ////////////////////////////////////////////////////////////////////////////
type TImmediateUxProtectionFn = () => boolean;
type TPermissionVerificationFn = ( redirect: string ) => Promise<void>;

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

const makeAdminVerificationFn = ( roles: TUserRole[] ): TPermissionVerificationFn => async ( redirect: string ) => {
  const email = currentUser.get().email || '';
  let authenticated = false;

  try {
    const response = await buildQuery( 'admin/get', { username: email }, 'POST' );
    const { data } = await response.json();
    const { role } = data;

    authenticated = roles.includes( role );

    // eslint-disable-next-line no-empty
  } catch ( err ) {}

  if ( !authenticated ) {
    window.location.assign( redirect );
  }
};

const partnerVerificationFn: TPermissionVerificationFn = async ( redirect: string ) => {
  const email = currentUser.get().email || '';
  let authenticated = false;

  try {
    const response = await fetch( `${constructUrl( 'guest' )}?id=${email}` );
    const { data } = await response.json();
    const { role } = data;

    authenticated = role === 'guest admin';

  // eslint-disable-next-line no-empty
  } catch ( err ) {}

  if ( !authenticated ) {
    window.location.assign( redirect );
  }
};

// ////////////////////////////////////////////////////////////////////////////
// Implementation
// ////////////////////////////////////////////////////////////////////////////
export const isLoggedIn = ( additionalCheck?: boolean ) => !userIsExpired() && ( additionalCheck ?? true );

export const userIsExternalPartner = () => {
  const { role } = currentUser.get();

  return role === 'guest admin';
};

export const userIsAdmin = () => {
  const { role } = currentUser.get();

  return role === 'super admin' || role === 'admin';
};

export const userIsSuperAdmin = () => {
  const { role } = currentUser.get();

  return role === 'super admin';
};

export const isLoggedInAsSuperAdmin = () => isLoggedIn( userIsSuperAdmin() );

export const isLoggedInAsAdmin = () => isLoggedIn( userIsAdmin() );

export const isLoggedInAsExternalPartner = () => isLoggedIn( userIsExternalPartner() );

/**
 * Checks whether the current user is authenticated and if not,
 * redirects them to the specified page.
 */
const protectPage = (
  immediateUxProtectionFn: TImmediateUxProtectionFn,
  redirect: string,
  permissionVerificationFn?: TPermissionVerificationFn,
) => () => {
  // Returns quickly to redirect well-formed users
  const authenticated = immediateUxProtectionFn();

  if ( !authenticated ) {
    window.location.replace( redirect.startsWith( '/' ) ? redirect : `/${redirect}` );
  }

  // Returns async to catch malicious users tryign to bypass normal login rules
  permissionVerificationFn && permissionVerificationFn( redirect );
};

export const adminOnlyPage = protectPage( isLoggedInAsAdmin, 'adminLogin', makeAdminVerificationFn( ['super admin', 'admin'] ) );
export const superAdminOnlyPage = protectPage( isLoggedInAsAdmin, 'adminLogin', makeAdminVerificationFn( ['super admin'] ) );
export const partnerOnlyPage = protectPage( isLoggedInAsExternalPartner, 'partnerLogin', partnerVerificationFn );
export const loggedInOnlyPage = protectPage( isLoggedIn, 'login' );
