import { useEffect, useState } from 'react';
import type { FC, FormEvent } from 'react';

import BackButton from './BackButton';

import { buildQuery } from '../utils/api';
import { addDaysToNow, dateSelectionIsValid, getYearMonthDay } from '../utils/dates';
import currentUser from '../stores/current-user';
import { showError } from '../utils/alert';
import { MAX_ACCESS_GRANT_DAYS } from '../utils/constants';

import '../styles/form.scss';
import styles from '../styles/button.module.scss';

// ////////////////////////////////////////////////////////////////////////////
// Interfaces and Types
// ////////////////////////////////////////////////////////////////////////////
interface IUserFormData {
  givenName: string;
  familyName: string;
  email: string;
  team: string;
  accessEndDate: string;
}

interface IUserFormProps {
  readonly user?: boolean;
}

const initialState = {
  givenName: '',
  familyName: '',
  email: '',
  team: '',
  accessEndDate: getYearMonthDay( new Date() ),
};

// ////////////////////////////////////////////////////////////////////////////
// Interface and Implementation
// ////////////////////////////////////////////////////////////////////////////
const UserForm: FC<IUserFormProps> = ( { user } ) => {
  const [isAdmin, setIsAdmin] = useState( false );
  const [teamList, setTeamList] = useState( [] );
  const [userData, setUserData] = useState<IUserFormData>( initialState );

  // Reset the form fields
  const clear = () => {
    setUserData( initialState );
  };

  // Check whether the user is an admin and set that value in state.
  // Doing so outside of a useEffect hook causes a mismatch in values
  // between the statically rendered portion and the client.
  useEffect( () => {
    setIsAdmin( currentUser.get().isAdmin === 'true' );
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

  // Clear the input fields.
  useEffect( () => {
    const getUser = async ( id: string ) => {
      const response = await buildQuery( `guest?id=${id}`, null, 'GET' );
      const { data } = await response.json();


      if ( data ) {
        setUserData( {
          ...data,
          accessEndDate: getYearMonthDay( new Date( data.expiration ) ),
        } );
      }
    };

    clear();

    if ( user ) {
      const urlSearchParams = new URLSearchParams( window.location.search );
      const { id } = Object.fromEntries( urlSearchParams.entries() );

      getUser( id );
    }
  }, [user] );

  const handleUpdate = ( key: keyof IUserFormData, value?: string|Date ) => {
    setUserData( { ...userData, [key]: value } );
  };

  /**
   * Ensure that the form submissions are valid before sending data to the API.
   */
  const validateSubmission = () => {
    if ( !userData.email?.match( /^.+@.+$/ ) ) {
      showError( 'Email address is not valid' );

      return false;
    } if ( !dateSelectionIsValid( userData.accessEndDate ) ) {
      showError( `Please select an access grant end date after today and no more than ${MAX_ACCESS_GRANT_DAYS} in the future` );

      return false;
    } if ( isAdmin && !userData.team ) {
      // Admin users have the option to set a team, so a team should be
      showError( 'Please assign this user to a valid team' );

      return false;
    }

    return true;
  };

  const handleSubmit = async ( e: FormEvent<HTMLFormElement> ) => {
    e.preventDefault();

    if ( !validateSubmission() ) {
      return;
    }

    const invitation = {
      inviter: currentUser.get().email,
      invitee: {
        email: userData.email,
        givenName: userData.givenName,
        familyName: userData.familyName,
        team: userData.team || currentUser.get().team,
      },
      expiration: new Date( userData.accessEndDate as string ).toISOString(), // Conversion to iso required by Lambda
    };

    await buildQuery( 'creds/provision', invitation, 'POST' )
      .then( () => window.location.assign( '/' ) )
      .catch( err => console.error( err ) );
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
      </div>
      <div style={ { textAlign: 'center' } }>
        <button id="login-btn" type="submit" className={ styles.btn }>Invite User</button>
        <BackButton showConfirmDialog />
      </div>
    </form>
  );
};

export default UserForm;
