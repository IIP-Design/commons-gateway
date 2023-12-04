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
import { escapeQueryStrings } from '../utils/string';
import { defaultColumnDef } from './Table';
import TableWrapper from './TableWrapper';

// ////////////////////////////////////////////////////////////////////////////
// Styles and CSS
// ////////////////////////////////////////////////////////////////////////////
import tableStyles from '../styles/table.module.scss';

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
  const [loading, setLoading] = useState<boolean>( true );

  const getTeams = async () => {
    const response = await buildQuery( 'teams', null, 'GET' );
    const { data } = await response.json();

    return data as ITeam[];
  };

  useEffect( () => {
    const body = {
      ...( userIsSuperAdmin() ? {} : { team: currentUser.get().team } ),
      ...( role ? { role } : {} ),
    };

    const getUsers = async ( teams: ITeam[] ) => {
      const response = await buildQuery( 'guests', body );
      const { data } = await response.json();

      if ( data ) {
        setUsers(
          data.map( ( user: IUserEntry ) => ( {
            ...user,
            name: `${user.givenName} ${user.familyName}`,
            active: isGuestActive( user.expires ),
            team: getTeamName( user.team, teams ),
          } ) ),
        );
      }
    };

    getTeams()
      .then( teams => getUsers( teams ) )
      .finally( () => setLoading( false ) );
  }, [role] );

  const columns = useMemo<ColumnDef<WithUiData<IUserEntry>>[]>(
    () => [
      {
        ...defaultColumnDef( 'name' ),
        cell: info => <a href={ `/edit-user?id=${escapeQueryStrings( info.row.getValue( 'email' ) )}` }>{ info.getValue() as string }</a>,
      },
      defaultColumnDef( 'email' ),
      defaultColumnDef( 'team' ),
      {
        ...defaultColumnDef( 'active', 'Status' ),
        cell: info => {
          const isPending = info.row.original.pending as boolean;
          const isActive = info.getValue() as boolean;

          const baseStyle = isPending ? tableStyles.pending : tableStyles.active;
          const baseLabel = isPending ? 'Pending' : 'Active';

          return (
            <span className={ tableStyles.status }>
              <span className={ isActive ? baseStyle : tableStyles.inactive } />
              { isActive ? baseLabel : 'Inactive' }
            </span>
          );
        },
      },
    ],
    [],
  );

  return (
    <div style={ { display: 'flex' } }>
      <TableWrapper
        loading={ loading }
        table={ { data: users, columns, additionalTableClasses: ['user-table'] } }
      />
    </div>
  );
};

export default UserTable;
