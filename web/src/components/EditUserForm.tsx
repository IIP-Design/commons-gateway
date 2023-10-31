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
import { showConfirm, showError, showSuccess } from '../utils/alert';
import { buildQuery } from '../utils/api';
import { userIsAdmin } from '../utils/auth';
import { MAX_ACCESS_GRANT_DAYS } from '../utils/constants';
import { addDaysToNow, dateSelectionIsValid, getYearMonthDay, userWillNeedNewPassword } from '../utils/dates';

import type { IInvite } from '../utils/types';
import { makeDummyUserForm, type IUserFormData, makeApproveUserHandler } from '../utils/users';
import { InviteModal } from './InviteModal';

// ////////////////////////////////////////////////////////////////////////////
// Styles and CSS
// ////////////////////////////////////////////////////////////////////////////
import '../styles/form.scss';
import btnStyles from '../styles/button.module.scss';
import tblStyles from '../styles/table.module.scss';

// ////////////////////////////////////////////////////////////////////////////
// Interfaces and Types
// ////////////////////////////////////////////////////////////////////////////
interface IStatusTokenProps {
  active: boolean;
}

interface IInviteWidgetParams {
  readonly userData: IUserFormData;
  readonly invite: IInvite;
  readonly isAdmin: boolean;
}

// ////////////////////////////////////////////////////////////////////////////
// Helpers
// ////////////////////////////////////////////////////////////////////////////
const StatusToken: FC<IStatusTokenProps> = ( { active }: IStatusTokenProps ) => {
  return (
    <span className={ tblStyles.status } style={{display: "inline"}}>
      <span className={ active ? tblStyles.active : tblStyles.inactive } />
    </span>
  );
}

