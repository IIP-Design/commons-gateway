// ////////////////////////////////////////////////////////////////////////////
// React Imports
// ////////////////////////////////////////////////////////////////////////////
import { useEffect, useState } from 'react';
import type { FC, FormEvent } from 'react';

// ////////////////////////////////////////////////////////////////////////////
// Local Imports
// ////////////////////////////////////////////////////////////////////////////
import BackButton from './BackButton';

import type { TUserRole } from '../utils/types';
import { showConfirm, showError } from '../utils/alert';
import { buildQuery } from '../utils/api';
import { userIsAdmin } from '../utils/auth';

// ////////////////////////////////////////////////////////////////////////////
// Styles and CSS
// ////////////////////////////////////////////////////////////////////////////
import '../styles/form.scss';
import styles from '../styles/button.module.scss';

// ////////////////////////////////////////////////////////////////////////////
// Interfaces and Types
// ////////////////////////////////////////////////////////////////////////////
interface IAdminFormData {
  givenName: string;
  familyName: string;
  email: string;
  team: string;
  role: TUserRole;
  active: boolean;
}

interface IAdminFormProps {
  readonly admin?: boolean;
}

const initialState = {
  active: true,
  givenName: '',
  familyName: '',
  email: '',
  team: '',
  role: 'admin' as TUserRole,
};

// ////////////////////////////////////////////////////////////////////////////
// Interface and Implementation
// ////////////////////////////////////////////////////////////////////////////
const AdminForm: FC<IAdminFormProps> = ( { admin } ) => {
  const [isAdmin, setIsAdmin] = useState( false );
  const [teamList, setTeamList] = useState( [] );
  const [adminData, setAdminData] = useState<IAdminFormData>( initialState );

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
        setTeamList( data );
      }
    };

    getTeams();
  }, [] );

  // Initialize the form.
  useEffect( () => {
    const getAdmin = async ( username: string ) => {
      const response = await buildQuery( `admin?username=${username}`, null, 'GET' );
      const { data } = await response.json();


      if ( data ) {
        setAdminData( {
          ...data,
          active: data.active === 'true', // value comes in as a string from lambda
        } );
      }
    };

    if ( admin ) {
      const urlSearchParams = new URLSearchParams( window.location.search );
      const { id } = Object.fromEntries( urlSearchParams.entries() );

      getAdmin( id );
    }
  }, [admin] );

  /**
   * Updates the user state on changed to the form inputs.
   * @param key The user property being updated.
   */
  const handleUpdate = ( key: keyof IAdminFormData, value?: string|Date ) => {
    setAdminData( { ...adminData, [key]: value } );
  };

  /**
   * Ensure that the form submissions are valid before sending data to the API.
   */
  const validateSubmission = () => {
    if ( !adminData.email?.match( /^.+@.+$/ ) ) {
      showError( 'Email address is not valid' );

      return false;
    } if ( !adminData.team ) {
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

    const newAdmin = {
      active: adminData.active,
      email: adminData.email.trim(),
      givenName: adminData.givenName.trim(),
      familyName: adminData.familyName.trim(),
      team: adminData.team,
      role: adminData.role,
    };

    if ( admin ) {
      await buildQuery( `admin?username=${adminData.email}`, { ...newAdmin }, 'PUT' )
        .then( () => window.location.assign( '/admins' ) )
        .catch( err => console.error( err ) );
    } else {
      await buildQuery( 'admin', { ...newAdmin, active: true }, 'POST' )
        .then( () => window.location.assign( '/admins' ) )
        .catch( err => console.error( err ) );
    }
  };

  const handleDeactivate = async () => {
    const { email, givenName, familyName } = adminData;
    const { isConfirmed } = await showConfirm( `Are you sure you want to deactivate ${givenName} ${familyName}?` );

    if ( !isConfirmed ) {
      return;
    }

    const { ok } = await buildQuery( `admin?username=${email}`, null, 'DELETE' );

    if ( !ok ) {
      showError( 'Unable to deactivate user' );
    } else {
      window.location.assign( '/admins' );
    }
  };

  const handleReactivate = async () => {
    const { email, givenName, familyName } = adminData;
    const { isConfirmed } = await showConfirm( `Are you sure you want to reactivate ${givenName} ${familyName}?` );

    if ( !isConfirmed ) {
      return;
    }

    await buildQuery( `admin?username=${email}`, { ...adminData, active: true }, 'PUT' )
      .then( () => window.location.assign( '/admins' ) )
      .catch( err => {
        showError( 'Unable to reactivate user' );
        console.error( err );
      } );
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
            value={ adminData.givenName }
            onChange={ e => handleUpdate( 'givenName', e.target.value ) }
          />
        </label>
        <label>
          <span>Family (Last) Name</span>
          <input
            id="family-name-input"
            type="text"
            required
            value={ adminData.familyName }
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
            disabled={ admin } // Email is the primary key for users so we prevent changes for existing admins.
            value={ adminData.email }
            onChange={ e => handleUpdate( 'email', e.target.value ) }
          />
        </label>
        <label>
          <span>Team</span>
          <select
            id="team-input"
            disabled={ !isAdmin }
            required
            value={ adminData.team }
            onChange={ e => handleUpdate( 'team', e.target.value ) }
          >
            <option value="">- Select a team -</option>
            { teamList.map( ( { id, name } ) => <option key={ id } value={ id }>{ name }</option> ) }
          </select>
        </label>
      </div>
      <div className="field-group">
        <label>
          <span>Admin Type</span>
          <select
            id="role-input"
            required
            value={ adminData.role }
            onChange={ e => handleUpdate( 'role', e.target.value ) }
          >
            <option value="admin">Admin</option>
            <option value="super admin">Super Admin</option>
          </select>
        </label>
      </div>
      <div style={ { textAlign: 'center' } }>
        <button
          className={ `${styles.btn} ${styles['spaced-btn']}` }
          id="update-btn"
          type="submit"
        >
          { admin ? 'Update Admin User' : 'Add Admin User' }
        </button>
        {
          admin && adminData.active && (
            <button
              className={ `${styles['btn-light']} ${styles['spaced-btn']}` }
              id="deactivate-btn"
              type="button"
              onClick={ handleDeactivate }
            >
              Deactivate Admin User
            </button>
          )
        }
        {
          admin && !adminData.active && (
            <button
              className={ `${styles['btn-light']} ${styles['spaced-btn']}` }
              id="deactivate-btn"
              type="button"
              onClick={ handleReactivate }
            >
              Reactivate Admin User
            </button>
          )
        }
        <BackButton text="Cancel" showConfirmDialog />
      </div>
    </form>
  );
};

export default AdminForm;
