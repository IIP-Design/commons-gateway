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
import { Table, defaultColumnDef } from './Table';
import { escapeQueryStrings } from '../utils/string';

// ////////////////////////////////////////////////////////////////////////////
// Styles and CSS
// ////////////////////////////////////////////////////////////////////////////
import style from '../styles/table.module.scss';

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

    getUsers();
  }, [] );

  const columns = useMemo<ColumnDef<WithUiData<IUploader>>[]>(
    () => [
      {
        ...defaultColumnDef( 'name' ),
        cell: info => <a href={ `/edit-user?id=${escapeQueryStrings( info.row.getValue( 'email' ) )}` }>{ info.getValue() as string }</a>,
      },
      defaultColumnDef( 'email' ),
      {
        ...defaultColumnDef( 'proposer' ),
        cell: info => {
          const { String, Valid } = info.getValue() as any;

          return Valid ? String : null;
        },
      },
      {
        ...defaultColumnDef( 'inviter' ),
        cell: info => {
          const { String, Valid } = info.getValue() as any;

          return Valid ? String : null;
        },
      },
      {
        ...defaultColumnDef( 'active', 'Status' ),
        cell: info => {
          const isPending = info.row.original.pending as boolean;
          const isActive = info.getValue() as boolean;


          return (
            <span className={ style.status }>
              <span className={ isActive ? ( isPending ? style.pending : style.active ) : style.inactive } />
              { isActive ? ( isPending ? 'Pending' : 'Active' ) : 'Inactive' }
            </span>
          );
        },
      },
    ],
    [],
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
        : <p className={ style['no-data'] }>No data to show</p> }
    </div>
  );
};

export default UploaderTable;
