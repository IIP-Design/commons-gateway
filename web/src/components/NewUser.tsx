// ////////////////////////////////////////////////////////////////////////////
// React Imports
// ////////////////////////////////////////////////////////////////////////////
import { useEffect, useState } from 'react';
import type { FC, FormEvent } from 'react';

// ////////////////////////////////////////////////////////////////////////////
// 3PP Imports
// ////////////////////////////////////////////////////////////////////////////
import { isAfter, isBefore, addDays, parse } from 'date-fns';

// ////////////////////////////////////////////////////////////////////////////
// Local Imports
// ////////////////////////////////////////////////////////////////////////////
import BackButton from './BackButton';
import { showError } from '../utils/alert';
import { MAX_ACCESS_GRANT_DAYS } from '../utils/constants';

// ////////////////////////////////////////////////////////////////////////////
// Styles and CSS
// ////////////////////////////////////////////////////////////////////////////

import '../styles/form.css';
import styles from '../styles/button.module.scss';

// ////////////////////////////////////////////////////////////////////////////
// Interfaces and Types
// ////////////////////////////////////////////////////////////////////////////
interface INewUserProps {
  readonly teams: ITeam[];
}

interface INewUserData {
  givenName: string;
  familyName: string;
  email: string;
  team: number;
  accessEndDate?: string;
}

interface ITeamElementProps {
  readonly teams: ITeam[];
  readonly setData: ( val: string ) => void;
}

// ////////////////////////////////////////////////////////////////////////////
// Helpers
// ////////////////////////////////////////////////////////////////////////////
const sortTeams = ( a: ITeam, b: ITeam ) => {
  if ( a.teamName > b.teamName ) {
    return 1;
  } if ( b.teamName > a.teamName ) {
    return -1;
  }

  return 0;
};

const TeamElement: FC<ITeamElementProps> = ( { teams, setData } ) => {
  if ( teams.length === 1 ) {
    return <input id="family-name-input" type="text" disabled value={ teams[0].teamName } />;
  }
  const sorted = teams.sort( sortTeams );

  return (
    <select id="team-input" onChange={ e => setData( e.target.value ) }>
      { sorted.map( ( { id, teamName } ) => <option key={ id } value={ id }>{ teamName }</option> ) }
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
const NewUser: FC<INewUserProps> = ( { teams } ) => {
  const [teamList, setTeamList] = useState( teams ); // eslint-disable-line @typescript-eslint/no-unused-vars
  const [userData, setUserData] = useState<Partial<INewUserData>>( {} );

  const clear = ( ids: string[] ) => {
    ids.forEach( id => {
      const el = document.getElementById( id ) as HTMLInputElement;

      el.value = '';
    } );
  };

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

    // TODO: Submit user data --> Does this need to vary b/w Bob and Sue types?
    console.log( userData );
  };

  return (
    <form onSubmit={ handleSubmit }>
      <label>
        <span>Given (First) Name</span>
        <input id="given-name-input" type="text" required onChange={ e => handleUpdate( 'givenName', e.target.value ) } />
      </label>
      <label>
        <span>Family (Last) Name</span>
        <input id="family-name-input" type="text" required onChange={ e => handleUpdate( 'familyName', e.target.value ) } />
      </label>
      <label>
        <span>Email</span>
        <input id="email-input" type="text" required onChange={ e => handleUpdate( 'email', e.target.value ) } />
      </label>
      <label>
        <span>Team</span>
        <TeamElement teams={ teamList } setData={ val => handleUpdate( 'team', val ) } />
      </label>
      <label>
        <span>Access End Date</span>
        <input id="date-input" type="date" onChange={ e => handleUpdate( 'accessEndDate', e.target.value ) } />
      </label>
      <div>
        <button id="login-btn" type="submit" className={ styles.btn }>Invite User</button>
        <BackButton showConfirmDialog />
      </div>
    </form>
  );
};

export default NewUser;
