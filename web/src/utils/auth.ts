import { Amplify, Auth } from 'aws-amplify';

import currentUser, { clearCurrentUser } from '../stores/current-user';

interface IIdToken {
  /** The user's email address. */
  email: string
  /** The expiration time for the token. */
  exp: number
  /** The issued-at time when Cognito issued the token. */
  iat: number
  /** The identity provider that issued the token. */
  iss: string
  /** A unique identifier (UUID), or subject, for the authenticated user. */
  sub: string
}

/**
 * Variables used to initialize AWS Amplify to work with our Cognito instance.
 */
const awsConfig = {
  aws_project_region: 'us-east-1',
  aws_cognito_region: 'us-east-1',
  aws_cognito_identity_pool_id: import.meta.env.PUBLIC_COGNITO_IDENTITY_POOL_ID,
  aws_user_pools_id: import.meta.env.PUBLIC_COGNITO_USER_POOLS_ID,
  aws_user_pools_web_client_id: import.meta.env.PUBLIC_COGNITO_USER_POOL_WEB_CLIENT_ID,
  oauth: {
    domain: import.meta.env.PUBLIC_COGNITO_CLIENT_DOMAIN,
    scope: [
      'email', 'openid', 'profile',
    ],
    redirectSignIn: import.meta.env.PUBLIC_COGNITO_CLIENT_REDIRECT_SIGNIN,
    redirectSignOut: import.meta.env.PUBLIC_COGNITO_CLIENT_REDIRECT_SIGNOUT,
    responseType: 'code',
  },
  federationTarget: 'COGNITO_USER_POOLS',
};

/**
 * Initializes AWS Amplify for use in authentication.
 */
export const initializeAmplify = () => Amplify.configure( awsConfig );

/**
 * Pulls out the relevant data from the Cognito id token.
 * This data is user to provide the app with the relevant user information.
 * @param payload
 */
const setUserFromToken = ( payload: IIdToken ) => {
  currentUser.setKey( 'email', payload.email );
};

/**
 * Checks whether the current user is authenticated using Cognito/Okta.
 */
export const isLoggedIn = async () => {
  try {
    const user = await Auth.currentAuthenticatedUser();

    if ( user ) {
      // Add the required data from the id token to the current user store.
      setUserFromToken( user?.signInUserSession?.idToken?.payload );

      return true;
    }
  } catch ( err ) {
    console.log( err );
  }

  return false;
};

/**
 * Initiates a federated Okta login through Cognito. If successful,
 * updates the current user store with the relevant data.
 */
export const handleAdminLogin = async () => {
  Auth.federatedSignIn( {
    provider: import.meta.env.PUBLIC_COGNITO_OKTA_PROVIDER_NAME,
  } );
};

export const signOut = () => {
  try {
    Auth.signOut();
    clearCurrentUser();
  } catch ( err ) {
    console.log( 'error signing out', err );
  }
};

/**
 * Checks whether the current user is authenticated and if not,
 * redirects them to the admin login page.
 */
export const adminOnlyPage = async () => {
  const authenticated = await isLoggedIn();

  if ( !authenticated ) {
    window.location.replace( '/admin' );
  }
};
