// ////////////////////////////////////////////////////////////////////////////
// React Imports
// ////////////////////////////////////////////////////////////////////////////
import { useEffect, useMemo, useState } from 'react';
import type { FC } from 'react';

// ////////////////////////////////////////////////////////////////////////////
// 3PP Imports
// ////////////////////////////////////////////////////////////////////////////
import { LiaUserEditSolid } from 'react-icons/lia';

// ////////////////////////////////////////////////////////////////////////////
// Local Imports
// ////////////////////////////////////////////////////////////////////////////
import currentUser from '../stores/current-user';
import type { TUserRole } from '../stores/current-user';
import { buildQuery } from '../utils/api';
import { userIsSuperAdmin } from '../utils/auth';
import { isGuestActive } from '../utils/guest';
import { getTeamName } from '../utils/team';

// ////////////////////////////////////////////////////////////////////////////
// Styles and CSS
// ////////////////////////////////////////////////////////////////////////////
import 'bootstrap/dist/css/bootstrap.css';

import style from '../styles/table.module.scss';
import type { ColumnDef } from '@tanstack/react-table';
import { Table, defaultColumnDef } from './Table';

// ////////////////////////////////////////////////////////////////////////////
// Types and Interfaces
// ////////////////////////////////////////////////////////////////////////////
interface IUserEntry {
  email: string;
  expiration: string;
  familyName: string;
  givenName: string;
  role: TUserRole;
  team: string;
}

interface ITeam {
  active: boolean;
  id: string;
  name: string;
}

type WithActiveTag<T> = T & { active: boolean; };

// ////////////////////////////////////////////////////////////////////////////
// Implementation
// ////////////////////////////////////////////////////////////////////////////

const UserTable: FC = () => {
  const [users, setUsers] = useState<WithActiveTag<IUserEntry>[]>( [] );
  const [teams, setTeams] = useState<ITeam[]>( [] );

  useEffect( () => {
    const body = userIsSuperAdmin() ? {} : { team: currentUser.get().team };

    const getUsers = async () => {
      const response = await buildQuery( 'guests', body );
      const { data } = await response.json();

      if ( data ) {
        setUsers(
          data.map( ( user: IUserEntry ) => {
            return {
              ...user,
              active: isGuestActive( user.expiration ),
            };
          } )
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

  const columns = useMemo<ColumnDef<WithActiveTag<IUserEntry>>[]>(
    () => [
      defaultColumnDef( 'givenName' ),
      defaultColumnDef( 'familyName' ),
      defaultColumnDef( 'email' ),
      {
        ...defaultColumnDef( 'team' ),
        cell: info => getTeamName( info.getValue() as string, teams ),
      },
      {
        ...defaultColumnDef( 'active' ),
        cell: info => {
          const isActive = info.getValue() as boolean;
          return (
            <span className={ style.status }>
              <span className={ isActive ? style.active : style.inactive } />
              { isActive ? 'Active' : 'Inactive' }
            </span>
          );
        },
      },
      {
        accessorFn: row => row.email,
        id: "_edit",
        cell: info => <a href={`/editUser?id=${info.getValue()}`}><LiaUserEditSolid /></a>,
        header: () => "Edit",
        enableSorting: false,
      }
    ],
    [teams]
  );

  return (
    <div className={ style.container }>
      <Table
        {
          ...{
            data: users,
            columns,
            additionalTableClasses: [ 'user-table' ],
          }
        }
      />
    </div>
  );
}

export default UserTable;
