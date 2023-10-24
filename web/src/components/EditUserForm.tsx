// ////////////////////////////////////////////////////////////////////////////
// React Imports
// ////////////////////////////////////////////////////////////////////////////
import { useEffect, useState } from 'react';
import type { FC, FormEvent } from 'react';

// ////////////////////////////////////////////////////////////////////////////
// 3PP Imports
// ////////////////////////////////////////////////////////////////////////////
import { addDays, parse } from 'date-fns';

// ////////////////////////////////////////////////////////////////////////////
// Local Imports
// ////////////////////////////////////////////////////////////////////////////
import BackButton from './BackButton';

import currentUser from '../stores/current-user';
import type { IInvite, TUserRole } from '../utils/types';
import { showConfirm, showError } from '../utils/alert';
import { buildQuery } from '../utils/api';
import { userIsAdmin } from '../utils/auth';
import { MAX_ACCESS_GRANT_DAYS } from '../utils/constants';
import { addDaysToNow, dateSelectionIsValid, getYearMonthDay, userWillNeedNewPassword } from '../utils/dates';

// ////////////////////////////////////////////////////////////////////////////
// Styles and CSS
// ////////////////////////////////////////////////////////////////////////////
import '../styles/form.scss';
import styles from '../styles/button.module.scss';
import { InviteModal } from './InviteModal';

// ////////////////////////////////////////////////////////////////////////////
// Interfaces and Types
// ////////////////////////////////////////////////////////////////////////////
interface IUserFormData {
  givenName: string;
  familyName: string;
  email: string;
  team: string;
  role: TUserRole,
}

const initialUserState: IUserFormData = {
  givenName: '',
  familyName: '',
  email: '',
  team: currentUser.get().team || '',

  role: 'guest' as TUserRole,
};

const initialInvites: IInvite[] = [
  {
    pending: false,
    expired: false,
    dateInvited: getYearMonthDay(new Date()),
    accessEndDate: getYearMonthDay(addDays(new Date(), 14)),
  }
];

