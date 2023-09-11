export const MAX_ACCESS_GRANT_DAYS = 60;

/**
 * Variables used to initialize AWS Amplify to work with our Cognito instance.
 */
export const AMPLIFY_CONFIG = {
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

export const MONITORING_CONSENT_MESSAGE = 'You are accessing an official US Government information system!  By entering this system, you acknowledge your responsibility under applicable laws and policies and agree to maintain the security and privacy of all information contained herein. You agree that you will only use this system for authorized purposes. Your activities on the system may be monitored, and unauthorized or illegal behavior may result in criminal prosecution or administrative sanction.';
