// ////////////////////////////////////////////////////////////////////////////
// React Imports
// ////////////////////////////////////////////////////////////////////////////
import { useEffect, useMemo, useState } from 'react';
import type { FC } from 'react';

// ////////////////////////////////////////////////////////////////////////////
// 3PP Imports
// ////////////////////////////////////////////////////////////////////////////
import type { ColumnDef } from '@tanstack/react-table';

// ////////////////////////////////////////////////////////////////////////////
// Local Imports
// ////////////////////////////////////////////////////////////////////////////
import currentUser from '../stores/current-user';
import { daysUntil } from '../utils/dates';
import type { IUserEntry, WithUiData } from '../utils/types';
import { buildQuery } from '../utils/api';
import { userIsSuperAdmin } from '../utils/auth';
import { getTeamName } from '../utils/team';
import { Table, defaultColumnDef } from './Table';
import { makeApproveUserHandler } from '../utils/users';

// ////////////////////////////////////////////////////////////////////////////
// Styles and CSS
// ////////////////////////////////////////////////////////////////////////////
import btnStyle from '../styles/button.module.scss';
import style from '../styles/table.module.scss';

// ////////////////////////////////////////////////////////////////////////////
// Types and Interfaces
// ////////////////////////////////////////////////////////////////////////////
interface IInvite extends IUserEntry {
  dateInvited: string;
  proposer: string;
}

// ////////////////////////////////////////////////////////////////////////////
// Implementation
// ////////////////////////////////////////////////////////////////////////////

const UserTable: FC = () => {
  const [users, setUsers] = useState<WithUiData<IInvite>[]>( [] );
  const [teams, setTeams] = useState<ITeam[]>( [] );

  useEffect( () => {
    const body = userIsSuperAdmin() ? {} : { team: currentUser.get().team };

    const getUsers = async () => {
      const response = await buildQuery( 'guests/pending', body );
      const { data } = await response.json();

      if ( data ) {
        setUsers(
          data.map( ( user: IInvite ) => ( {
            ...user,
            name: `${user.givenName} ${user.familyName}`,
          } ) ),
        );
      }
    };

    getUsers();
  }, [] );

  useEffect( () => {
    const getTeams = async () => {
      const response = await buildQuery( 'teams', null, 'GET' );
      const { data } = await response.json();

      if ( data ) {
        setTeams( data );
      }
    };

    getTeams();
  }, [] );

  const columns = useMemo<ColumnDef<WithUiData<IInvite>>[]>(
    () => [
      {
        ...defaultColumnDef( 'name' ),
        cell: info => (
          <button
            className={ btnStyle['link-btn'] }
            onClick={ makeApproveUserHandler( info.row.getValue( 'email' ) ) }
            type="button"
          >
            { info.getValue() as string }
          </button>
        ),
      },
      defaultColumnDef( 'email' ),
      {
        ...defaultColumnDef( 'team' ),
        cell: info => getTeamName( info.getValue() as string, teams ),
      },
      defaultColumnDef( 'proposer' ),
      {
        accessorFn: row => row.expiration,
        id: '_exp',
        header: 'Days Till Expiration',
        footer: props => props.column.id,
        enableSorting: true,
        cell: info => daysUntil( info.getValue() as string ),
      },
      {
        accessorFn: row => row.dateInvited,
        id: '_inv',
        header: 'Days Since Invite',
        footer: props => props.column.id,
        enableSorting: true,
        cell: info => daysUntil( info.getValue() as string ) * -1,
      },
    ],
    [teams],
  );

  return (
    <div style={ { display: 'flex', marginBottom: '0.75em' } }>
      { users.length
        ? (
          <Table
            {
              ...{
                data: users,
                columns,
                additionalTableClasses: ['user-table'],
              }
            }
          />
        )
        : <p className={ style['no-data'] }>No pending invites at this time</p> }
    </div>
  );
};

export default UserTable;
