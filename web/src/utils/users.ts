// ////////////////////////////////////////////////////////////////////////////
// Local Imports
// ////////////////////////////////////////////////////////////////////////////
import currentUser from '../stores/current-user';
import { showError, showTernary } from './alert';
import { buildQuery } from './api';
import { escapeQueryStrings } from './string';
import type { TUserRole } from './types';

// ////////////////////////////////////////////////////////////////////////////
// Interfaces and Types
// ////////////////////////////////////////////////////////////////////////////
export interface IUserFormData {
  givenName: string;
  familyName: string;
  email: string;
  team: string;
  role: TUserRole;
}

// ////////////////////////////////////////////////////////////////////////////
// Functions
// ////////////////////////////////////////////////////////////////////////////
export const makeDummyUserForm = (): IUserFormData => ( {
  givenName: '',
  familyName: '',
  email: '',
  team: currentUser.get().team || '',
  role: 'guest' as TUserRole,
} );

export const makeApproveUserHandler = ( inviteeEmail: string ) => {
  const inviterEmail = currentUser.get().email;

  return async () => {
    const { isConfirmed, isDenied } = await showTernary(
      'By approving this user they will be allowed to upload media to the Content Commons system until deactivated or their login expires.  Denying access will allow an external partner to re-propose an invitation.',
      { confirmButtonText: 'Approve' },
    );
    let wasUpdated = false;

    if ( isConfirmed ) {
      const { ok } = await buildQuery( 'guest/approve', { inviteeEmail, inviterEmail }, 'POST' );

      if ( !ok ) {
        showError( 'Unable to accept invite' );
      } else {
        wasUpdated = true;
      }
    } else if ( isDenied ) {
      const escaped = escapeQueryStrings( inviteeEmail );
      const { ok } = await buildQuery( `guest?id=${escaped}`, null, 'DELETE' );

      if ( !ok ) {
        showError( 'Unable to reject invite' );
      } else {
        wasUpdated = true;
      }
    }

    if ( wasUpdated ) {
      window.location.reload();
    }
  };
};
