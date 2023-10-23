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
import { addDaysToNow, dateSelectionIsValid, getYearMonthDay } from '../utils/dates';

// ////////////////////////////////////////////////////////////////////////////
// Styles and CSS
// ////////////////////////////////////////////////////////////////////////////
import '../styles/form.scss';
import styles from '../styles/button.module.scss';
import tableStyles from '../styles/table.module.scss';
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

interface IUserFormProps {
  readonly user?: boolean;
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
const UserForm: FC<IUserFormProps> = ({ user }) => {
  const [isAdmin, setIsAdmin] = useState(false);
  const [userExists, setUserExists] = useState(false);
  const [teamList, setTeamList] = useState([]);

  const [userData, setUserData] = useState<IUserFormData>(initialUserState);
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

        setInvites(data.invites.map((invite) => ({
          pending: invite.pending,
          expired: invite.expired,
          dateInvited: getYearMonthDay(parse(invite.dateInvited, "yyyy-MM-dd'T'HH:mm:ssX", new Date())),
          accessEndDate: getYearMonthDay(parse(invite.expiration, "yyyy-MM-dd'T'HH:mm:ssX", new Date())),
        })));

        setUserExists(true);
      }
    };

    if (user) {
      const urlSearchParams = new URLSearchParams(window.location.search);
      const { id } = Object.fromEntries(urlSearchParams.entries());

      getUser(id);
    }
  }, [user]);

  /**
   * Updates the user state on changed to the form inputs.
   * @param key The user property being updated.
   */
  // FIXME LOGIC
  const handleUpdate = (key: keyof IUserFormData, value?: string | Date) => {
    setUserData({ ...userData, [key]: value });
  };

  /**
   * Ensure that the form submissions are valid before sending data to the API.
   */
  // FIXME LOGIC
  const validateSubmission = () => {
    if (!userData.email?.match(/^.+@.+$/)) {
      showError('Email address is not valid');

      return false;
    } if (!dateSelectionIsValid(userData.accessEndDate)) {
      showError(`Please select an access grant end date after today and no more than ${MAX_ACCESS_GRANT_DAYS} in the future`);

      return false;
    } if (isAdmin && !userData.team) {
      // Admin users have the option to set a team, so a team should be
      showError('Please assign this user to a valid team');

      return false;
    }

    return true;
  };

  /**
   * Validates the form inputs and sends the provided data to the API.
   * @param e The submission event - simply used to prevent the default page refresh.
   */
  // FIXME LOGIC
  const handleSubmit = async (e: FormEvent<HTMLFormElement>) => {
    e.preventDefault();

    if (!validateSubmission()) {
      return;
    }

    // Conversion to iso required by Lambda
    const expiration = new Date(userData.accessEndDate).toISOString();

    const invitee = {
      email: userData.email.trim(),
      givenName: userData.givenName.trim(),
      familyName: userData.familyName.trim(),
      team: userData.team || currentUser.get().team,
      role: userData.role,
    };

    if (user) {
      buildQuery('guest', { ...invitee, expiration }, 'PUT')
        .then(() => window.location.assign((isAdmin ? '/' : '/uploader-users')))
        .catch(err => console.error(err));
    } else if (isAdmin) {
      const invitation = {
        inviter: currentUser.get().email,
        invitee,
        expiration,
      };

      buildQuery('creds/provision', invitation, 'POST')
        .then(() => window.location.assign('/'))
        .catch(err => console.error(err));
    } else {
      const invitation = {
        proposer: currentUser.get().email,
        invitee,
        expiration,
      };

      try {
        await buildQuery('creds/propose', invitation, 'POST');
        window.location.assign('/uploader-users');
      } catch (err: any) {
        console.error(err);
      }
    }
  };

  // FIXME LOGIC
  const handleDeactivate = async () => {
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
        <form onSubmit={handleSubmit}>
          <div className="field-group">
            <label>
              <span>Given (First) Name</span>
              <input
                id="given-name-input"
                type="text"
                required
                value={userData.givenName}
                onChange={e => handleUpdate('givenName', e.target.value)}
              />
            </label>
            <label>
              <span>Family (Last) Name</span>
              <input
                id="family-name-input"
                type="text"
                required
                value={userData.familyName}
                onChange={e => handleUpdate('familyName', e.target.value)}
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
                disabled={user} // Email is the primary key for users so we prevent changes for existing users.
                value={userData.email}
                onChange={e => handleUpdate('email', e.target.value)}
              />
            </label>
            <label>
              <span>Team</span>
              <select
                id="team-input"
                disabled={!isAdmin}
                required
                value={userData.team}
                onChange={e => handleUpdate('team', e.target.value)}
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
                  onChange={e => handleUpdate('role', e.target.value)}
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
      <div id="guest-data-form">
        <h3>{user ? (invites.length && !invites[0].expired ? 'Current' : 'Recent') : 'Proposed'} Access</h3>
        <label>
          <span>Access End Date</span>
          <input
            id="date-input"
            type="date"
            disabled={!isAdmin && userExists}
            min={getYearMonthDay(new Date())}
            max={getYearMonthDay(addDaysToNow(60))}
            value={invites[0].accessEndDate}
            onChange={e => handleUpdate('accessEndDate', e.target.value)}
          />
        </label>
        <button
          className={`${styles.btn} ${styles['spaced-btn']}`}
          id="update-btn"
          type="button"
        >
          Reauthorize
        </button>
        {user && (
          <button
            className={`${styles['btn-light']} ${styles['spaced-btn']}`}
            id="deactivate-btn"
            type="button"
            onClick={handleDeactivate}
          >
            Revoke Access
          </button>
        )}
      </div>
      <div>
        <h3>Additional Options</h3>
        <InviteModal invites={invites} anchor={"Invite History"} />
        <BackButton text="Cancel" showConfirmDialog />
      </div>
    </div>
  );
};

export default UserForm;
