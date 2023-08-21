// ////////////////////////////////////////////////////////////////////////////
// React Imports
// ////////////////////////////////////////////////////////////////////////////
import { useEffect, useState } from 'react';
import type { FC, FormEvent } from 'react';

// ////////////////////////////////////////////////////////////////////////////
// Local Imports
// ////////////////////////////////////////////////////////////////////////////
import { showError } from '../utils/alert';

// ////////////////////////////////////////////////////////////////////////////
// Styles and CSS
// ////////////////////////////////////////////////////////////////////////////
import '../styles/form.css';
import '../styles/button.module.scss'
import BackButton from './BackButton';

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
  accessGrantDate: Date;
  accessEndDate: Date;
}

interface ITeamElementProps {
  readonly teams: ITeam[];
  setData: ( val: string ) => void;
}

// ////////////////////////////////////////////////////////////////////////////
// Helpers
// ////////////////////////////////////////////////////////////////////////////
const sortTeams = ( a: ITeam, b: ITeam ) => {
  if( a.teamName > b.teamName ) {
    return 1;
  } else if( b.teamName > a.teamName ) {
    return -1;
  } else {
    return 0;
  }
}

const TeamElement: FC<ITeamElementProps> = ( { teams, setData } ) => {
  if (1 === teams.length ) {
    return <input id="family-name-input" type="text" disabled value={teams[0].teamName}/>;
  } else {
    const sorted = teams.sort(sortTeams);
    return <select id="team-input" onChange={e => setData( e.target.value )}>
      { sorted.map( ( { id, teamName } ) => <option key={id} value={id}>{teamName}</option>) }
    </select>
  }
}

// ////////////////////////////////////////////////////////////////////////////
// Interface and Implementation
// ////////////////////////////////////////////////////////////////////////////
const NewUser: FC<INewUserProps> = ( props ) => {
  const [ teams ] = useState( props.teams );
  const [ userData, setUserData ] = useState<Partial<INewUserData>>( {} );

  useEffect( () => {
    clear( [ 'given-name-input', 'family-name-input', 'email-input' ] );
  }, [] );

  const clear = ( ids: string[] ) => {
    ids.forEach( id => (document.getElementById( id ) as HTMLInputElement).value = '' );
  } 

  const handleUpdate = ( key: keyof INewUserData, value: string ) => {
    setUserData( { ...userData, [key]: value } );
  };

  const validateSubmission = () => {
    if( !userData['email']?.match( /^.+@.+$/ ) ) {
      showError("Email address is not valid");
      return false;
    }

    return true;
  }

  const handleSubmit = async ( e: FormEvent<HTMLFormElement> ) => {
    e.preventDefault();
    if( !validateSubmission() ) {
      return;
    }

    // TODO: Submit user data --> Does this need to vary b/w Bob and Sue types?
    console.log( userData );
  }

  return <>
    <form onSubmit={handleSubmit}>
      <label>
        Given (First) Name
        <input id="given-name-input" type="text" required onChange={e => handleUpdate('givenName', e.target.value)} />
      </label>
      <label>
        Family (Last) Name
        <input id="family-name-input" type="text" required onChange={e => handleUpdate('familyName', e.target.value)} />
      </label>
      <label>
        Email
        <input id="email-input" type="text" required onChange={e => handleUpdate('email', e.target.value)} />
      </label>
      <label>
        Team
        <TeamElement teams={teams} setData={val => handleUpdate('team', val)} />
      </label>
      <button id="login-btn" type="submit">Invite User</button>
    </form>
    <BackButton showConfirmDialog />
  </>;
}

export default NewUser;