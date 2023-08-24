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

export default awsConfig;
