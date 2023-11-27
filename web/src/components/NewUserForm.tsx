// ////////////////////////////////////////////////////////////////////////////
// React Imports
// ////////////////////////////////////////////////////////////////////////////
import { useEffect, useState } from 'react';
import type { FC, FormEvent } from 'react';

// ////////////////////////////////////////////////////////////////////////////
// 3PP Imports
// ////////////////////////////////////////////////////////////////////////////
import { addDays } from 'date-fns';

// ////////////////////////////////////////////////////////////////////////////
// Local Imports
// ////////////////////////////////////////////////////////////////////////////
import BackButton from './BackButton';

import currentUser from '../stores/current-user';
import { showError } from '../utils/alert';
import { buildQuery } from '../utils/api';
import { userIsAdmin } from '../utils/auth';
import { MAX_ACCESS_GRANT_DAYS } from '../utils/constants';
import { addDaysToNow, dateSelectionIsValid, getYearMonthDay } from '../utils/dates';
import { makeDummyUserForm } from '../utils/users';
import type { IUserFormData } from '../utils/users';

// ////////////////////////////////////////////////////////////////////////////
// Styles and CSS
// ////////////////////////////////////////////////////////////////////////////
import '../styles/form.scss';
import styles from '../styles/button.module.scss';

// ////////////////////////////////////////////////////////////////////////////
// Types and Interfaces
// ////////////////////////////////////////////////////////////////////////////
export interface INewUserFormData extends IUserFormData {
  accessEndDate: string;
}

const initialState = {
  ...makeDummyUserForm(),
  accessEndDate: getYearMonthDay( addDays( new Date(), 14 ) ),
};

// ////////////////////////////////////////////////////////////////////////////
// Interface and Implementation
// ////////////////////////////////////////////////////////////////////////////
const UserForm: FC = () => {
  const [isAdmin, setIsAdmin] = useState( false );
  const [teamList, setTeamList] = useState( [] );
  const [userData, setUserData] = useState<INewUserFormData>( initialState );

  const partnerRoles = [{ name: 'External Partner', value: 'guest' }, { name: 'External Team Lead', value: 'guest admin' }];

  // Check whether the user is an admin and set that value in state.
  // Doing so outside of a useEffect hook causes a mismatch in values
  // between the statically rendered portion and the client.
  useEffect( () => {
    setIsAdmin( userIsAdmin() );
  }, [] );

  // Generate the teams list.
  useEffect( () => {
    const getTeams = async () => {
      const response = await buildQuery( 'teams', null, 'GET' );
      const { data } = await response.json();

      if ( data ) {
        // If the user is not an admin user, only allow them to invite users to their own team.
        const filtered = isAdmin ? data : data.filter( ( t: ITeam ) => t.id === currentUser.get().team );

        setTeamList( filtered );
      }
    };

    getTeams();
  }, [isAdmin] );

  /**
   * Updates the user state on changed to the form inputs.
   * @param key The user property being updated.
   */
  const handleUpdate = ( key: keyof INewUserFormData, value?: string|Date ) => {
    setUserData( { ...userData, [key]: value } );
  };

  /**
   * Ensure that the form submissions are valid before sending data to the API.
   */
  const validateSubmission = () => {
    if ( !userData.email?.match( /^.+@.+$/ ) ) {
      showError( 'Email address is not valid' );

      return false;
    }

    if ( !dateSelectionIsValid( userData.accessEndDate ) ) {
      showError( `Please select an access grant end date after today and no more than ${MAX_ACCESS_GRANT_DAYS} in the future` );

      return false;
    }

    if ( isAdmin && !userData.team ) {
      // Admin users have the option to set a team, so a team should be
      showError( 'Please assign this user to a valid team' );

      return false;
    }

    return true;
  };

  /**
   * Validates the form inputs and sends the provided data to the API.
   * @param e The submission event - simply used to prevent the default page refresh.
   */
  const handleSubmit = async ( e: FormEvent<HTMLFormElement> ) => {
    e.preventDefault();

    if ( !validateSubmission() ) {
      return;
    }

    // Conversion to iso required by Lambda
    const expiration = new Date( userData.accessEndDate ).toISOString();

    const invitee = {
      email: userData.email.trim(),
      givenName: userData.givenName.trim(),
      familyName: userData.familyName.trim(),
      team: userData.team || currentUser.get().team,
      role: userData.role,
    };

    if ( isAdmin ) {
      const invitation = {
        inviter: currentUser.get().email,
        invitee,
        expiration,
      };

      buildQuery( 'creds/provision', invitation, 'POST' )
        .then( () => window.location.assign( '/' ) )
        .catch( err => console.error( err ) );
    } else {
      const invitation = {
        proposer: currentUser.get().email,
        invitee,
        expiration,
      };

      try {
        await buildQuery( 'creds/propose', invitation, 'POST' );
        window.location.assign( '/uploader-users' );
      } catch ( err ) {
        console.error( err );
      }
    }
  };

  return (
    <form onSubmit={ handleSubmit }>
      <div className="field-group">
        <label>
          <span>Given (First) Name</span>
          <input
            id="given-name-input"
            type="text"
            required
            value={ userData.givenName }
            onChange={ e => handleUpdate( 'givenName', e.target.value ) }
          />
        </label>
        <label>
          <span>Family (Last) Name</span>
          <input
            id="family-name-input"
            type="text"
            required
            value={ userData.familyName }
            onChange={ e => handleUpdate( 'familyName', e.target.value ) }
          />
        </label>
      </div>
      <div className="field-group">
        <label>
          <span>Email</span>
          <input
            id="email-input"
            type="text"
            required
            value={ userData.email }
            onChange={ e => handleUpdate( 'email', e.target.value ) }
          />
        </label>
        <label>
          <span>Team</span>
          <select
            id="team-input"
            disabled={ !isAdmin }
            required
            value={ userData.team }
            onChange={ e => handleUpdate( 'team', e.target.value ) }
          >
            <option value="">- Select a team -</option>
            { teamList.map( ( { id, name } ) => <option key={ id } value={ id }>{ name }</option> ) }
          </select>
        </label>
      </div>
      <div className="field-group">
        <label>
          <span>Access End Date</span>
          <input
            id="date-input"
            type="date"
            min={ getYearMonthDay( new Date() ) }
            max={ getYearMonthDay( addDaysToNow( 60 ) ) }
            value={ userData.accessEndDate }
            onChange={ e => handleUpdate( 'accessEndDate', e.target.value ) }
          />
        </label>
        { isAdmin && (
          <label>
            <span>User Role</span>
            <select
              id="role-input"
              required
              value={ userData.role }
              onChange={ e => handleUpdate( 'role', e.target.value ) }
            >
              { partnerRoles.map( ( { name, value } ) => <option key={ value } value={ value }>{ name }</option> ) }
            </select>
          </label>
        ) }
      </div>
      <div style={ { textAlign: 'center' } }>
        <button
          className={ `${styles.btn} ${styles['spaced-btn']}` }
          id="update-btn"
          type="submit"
        >
          { isAdmin ? 'Invite' : 'Propose' }
        </button>
        <BackButton text="Cancel" showConfirmDialog />
      </div>
    </form>
  );
};

export default UserForm;
