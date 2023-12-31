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
import type { IUserEntry, WithUiData } from '../utils/types';
import { buildQuery } from '../utils/api';
import { isGuestActive } from '../utils/guest';
import { defaultColumnDef } from './Table';
import { escapeQueryStrings } from '../utils/string';
import TableWrapper from './TableWrapper';

// ////////////////////////////////////////////////////////////////////////////
// Styles and CSS
// ////////////////////////////////////////////////////////////////////////////
import tableStyles from '../styles/table.module.scss';

// ////////////////////////////////////////////////////////////////////////////
// Types and Interfaces
// ////////////////////////////////////////////////////////////////////////////
interface IUploader extends IUserEntry {
  dateInvited: string;
  proposer: string;
  inviter: string;
  pending: boolean;
}

// ////////////////////////////////////////////////////////////////////////////
// Implementation
// ////////////////////////////////////////////////////////////////////////////

const UploaderTable: FC = () => {
  const [users, setUsers] = useState<WithUiData<IUploader>[]>( [] );
  const [loading, setLoading] = useState<boolean>( true );

  useEffect( () => {
    const body = { team: currentUser.get().team };

    const getUsers = async () => {
      const response = await buildQuery( 'guests/uploaders', body );
      const { data } = await response.json();

      if ( data ) {
        setUsers(
          data.map( ( user: IUploader ) => ( {
            ...user,
            name: `${user.givenName} ${user.familyName}`,
            active: isGuestActive( user.expires ),
          } ) ),
        );
      }
    };

    getUsers().finally( () => setLoading( false ) );
  }, [] );

  const columns = useMemo<ColumnDef<WithUiData<IUploader>>[]>(
    () => [
      {
        ...defaultColumnDef( 'name' ),
        cell: info => <a href={ `/edit-user?id=${escapeQueryStrings( info.row.getValue( 'email' ) )}` }>{ info.getValue() as string }</a>,
      },
      defaultColumnDef( 'email' ),
      {
        ...defaultColumnDef( 'proposer', 'Invited By' ),
        cell: info => info.row.getValue( 'proposer' ),
      },
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
    <div style={ { display: 'flex', marginBottom: '0.75em' } }>
      <TableWrapper
        loading={ loading }
        table={ { data: users, columns, additionalTableClasses: ['user-table'] } }
      />
    </div>
  );
};

export default UploaderTable;
