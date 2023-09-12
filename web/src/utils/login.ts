// ////////////////////////////////////////////////////////////////////////////
// AWS Imports
// ////////////////////////////////////////////////////////////////////////////
import { Amplify, Auth } from 'aws-amplify';

// ////////////////////////////////////////////////////////////////////////////
// Local Imports
// ////////////////////////////////////////////////////////////////////////////
import {
  accessToken,
  clearCurrentUser,
  loginStatus,
  setCurrentUser,
} from '../stores/current-user';
import { buildQuery } from './api';
import type { TActions } from './api';
import { AMPLIFY_CONFIG } from './constants';
import { derivePasswordHash } from './hashing';
import { tokenExpiration } from './jwt';
import { isLoggedInAsAdmin } from './auth';

// ////////////////////////////////////////////////////////////////////////////
// Admin/DoS
// ////////////////////////////////////////////////////////////////////////////

export const handleFederatedLogin = async () => {
  Amplify.configure( AMPLIFY_CONFIG );

  await Auth.federatedSignIn( {
    provider: import.meta.env.PUBLIC_COGNITO_OKTA_PROVIDER_NAME,
  } );

  // Needed due to limitations of Cognito
  loginStatus.set( 'loggedIn' );
};

/**
 * Initiates a federated Okta login through Cognito. If successful,
 * updates the current user store with the relevant data.
 */
export const handleAdminLogin = async () => {
  // Needed due to limitations of Cognito
  if ( loginStatus.get() === 'loggedOut' ) {
    return false;
  }

  let authenticated = false;

  Amplify.configure( AMPLIFY_CONFIG );

  try {
    const user = await Auth.currentAuthenticatedUser( { bypassCache: true } );

    if ( user ) {
      const payload = user?.signInUserSession?.idToken?.payload;
      const { email, exp } = payload;

      // Retrieve additional data from the application.
      const response = await buildQuery( `admin?username=${email}`, null, 'GET' );
      const { data } = await response.json();
      const { active, role, team, token } = data;

      // Add the required data from the id token to the current user store.
      if ( active ) {
        setCurrentUser( { email, team, role, exp } );
        loginStatus.set( 'loggedIn' );
        accessToken.set( token );

        authenticated = true;
      }
    }
  } catch ( err ) {
    console.error( err );
  }

  return authenticated;
};

// ////////////////////////////////////////////////////////////////////////////
// Partner
// ////////////////////////////////////////////////////////////////////////////
/**
 * Retrieves the salt value used to hash the user's password.
 * @param username The name of the user to look up.
 * @returns The salt value (if the user exits).
 */
const getUserPasswordSalt = async ( username: string ) => {
  const response = await buildQuery( 'creds/salt', { username } );
  const { data } = await response.json();

  return data;
};

/**
 * Send the locally generated password hash to the server to authenticate user and request access.
 * @param action Whether to initiate a authenticated session or confirm an existing session.
 * @param hash The locally generated password hash.
 * @param username The email of the user attempting to log in.
 */
const submitUserPasswordHash = async (
  action: TActions,
  hash: string,
  username: string,
  token: string,
): Promise<Nullable<string>> => {
  const response = await buildQuery( 'guest/auth', {
    action,
    hash,
    username,
    token,
  } );

  const { data } = await response.json();

  if ( data ) {
    const parsed = JSON.parse( data );

    return parsed.token ?? null;
  }

  return null;
};

export const handlePartnerLogin = async ( username: string, password: string, token: string ) => {
  let authenticated = false;

  try {
    const salt = await getUserPasswordSalt( username );

    if ( !salt ) { return authenticated; }

    const localHash = await derivePasswordHash( password, salt );
    const jwt = await submitUserPasswordHash( 'create', localHash, username, token );

    if ( !jwt ) {
      return authenticated;
    }
    accessToken.set( jwt );


    const exp = tokenExpiration( jwt );

    // Retrieve additional data from the application.
    const response = await buildQuery( `guest?id=${username}`, null, 'GET' );
    const { data } = await response.json();
    const { role, team } = data;

    // Add the required data from the id token to the current user store.
    setCurrentUser( { email: username, team, role, exp } );
    loginStatus.set( 'loggedIn' );

    authenticated = true;
  } catch ( err ) {
    console.error( err );
  }

  return authenticated;
};

// ////////////////////////////////////////////////////////////////////////////
// Common
// ////////////////////////////////////////////////////////////////////////////

export const logout = async () => {
  try {
    // For UX
    const isAdmin = isLoggedInAsAdmin();

    // Admin signout, though this doesn't always work
    await Auth.signOut();

    // Partner signout
    accessToken.set( '' );

    // Common signout
    clearCurrentUser();
    loginStatus.set( 'loggedOut' );
    window.location.replace( isAdmin ? '/adminLogin' : '/partnerLogin' );
  } catch ( err ) {
    console.error( 'error signing out', err );
  }
};
