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
import type { TUserRole, WithUiData } from '../utils/types';
import { buildQuery } from '../utils/api';
import { getTeamName } from '../utils/team';
import { escapeQueryStrings } from '../utils/string';
import { defaultColumnDef } from './Table';
import TableWrapper from './TableWrapper';

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

// ////////////////////////////////////////////////////////////////////////////
// Implementation
// ////////////////////////////////////////////////////////////////////////////
const AdminTable: FC = () => {
  const [admins, setAdmins] = useState<WithUiData<IAdminUser>[]>( [] );
  const [teams, setTeams] = useState<ITeam[]>( [] );
  const [loading, setLoading] = useState<boolean>( true );

  useEffect( () => {
    const getAdmins = async () => {
      const response = await buildQuery( 'admins', null, 'GET' );
      const { data } = await response.json();

      if ( data ) {
        setAdmins(
          data.map( ( user: IAdminUser ) => ( {
            ...user,
            name: `${user.givenName} ${user.familyName}`,
          } ) ),
        );
      }
    };

    getAdmins().finally( () => setLoading( false ) );
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

  const columns = useMemo<ColumnDef<WithUiData<IAdminUser>>[]>(
    () => [
      {
        ...defaultColumnDef( 'name' ),
        cell: info => <a href={ `/edit-admin?id=${escapeQueryStrings( info.row.getValue( 'email' ) )}` }>{ info.getValue() as string }</a>,
      },
      {
        ...defaultColumnDef( 'role' ),
        cell: info => <span style={ { textTransform: 'capitalize' } }>{ info.getValue() as string }</span>,
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
    [teams],
  );

  return (
    <div style={ { display: 'flex', marginBottom: '0.75em' } }>
      <TableWrapper
        loading={ loading }
        table={ { data: admins, columns } }
      />
    </div>
  );
};

export default AdminTable;
