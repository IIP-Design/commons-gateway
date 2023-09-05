import { persistentAtom, persistentMap } from '@nanostores/persistent';

export type LoginStatus = 'loggedIn' | 'loggedOut';
export type UserRole = 'super admin' | 'admin' | 'guest admin' | 'guest';

export type ICurrentUser = {
  email: string;
  team: string;
  role: UserRole;
  exp: number;
}

const STORAGE_KEY_PREFIX = 'CommonsGatewayCurrentUser:';

const currentUser = persistentMap<Partial<ICurrentUser>>( STORAGE_KEY_PREFIX,
  {},
  { listen: false } ); // Do not sync across tabs

// This is needed due to issues with Cognito logout: https://stackoverflow.com/questions/58154256/aws-cognito-how-to-force-select-account-when-signing-in-with-google
export const loginStatus = persistentAtom<LoginStatus>( STORAGE_KEY_PREFIX, 'loggedOut' );

/**
 * Removes user data from local storage.
 */
export const clearCurrentUser = () => {
  localStorage.removeItem( 'CommonsGatewayCurrentUser:email' );
  localStorage.removeItem( 'CommonsGatewayCurrentUser:team' );
  localStorage.removeItem( 'CommonsGatewayCurrentUser:role' );
  localStorage.removeItem( 'CommonsGatewayCurrentUser:exp' );
};

export const setCurrentUser = ( { email, team, role, exp }: ICurrentUser ) => {
  currentUser.setKey( 'email', email );
  currentUser.setKey( 'team', team );
  currentUser.setKey( 'role', role );
  currentUser.setKey( 'exp', exp );
};

export default currentUser;
