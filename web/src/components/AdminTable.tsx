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
import type { TUserRole } from '../stores/current-user';
import { buildQuery } from '../utils/api';
import { getTeamName } from '../utils/team';
import { Table, defaultColumnDef } from './Table';

// ////////////////////////////////////////////////////////////////////////////
// Styles and CSS
// ////////////////////////////////////////////////////////////////////////////
import style from '../styles/table.module.scss';

// ////////////////////////////////////////////////////////////////////////////
// Interfaces and Types
// ////////////////////////////////////////////////////////////////////////////
interface IAdminUser {
  active: boolean;
  email: string;
  familyName: string;
  givenName: string;
  role: TUserRole;
  team: string;
}

interface IAdminWithUiData extends IAdminUser {
  name: string;
}

// ////////////////////////////////////////////////////////////////////////////
// Implementation
// ////////////////////////////////////////////////////////////////////////////
const AdminTable: FC = () => {
  const [admins, setAdmins] = useState<IAdminWithUiData[]>( [] );
  const [teams, setTeams] = useState<ITeam[]>( [] );

  useEffect( () => {
    const getAdmins = async () => {
      const response = await buildQuery( 'admins', null, 'GET' );
      const { data } = await response.json();

      if ( data ) {
        setAdmins(
          data.map( ( user: IAdminUser ) => {
            return {
              ...user,
              name: `${user.givenName} ${user.familyName}`,
            };
          } )
        );
      }
    };

    getAdmins();
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

  const columns = useMemo<ColumnDef<IAdminWithUiData>[]>(
    () => [
      {
        ...defaultColumnDef( 'name' ),
        cell: info => <a href={`/editAdmin?id=${info.row.getValue('email')}`}>{info.getValue() as string}</a>,
      },
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
    ],
    [teams]
  );

  return (
    <div className={ style.container }>
      <Table
        {
          ...{
            data: admins,
            columns,
          }
        }
      />
    </div>
  );
}

export default AdminTable;