// ////////////////////////////////////////////////////////////////////////////
// Interface and Implementation
// ////////////////////////////////////////////////////////////////////////////
const UserForm: FC = () => {
  const [isAdmin, setIsAdmin] = useState(false);
  const [teamList, setTeamList] = useState([]);

  const [userData, setUserData] = useState<IUserFormData>(initialUserState);
  const [currentInvite, setCurrentInvite] = useState<IInvite>(initialInvites[0]);
  const [invites, setInvites] = useState<IInvite[]>(initialInvites);

  const partnerRoles = [{ name: 'External Partner', value: 'guest' }, { name: 'External Team Lead', value: 'guest admin' }];

  // Check whether the user is an admin and set that value in state.
  // Doing so outside of a useEffect hook causes a mismatch in values
  // between the statically rendered portion and the client.
  useEffect(() => {
    setIsAdmin(userIsAdmin());
  }, []);

  // Generate the teams list.
  useEffect(() => {
    const getTeams = async () => {
      const response = await buildQuery('teams', null, 'GET');
      const { data } = await response.json();

      if (data) {
        // If the user is not an admin user, only allow them to invite users to their own team.
        const filtered = isAdmin ? data : data.filter((t: ITeam) => t.id === currentUser.get().team);

        setTeamList(filtered);
      }
    };

    getTeams();
  }, [isAdmin]);

  // Initialize the form.
  useEffect(() => {
    const getUser = async (id: string) => {
      const response = await buildQuery(`guest?id=${id}`, null, 'GET');
      const { data } = await response.json();

      if (data) {
        setUserData({
          ...data,
          accessEndDate: getYearMonthDay(parse(data.expiration, "yyyy-MM-dd'T'HH:mm:ssX", new Date())),
        });

        const fmtInvites: IInvite[] = data.invites.map((invite: any) => ({
          pending: invite.pending,
          expired: invite.expired,
          dateInvited: getYearMonthDay(parse(invite.dateInvited, "yyyy-MM-dd'T'HH:mm:ssX", new Date())),
          accessEndDate: getYearMonthDay(parse(invite.expiration, "yyyy-MM-dd'T'HH:mm:ssX", new Date())),
        }))

        setCurrentInvite(fmtInvites[0])
        setInvites(fmtInvites);
      }
    };

    const urlSearchParams = new URLSearchParams(window.location.search);
    const { id } = Object.fromEntries(urlSearchParams.entries());

    getUser(id);
  }, []);

  /**
   * Updates the user state on changed to the form inputs.
   * @param key The user property being updated.
   */
  const handleUserUpdate = (key: keyof IUserFormData, value?: string | Date) => {
    setUserData({ ...userData, [key]: value });
  };

  const handleAccessUpdate = (value: string) => {
    setCurrentInvite({ ...currentInvite, accessEndDate: value });
  }

  /**
   * Ensure that the form submissions are valid before sending data to the API.
   */
  const validateUserSub = () => {
    if (!userData.email?.match(/^.+@.+$/)) {
      showError('Email address is not valid');

      return false;
    } 
    
    if (isAdmin && !userData.team) {
      // Admin users have the option to set a team, so a team should be set
      showError('Please assign this user to a valid team');

      return false;
    }

    return true;
  };

  const validateAccessSub = () => {
    if ( !dateSelectionIsValid( currentInvite.accessEndDate ) ) {
      showError( `Please select an access grant end date after today and no more than ${MAX_ACCESS_GRANT_DAYS} in the future` );

      return false;
    }

    return true;
  };

  /**
   * Validates the form inputs and sends the provided data to the API.
   * @param e The submission event - simply used to prevent the default page refresh.
   */
  const handleUserSubmit = async (e: FormEvent<HTMLFormElement>) => {
    e.preventDefault();

    if (!validateUserSub()) {
      return;
    }

    const invitee = {
      email: userData.email.trim(),
      givenName: userData.givenName.trim(),
      familyName: userData.familyName.trim(),
      team: userData.team || currentUser.get().team,
      role: userData.role,
    };

    buildQuery('guest', invitee, 'PUT').catch(err => console.error(err));
  };

  const handleReauth = async () => {
    const { email } = userData;
    const { dateInvited, accessEndDate, expired } = currentInvite;

    const shouldPrompt = userWillNeedNewPassword( dateInvited, accessEndDate, expired );
    if( shouldPrompt ) {
      const { isConfirmed } = await showConfirm(`Updating this user's access end date to ${accessEndDate} will trigger a password reset.  Continue?`);

      if (!isConfirmed) {
        return;
      }
    }

    // Conversion to iso required by Lambda
    const expiration = new Date( accessEndDate ).toISOString();
    
    const body = {
      email,
      expiration,
    };

    const { ok } = await buildQuery('guest/reauth', body, 'POST');

    if (!ok) {
      showError('Unable to re-authorize user');
    }
  };

  const handleRevoke = async () => {
    const { email, givenName, familyName } = userData;
    const { isConfirmed } = await showConfirm(`Are you sure you want to deactivate ${givenName} ${familyName}?`);

    if (!isConfirmed) {
      return;
    }

    const { ok } = await buildQuery(`guest?id=${email}`, null, 'DELETE');

    if (!ok) {
      showError('Unable to deactivate user');
    } else {
      window.location.assign((isAdmin ? '/' : '/uploader-users'));
    }
  };

  return (
    <div>
      <div id="guest-data-form">
        <h3>Guest Information</h3>
        <form onSubmit={handleUserSubmit}>
          <div className="field-group">
            <label>
              <span>Given (First) Name</span>
              <input
                id="given-name-input"
                type="text"
                required
                value={userData.givenName}
                onChange={e => handleUserUpdate('givenName', e.target.value)}
              />
            </label>
            <label>
              <span>Family (Last) Name</span>
              <input
                id="family-name-input"
                type="text"
                required
                value={userData.familyName}
                onChange={e => handleUserUpdate('familyName', e.target.value)}
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
                disabled={true} // Email is the primary key for users so we prevent changes for existing users.
                value={userData.email}
              />
            </label>
            <label>
              <span>Team</span>
              <select
                id="team-input"
                disabled={!isAdmin}
                required
                value={userData.team}
                onChange={e => handleUserUpdate('team', e.target.value)}
              >
                <option value="">- Select a team -</option>
                {teamList.map(({ id, name }) => <option key={id} value={id}>{name}</option>)}
              </select>
            </label>
          </div>
          <div className="field-group">
            {isAdmin && (
              <label>
                <span>User Role</span>
                <select
                  id="role-input"
                  required
                  value={userData.role}
                  onChange={e => handleUserUpdate('role', e.target.value)}
                >
                  {partnerRoles.map(({ name, value }) => <option key={value} value={value}>{name}</option>)}
                </select>
              </label>
            )}
          </div>
          <button
            className={`${styles.btn} ${styles['spaced-btn']}`}
            id="update-btn"
            type="submit"
          >
            Update Guest
          </button>
        </form>
      </div>
      <div id="invite-data-form">
        <h3>{currentInvite.expired ? 'Recent' : 'Current'} Access</h3>
        <div className="field-group">
          <label>
            <span>Access End Date</span>
            <input
              id="date-input"
              type="date"
              disabled={!isAdmin}
              min={getYearMonthDay(new Date())}
              max={getYearMonthDay(addDaysToNow(60))}
              value={currentInvite.accessEndDate}
              onChange={e => handleAccessUpdate(e.target.value)}
            />
          </label>
          <label>
            <span>Date Invited</span>
            <input
              id="date-invited-input"
              type="date"
              disabled={true}
              value={currentInvite.dateInvited}
            />
          </label>
        </div>
        <button
          className={`${styles.btn} ${styles['spaced-btn']}`}
          id="update-btn"
          type="button"
          onClick={handleReauth}
        >
          Reauthorize
        </button>
        <button
          className={`${styles['btn-light']} ${styles['spaced-btn']}`}
          id="deactivate-btn"
          type="button"
          onClick={handleRevoke}
        >
          Revoke Access
        </button>
      </div>
      <div id="additional-options">
        <h3>Additional Options</h3>
        <InviteModal invites={invites} anchor={"Invite History"} />
        <BackButton text="Cancel" showConfirmDialog />
      </div>
    </div>
  );
};

export default UserForm;
