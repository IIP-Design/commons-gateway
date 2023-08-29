import { persistentMap } from '@nanostores/persistent';

export type ICurrentUser = {
  email: string
  team: string
  isAdmin: 'true' | 'false'
}

const currentUser = persistentMap<Partial<ICurrentUser>>( 'CommonsGatewayCurrentUser:',
  {},
  { listen: false } ); // Do not sync across tabs

/**
 * Removes user data from local storage.
 */
export const clearCurrentUser = () => {
  localStorage.removeItem( 'CommonsGatewayCurrentUser:email' );
  localStorage.removeItem( 'CommonsGatewayCurrentUser:team' );
  localStorage.removeItem( 'CommonsGatewayCurrentUser:isAdmin' );
};

export default currentUser;
