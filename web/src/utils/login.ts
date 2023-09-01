// ////////////////////////////////////////////////////////////////////////////
// AWS Imports
// ////////////////////////////////////////////////////////////////////////////
import { Amplify, Auth } from 'aws-amplify';

// ////////////////////////////////////////////////////////////////////////////
// 3PP Imports
// ////////////////////////////////////////////////////////////////////////////
import Cookie from 'js-cookie';

// ////////////////////////////////////////////////////////////////////////////
// Local Imports
// ////////////////////////////////////////////////////////////////////////////
import { clearCurrentUser, loginStatus, setCurrentUser } from '../stores/current-user';
import { showError } from './alert';
import { buildQuery } from './api';
import { AMPLIFY_CONFIG } from './constants';

// ////////////////////////////////////////////////////////////////////////////
// Implementation
// ////////////////////////////////////////////////////////////////////////////

export const handleFederatedLogin = async () => {
  Amplify.configure( AMPLIFY_CONFIG );
  
  await Auth.federatedSignIn( {
    provider: import.meta.env.PUBLIC_COGNITO_OKTA_PROVIDER_NAME,
  } );

  // Needed due to limitations of Cognito
  loginStatus.set( 'loggedIn' );
}

/**
 * Initiates a federated Okta login through Cognito. If successful,
 * updates the current user store with the relevant data.
 */
export const handleAdminLogin = async () => {
  // Needed due to limitations of Cognito
  if( loginStatus.get() === 'loggedOut' ) {
    return false;
  }

  let authenticated = false;
  Amplify.configure( AMPLIFY_CONFIG );

  try {
    const user = await Auth.currentAuthenticatedUser( { bypassCache: true } );
  
    if ( user ) {
      const { payload: { email, exp } } = user?.signInUserSession?.idToken;
  
      // Retrieve additional data from the application.
      const response = await buildQuery( 'admin/get', { username: email }, 'POST' );
      const { data } = await response.json();
      const { active, team } = data;
  
      // Add the required data from the id token to the current user store.
      setCurrentUser( { email, team, role: active ? 'superAdmin' : 'admin', exp } );
      loginStatus.set( 'loggedIn' );
      authenticated = true;
    }
  } catch ( err ) {
    console.log( err );
  }

  return authenticated;
};
  
export const handlePartnerLogin = async ( email: string, password: string ) => {
  const response = await buildQuery( 'partner/login', { email, password }, 'POST' );
  const { data } = await response.json();

  if( data ) {
    const { role, team, exp, token } = data;
    Cookie.set( 'token', token );

    // Add the required data from the id token to the current user store.
    setCurrentUser( { email, team, role, exp } );
  } else {
    showError( 'Invalid username or password' );
  }
}
  
export const logout = async () => {
  try {
    // Admin signout, though this doesn't always work
    await Auth.signOut();

    // Partner signout
    Cookie.remove( 'token' );
    
    // Common signout
    clearCurrentUser();
    loginStatus.set( 'loggedOut' );
    window.location.replace( '/' );
  } catch ( err ) {
    console.log( 'error signing out', err );
  }
};