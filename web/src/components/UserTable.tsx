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
import type { IUserEntry, TUserRole, WithUiData } from '../utils/types';
import { buildQuery } from '../utils/api';
import { userIsSuperAdmin } from '../utils/auth';
import { isGuestActive } from '../utils/guest';
import { getTeamName } from '../utils/team';
import { Table, defaultColumnDef } from './Table';

// ////////////////////////////////////////////////////////////////////////////
// Styles and CSS
// ////////////////////////////////////////////////////////////////////////////
import style from '../styles/table.module.scss';

// ////////////////////////////////////////////////////////////////////////////
// Types and Interfaces
// ////////////////////////////////////////////////////////////////////////////
interface IUserTableProps {
  readonly role?: TUserRole;
}

// ////////////////////////////////////////////////////////////////////////////
// Implementation
// ////////////////////////////////////////////////////////////////////////////

const UserTable: FC<IUserTableProps> = ( { role }: IUserTableProps ) => {
  const [users, setUsers] = useState<WithUiData<IUserEntry>[]>( [] );
  const [teams, setTeams] = useState<ITeam[]>( [] );

  useEffect( () => {
    const body = {
      ...( userIsSuperAdmin() ? {} : { team: currentUser.get().team } ),
      ...( role ? { role } : {} ),
    };

    const getUsers = async () => {
      const response = await buildQuery( 'guests', body );
      const { data } = await response.json();

      if ( data ) {
        setUsers(
          data.map( ( user: IUserEntry ) => ( {
            ...user,
            name: `${user.givenName} ${user.familyName}`,
            active: isGuestActive( user.expires ),
          } ) ),
        );
      }
    };

    getUsers();
  }, [role] );

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

  const columns = useMemo<ColumnDef<WithUiData<IUserEntry>>[]>(
    () => [
      {
        ...defaultColumnDef( 'name' ),
        cell: info => <a href={ `/edit-user?id=${info.row.getValue( 'email' )}` }>{ info.getValue() as string }</a>,
      },
      defaultColumnDef( 'email' ),
      {
        ...defaultColumnDef( 'team' ),
        cell: info => getTeamName( info.getValue() as string, teams ),
      },
      {
        ...defaultColumnDef( 'active', 'Status' ),
        cell: info => {
          const isPending = info.row.original["pending"] as boolean;
          const isActive = info.getValue() as boolean;

          return (
            <span className={ style.status }>
              <span className={ isActive ? ( isPending ? style.pending : style.active ) : style.inactive } />
              { isPending ? 'Pending' : ( isActive ? 'Active' : 'Inactive' ) }
            </span>
          );
        },
      },
    ],
    [teams],
  );

  return (
    <div style={ { display: 'flex' } }>
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
        : <p className={ style['no-data'] }>No data to show</p> }
    </div>
  );
};

export default UserTable;
