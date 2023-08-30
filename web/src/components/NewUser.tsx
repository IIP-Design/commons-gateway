import { useEffect, useState } from 'react';
import type { FC, FormEvent } from 'react';
import { isAfter, isBefore, addDays, parse } from 'date-fns';

import BackButton from './BackButton';

import { buildQuery } from '../utils/api';
import { addDaysToNow, getYearMonthDay } from '../utils/dates';
import currentUser from '../stores/current-user';
import { showError } from '../utils/alert';
import { MAX_ACCESS_GRANT_DAYS } from '../utils/constants';

import '../styles/form.scss';
import styles from '../styles/button.module.scss';

// ////////////////////////////////////////////////////////////////////////////
// Interfaces and Types
// ////////////////////////////////////////////////////////////////////////////
interface INewUserProps {
  readonly isAdmin: boolean;
}

interface INewUserData {
  givenName: string;
  familyName: string;
  email: string;
  team: number;
  accessEndDate: string;
}

interface ITeamElementProps {
  readonly teams: ITeam[];
  readonly setData: ( val: string ) => void;
}

// ////////////////////////////////////////////////////////////////////////////
// Helpers
// ////////////////////////////////////////////////////////////////////////////
const sortTeams = ( a: ITeam, b: ITeam ) => {
  if ( a.name > b.name ) {
    return 1;
  } if ( b.name > a.name ) {
    return -1;
  }

  return 0;
};

/**
 * Render the team selection field. If there is only one team available, it will
 * return a readonly text input. Otherwise, it will provide a select element.
 * @param param.teams A list of possible teams.
 * @param param.setData Function to set the form values on input.
 */
const TeamElement: FC<ITeamElementProps> = ( { teams, setData } ) => {
  if ( teams.length === 1 ) {
    return <input id="family-name-input" type="text" disabled value={ teams[0].name } />;
  }

  const sorted = teams.sort( sortTeams );

  return (
    <select id="team-input" onChange={ e => setData( e.target.value ) }>
      { sorted.map( ( { id, name } ) => <option key={ id } value={ id }>{ name }</option> ) }
    </select>
  );
};

const dateSelectionIsValid = ( dateStr?: string ) => {
  const now = new Date();
  const date = parse( dateStr || '', 'yyyy-MM-dd', new Date() );

  return date && isAfter( date, now ) && isBefore( date, addDays( now, MAX_ACCESS_GRANT_DAYS ) );
};

// ////////////////////////////////////////////////////////////////////////////
// Interface and Implementation
// ////////////////////////////////////////////////////////////////////////////
const NewUser: FC<INewUserProps> = ( { isAdmin } ) => {
  const [teamList, setTeamList] = useState( [] );
  const [userData, setUserData] = useState<Partial<INewUserData>>( {} );

  const clear = ( ids: string[] ) => {
    ids.forEach( id => {
      const el = document.getElementById( id ) as HTMLInputElement;

      el.value = '';
    } );
  };

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
    clear( [
      'given-name-input', 'family-name-input', 'email-input', 'date-input',
    ] );
  }, [] );

  const handleUpdate = ( key: keyof INewUserData, value?: string|Date ) => {
    setUserData( { ...userData, [key]: value } );
  };

  const validateSubmission = () => {
    if ( !userData.email?.match( /^.+@.+$/ ) ) {
      showError( 'Email address is not valid' );

      return false;
    } if ( !dateSelectionIsValid( userData.accessEndDate ) ) {
      showError( `Please select an access grant end date after today and no more than ${MAX_ACCESS_GRANT_DAYS} in the future` );

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
        team: currentUser.get().team,
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
          <input id="given-name-input" type="text" required onChange={ e => handleUpdate( 'givenName', e.target.value ) } />
        </label>
        <label>
          <span>Family (Last) Name</span>
          <input id="family-name-input" type="text" required onChange={ e => handleUpdate( 'familyName', e.target.value ) } />
        </label>
      </div>
      <div className="field-group">
        <label>
          <span>Email</span>
          <input id="email-input" type="text" required onChange={ e => handleUpdate( 'email', e.target.value ) } />
        </label>
        <label>
          <span>Team</span>
          <TeamElement teams={ teamList } setData={ val => handleUpdate( 'team', val ) } />
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

export default NewUser;