const CurrentInvite: FC<IInviteWidgetParams> = ( { userData, invite, isAdmin }: IInviteWidgetParams ) => {
  const [currentInvite, setCurrentInvite] = useState<IInvite>( invite );
  const [updated, setUpdated] = useState( false );

  const handleAccessUpdate = (value: string) => {
    setCurrentInvite({ ...currentInvite, accessEndDate: value });
    setUpdated(true);
  };

  const validateAccessSub = () => {
    if ( !dateSelectionIsValid( currentInvite.accessEndDate ) ) {
      showError( `Please select an access grant end date after today and no more than ${MAX_ACCESS_GRANT_DAYS} in the future` );

      return false;
    }

    return true;
  };

  const handleReauth = async () => {
    if (!validateAccessSub()) {
      return;
    }

    const { email } = userData;
    const { dateInvited, accessEndDate, expired, passwordReset } = currentInvite;

    const shouldPrompt = userWillNeedNewPassword( dateInvited, accessEndDate, expired, passwordReset );
    if( shouldPrompt ) {
      let prompText = '';
      if( expired ) {
        prompText = 'Because the user\'s access has expired, they will need to reset their password.';
      } else if( !passwordReset ) {
        prompText = 'Because the user did not reset their password following the last invite, they must reset it after this one.';
      } else {
        prompText = `Updating this user's access end date to ${accessEndDate} will trigger a password reset.  Continue?`;
      }

      const { isConfirmed } = await showConfirm(prompText);
      if (!isConfirmed) {
        return;
      }
    }

    // Conversion to iso required by Lambda
    const expiration = new Date( accessEndDate ).toISOString();
    
    const body = {
      email,
      admin: currentUser.get().email,
      expiration,
    };

    const { ok } = await buildQuery('guest/reauth', body, 'POST');

    if (!ok) {
      showError('Unable to re-authorize user');
    } else {
      showSuccess(`User reauthorized until ${accessEndDate}`)
      window.location.reload();
    }
  };

  const handlePasswordReset = async () => {
    const { email } = userData;
    const { ok } = await buildQuery(`passwordReset?id=${email}`, null, 'POST');

    if (!ok) {
      showError('Unable to reset password');
    } else {
      showSuccess('User password reset')
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
    <div id="invite-data-form">
        <h3>{currentInvite.expired ? 'Most Recent' : 'Current'} Access <StatusToken active={!currentInvite.expired} /></h3>
        <div className="field-group">
          <label>
            <span>Access End Date</span>
            <input
              id="date-input"
              type="date"
              disabled={!isAdmin && !currentInvite.expired}
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
          className={`${btnStyles.btn} ${updated ? "" : btnStyles['disabled-btn']} ${btnStyles['spaced-btn']}`}
          id="update-btn"
          type="button"
          onClick={handleReauth}
          disabled={!updated}
        >
          { isAdmin ? "Reauthorize" : "Re-Propose" }
        </button>
        <button
          className={`${btnStyles['btn-light']} ${!currentInvite.expired ? "" : btnStyles['disabled-btn']} ${btnStyles['spaced-btn']}`}
          id="reset-password-btn"
          type="button"
          onClick={handlePasswordReset}
          disabled={currentInvite.expired}
        >
          Reset Password
        </button>
        <button
          className={`${btnStyles.btn} ${btnStyles['back-btn']} ${btnStyles['spaced-btn']}`}
          id="deactivate-btn"
          type="button"
          onClick={handleRevoke}
        >
          Revoke Access
        </button>
      </div>
  );
};

const PendingInvite: FC<IInviteWidgetParams> = ( { userData: { email }, invite, isAdmin }: IInviteWidgetParams ) => {
  const [pendingInvite] = useState<IInvite>( invite );
  
  return (
    <div id="invite-data-form">
        <h3>Proposed Access</h3>
        <div className="field-group">
          <label>
            <span>Access End Date</span>
            <input
              id="date-input"
              type="date"
              disabled={true}
              value={pendingInvite.accessEndDate}
            />
          </label>
          <label>
            <span>Date Invited</span>
            <input
              id="date-invited-input"
              type="date"
              disabled={true}
              value={pendingInvite.dateInvited}
            />
          </label>
        </div>
        {
          isAdmin ?
            <button
              className={`${btnStyles['btn-light']}`}
              id="finalize-btn"
              type="button"
              onClick={makeApproveUserHandler(email)}
            >
              Finalize
            </button>
          : <p>Proposed by { pendingInvite.proposer || "N/A" }</p>
        }
        
      </div>
  );
};

// ////////////////////////////////////////////////////////////////////////////
// Interface and Implementation
// ////////////////////////////////////////////////////////////////////////////
const UserForm: FC = () => {
  const [isAdmin, setIsAdmin] = useState(false);
  const [updated, setUpdated] = useState(false);
  const [teamList, setTeamList] = useState([]);

  const [userData, setUserData] = useState<IUserFormData>(makeDummyUserForm());
  const [pendingInvite, setPendingInvite] = useState<IInvite|null>(null);
  const [currentInvite, setCurrentInvite] = useState<IInvite|null>(null);
  const [invites, setInvites] = useState<IInvite[]>([]);

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
          proposer: invite.proposer?.String || null,
          pending: invite.pending,
          expired: invite.expired,
          dateInvited: getYearMonthDay(parse(invite.dateInvited, "yyyy-MM-dd'T'HH:mm:ssX", new Date())),
          accessEndDate: getYearMonthDay(parse(invite.expiration, "yyyy-MM-dd'T'HH:mm:ssX", new Date())),
        }));

        setPendingInvite(fmtInvites[0].pending ? fmtInvites[0] : null);
        setCurrentInvite(fmtInvites.find( val => !val.pending ) || null);
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
    setUpdated( true );
  };

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

    buildQuery('guest', invitee, 'PUT')
      .then( () => setUpdated( false ) )
      .catch(err => console.error(err));
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
            className={`${btnStyles.btn} ${btnStyles['spaced-btn']}`}
            id="update-btn"
            type="submit"
          >
            Update Guest
          </button>
        </form>
      </div>
      { pendingInvite && <PendingInvite userData={userData} invite={pendingInvite} isAdmin={isAdmin} /> }
      { currentInvite && <CurrentInvite userData={userData} invite={currentInvite} isAdmin={isAdmin} /> }
      <div id="additional-options">
        <h3>Additional Options</h3>
        <InviteModal invites={invites} anchor={"Invite History"} />
        <BackButton text="Cancel" showConfirmDialog={updated} />
      </div>
    </div>
  );
};

export default UserForm;
