import { persistentAtom, persistentMap } from '@nanostores/persistent';

export type TLoginStatus = 'loggedIn' | 'loggedOut';
export type TUserRole = 'super admin' | 'admin' | 'guest admin' | 'guest';

export interface ICurrentUser {
  email: string;
  team: string;
  role: TUserRole;
  exp: number;
}

const STORAGE_KEY_PREFIX = 'CommonsGatewayCurrentUser';

const currentUser = persistentMap<Partial<ICurrentUser>>( `${STORAGE_KEY_PREFIX}:`,
  {},
  { listen: false } ); // Do not sync across tabs

// This is needed due to issues with Cognito logout: https://stackoverflow.com/questions/58154256/aws-cognito-how-to-force-select-account-when-signing-in-with-google
export const loginStatus = persistentAtom<TLoginStatus>( `${STORAGE_KEY_PREFIX}:login`, 'loggedOut' );

export const accessToken = persistentAtom<string>( `${STORAGE_KEY_PREFIX}:access`, '' );

/**
 * Removes user data from local storage.
 */
export const clearCurrentUser = () => {
  localStorage.removeItem( `${STORAGE_KEY_PREFIX}:email` );
  localStorage.removeItem( `${STORAGE_KEY_PREFIX}:team` );
  localStorage.removeItem( `${STORAGE_KEY_PREFIX}:role` );
  localStorage.removeItem( `${STORAGE_KEY_PREFIX}:exp` );
};

export const setCurrentUser = ( { email, team, role, exp }: ICurrentUser ) => {
  currentUser.setKey( 'email', email );
  currentUser.setKey( 'team', team );
  currentUser.setKey( 'role', role );
  currentUser.setKey( 'exp', exp );
};

export default currentUser;
